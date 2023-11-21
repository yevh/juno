package starknet

import (
	"crypto/rand"
	"slices"

	"github.com/NethermindEth/juno/adapters/core2p2p"
	"github.com/NethermindEth/juno/blockchain"
	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/juno/p2p/starknet/spec"
	"github.com/NethermindEth/juno/utils"
	"google.golang.org/protobuf/proto"
)

type blockBodyStep int

const (
	_ blockBodyStep = iota

	sendDiff // initial
	sendClasses
	sendProof
	sendTxs
	sendReceipts
	sendBlockFin
	terminal // final step
)

type blockBodyIterator struct {
	log         utils.SimpleLogger
	stateReader core.StateReader
	stateCloser func() error

	step  blockBodyStep
	block *core.Block

	stateUpdate *core.StateUpdate
}

func newBlockBodyIterator(bcReader blockchain.Reader, block *core.Block, log utils.SimpleLogger) (*blockBodyIterator, error) {
	stateUpdate, err := bcReader.StateUpdateByNumber(block.Number)
	if err != nil {
		return nil, err
	}

	stateReader, closer, err := bcReader.StateAtBlockNumber(block.Number)
	if err != nil {
		return nil, err
	}

	return &blockBodyIterator{
		step:        sendDiff,
		block:       block,
		log:         log,
		stateReader: stateReader,
		stateCloser: closer,
		stateUpdate: stateUpdate,
	}, nil
}

func (b *blockBodyIterator) hasNext() bool {
	return slices.Contains([]blockBodyStep{
		sendDiff,
		sendClasses,
		sendProof,
		sendBlockFin,
		sendTxs,
		sendReceipts,
	}, b.step)
}

// Either BlockBodiesResponse_Diff, *_Classes, *_Proof, *_Fin
func (b *blockBodyIterator) next() (msg proto.Message, valid bool) {
	switch b.step {
	case sendDiff:
		msg, valid = b.diff()
		b.step = sendClasses
	case sendClasses:
		msg, valid = b.classes()
		b.step = sendProof
	case sendProof:
		msg, valid = b.proof()
		b.step = sendTxs
	case sendTxs:
		msg, valid = b.transactions()
		b.step = sendReceipts
	case sendReceipts:
		msg, valid = b.receipts()
		b.step = sendBlockFin
	case sendBlockFin:
		// fin changes step to terminal internally
		msg, valid = b.fin()
	case terminal:
		panic("next called on terminal step")
	default:
		b.log.Errorw("Unknown step in blockBodyIterator", "step", b.step)
	}

	return
}

func (b *blockBodyIterator) transactions() (proto.Message, bool) {
	return &spec.BlockBodiesResponse{
		Id: b.blockID(),
		BodyMessage: &spec.BlockBodiesResponse_Transactions{
			Transactions: &spec.Transactions{
				Items: utils.Map(b.block.Transactions, core2p2p.AdaptTransaction),
			},
		},
	}, true
}

func (b *blockBodyIterator) receipts() (proto.Message, bool) {
	var receipts []*spec.Receipt
	for i, receipt := range b.block.Receipts {
		receipts = append(receipts, core2p2p.AdaptReceipt(receipt, b.block.Transactions[i]))
	}

	return &spec.BlockBodiesResponse{
		Id: b.blockID(),
		BodyMessage: &spec.BlockBodiesResponse_Receipts{
			Receipts: &spec.Receipts{
				Items: receipts,
			},
		},
	}, true
}

func (b *blockBodyIterator) classes() (proto.Message, bool) {
	classesM := make(map[felt.Felt]*spec.Class)

	stateDiff := b.stateUpdate.StateDiff

	for _, hash := range stateDiff.DeclaredV0Classes {
		cls, err := b.stateReader.Class(hash)
		if err != nil {
			return b.fin()
		}

		classesM[*hash] = core2p2p.AdaptClass(cls.Class)
	}
	for classHash := range stateDiff.DeclaredV1Classes {
		cls, err := b.stateReader.Class(&classHash)
		if err != nil {
			return b.fin()
		}

		hash, err := cls.Class.Hash()
		if err != nil {
			return b.fin()
		}
		classesM[classHash] = core2p2p.AdaptClass(cls.Class)
	}
	for _, classHash := range stateDiff.DeployedContracts {
		if _, ok := classesM[*classHash]; ok {
			// Skip if the class was declared and deployed in the same block. Otherwise, there will be duplicate classes.
			continue
		}
		cls, err := b.stateReader.Class(classHash)
		if err != nil {
			return b.fin()
		}

		var compiledHash *felt.Felt
		switch cairoClass := cls.Class.(type) {
		case *core.Cairo0Class:
			compiledHash = classHash
		case *core.Cairo1Class:
			compiledHash, err = cairoClass.Hash()
			if err != nil {
				return b.fin()
			}
		default:
			b.log.Errorw("Unknown cairo class", "cairoClass", cairoClass)
			return b.fin()
		}

		classesM[*compiledHash] = core2p2p.AdaptClass(cls.Class)
	}

	var classes []*spec.Class
	for _, c := range classesM {
		c := c
		classes = append(classes, c)
	}

	return &spec.BlockBodiesResponse{
		Id: b.blockID(),
		BodyMessage: &spec.BlockBodiesResponse_Classes{
			Classes: &spec.Classes{
				Domain:  0,
				Classes: classes,
			},
		},
	}, true
}

type contractDiff struct {
	address      *felt.Felt
	classHash    *felt.Felt
	storageDiffs map[felt.Felt]*felt.Felt
	nonce        *felt.Felt
}

func (b *blockBodyIterator) diff() (proto.Message, bool) {
	var err error
	diff := b.stateUpdate.StateDiff

	modifiedContracts := make(map[felt.Felt]*contractDiff)
	initContractDiff := func(addr *felt.Felt) (*contractDiff, error) {
		var cHash *felt.Felt
		cHash, err = b.stateReader.ContractClassHash(addr)
		if err != nil {
			return nil, err
		}
		return &contractDiff{address: addr, classHash: cHash}, nil
	}

	for addr, n := range diff.Nonces {
		addr := addr // copy
		cDiff, ok := modifiedContracts[addr]
		if !ok {
			cDiff, err = initContractDiff(&addr)
			if err != nil {
				b.log.Errorw("Failed to get class hash", "err", err)
				return b.fin()
			}
			modifiedContracts[addr] = cDiff
		}
		cDiff.nonce = n
	}

	for addr, sDiff := range diff.StorageDiffs {
		addr := addr // copy
		cDiff, ok := modifiedContracts[addr]
		if !ok {
			cDiff, err = initContractDiff(&addr)
			if err != nil {
				b.log.Errorw("Failed to get class hash", "err", err)
				return b.fin()
			}
			modifiedContracts[addr] = cDiff
		}
		cDiff.storageDiffs = sDiff
	}

	var contractDiffs []*spec.StateDiff_ContractDiff
	for _, c := range modifiedContracts {
		contractDiffs = append(contractDiffs, core2p2p.AdaptStateDiff(c.address, c.classHash, c.nonce, c.storageDiffs))
	}

	return &spec.BlockBodiesResponse{
		Id: b.blockID(),
		BodyMessage: &spec.BlockBodiesResponse_Diff{
			Diff: &spec.StateDiff{
				Domain:        0,
				ContractDiffs: contractDiffs,
			},
		},
	}, true
}

func (b *blockBodyIterator) fin() (proto.Message, bool) {
	b.step = terminal
	if err := b.stateCloser(); err != nil {
		b.log.Errorw("Call to state closer failed", "err", err)
	}
	return &spec.BlockBodiesResponse{
		Id:          b.blockID(),
		BodyMessage: &spec.BlockBodiesResponse_Fin{},
	}, true
}

func (b *blockBodyIterator) proof() (proto.Message, bool) {
	// proof size is currently 142K
	proof := make([]byte, 142*1024) //nolint:gomnd
	_, err := rand.Read(proof)
	if err != nil {
		b.log.Errorw("Failed to generate rand proof", "err", err)
		return b.fin()
	}

	return &spec.BlockBodiesResponse{
		Id: b.blockID(),
		BodyMessage: &spec.BlockBodiesResponse_Proof{
			Proof: &spec.BlockProof{
				Proof: proof,
			},
		},
	}, true
}

func (b *blockBodyIterator) blockID() *spec.BlockID {
	return core2p2p.AdaptBlockID(b.block.Header)
}
