package blockchain

import (
	"context"
	"github.com/NethermindEth/juno/adapters/core2p2p"
	"github.com/NethermindEth/juno/adapters/p2p2core"
	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/juno/core/trie"
	"github.com/NethermindEth/juno/db"
	"github.com/NethermindEth/juno/p2p/starknet/spec"
	"github.com/NethermindEth/juno/utils/iter"
	"github.com/pkg/errors"
)

type ContractRangeStreamingResult struct {
	ContractsRoot *felt.Felt
	ClassesRoot   *felt.Felt
	Range         []*spec.ContractState
	RangeProof    *spec.PatriciaRangeProof
}

type StorageRangeRequest struct {
	StateRoot     *felt.Felt
	ChunkPerProof uint64 // Missing in spec
	Queries       []*spec.StorageRangeQuery
}

type StorageRangeStreamingResult struct {
	ContractsRoot *felt.Felt
	ClassesRoot   *felt.Felt
	Range         []*spec.ContractStoredValue
	RangeProof    *spec.PatriciaRangeProof
}

type ClassRangeStreamingResult struct {
	ContractsRoot *felt.Felt
	ClassesRoot   *felt.Felt
	Range         *spec.Classes
	RangeProof    *spec.PatriciaRangeProof
}

type SnapServer interface {
	GetContractRange(ctx context.Context, request *spec.ContractRangeRequest) iter.Seq2[*ContractRangeStreamingResult, error]
	GetStorageRange(ctx context.Context, request *StorageRangeRequest) iter.Seq2[*StorageRangeStreamingResult, error]
	GetClassRange(ctx context.Context, request *spec.ClassRangeRequest) iter.Seq2[*ClassRangeStreamingResult, error]
	GetClasses(ctx context.Context, classHashes []*felt.Felt) ([]*spec.Class, error)
}

var (
	_                  SnapServer = &Blockchain{}
	ErrMissingSnapshot            = errors.New("missing snapshot")
)

const maxNodePerRequest = 1024 * 1024 // I just want it to process faster
func determineMaxNodes(specifiedMaxNodes uint64) uint64 {
	if specifiedMaxNodes == 0 {
		return 1024 * 16
	}

	if specifiedMaxNodes < maxNodePerRequest {
		return specifiedMaxNodes
	}
	return maxNodePerRequest
}

func (b *Blockchain) findSnapshotMatching(filter func(record *snapshotRecord) bool) (*snapshotRecord, error) {
	var snapshot *snapshotRecord
	for _, record := range b.snapshots {
		if filter(record) {
			snapshot = record
			break
		}
	}

	if snapshot == nil {
		return nil, ErrMissingSnapshot
	}

	return snapshot, nil
}

func iterateWithLimit(
	srcTrie *trie.Trie,
	startAddr *felt.Felt,
	limitAddr *felt.Felt,
	maxNode uint64,
	consumer func(key, value *felt.Felt) error,
) ([]trie.ProofNode, error) {
	pathes := make([]*felt.Felt, 0)
	hashes := make([]*felt.Felt, 0)

	// TODO: Verify class trie
	var startPath *felt.Felt
	var endPath *felt.Felt
	count := uint64(0)
	neverStopped, err := srcTrie.Iterate(startAddr, func(key *felt.Felt, value *felt.Felt) (bool, error) {
		// Need at least one.
		if limitAddr != nil && key.Cmp(limitAddr) > 1 && count > 0 {
			return false, nil
		}

		if startPath == nil {
			startPath = key
		}

		pathes = append(pathes, key)
		hashes = append(hashes, value)

		err := consumer(key, value)
		if err != nil {
			return false, err
		}

		endPath = key
		count++
		if count >= maxNode {
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		return nil, err
	}

	if neverStopped && startAddr.Equal(&felt.Zero) {
		return nil, nil // No need for proof
	}
	if startPath == nil {
		return nil, nil // No need for proof
	}

	return srcTrie.RangeProof(startPath, endPath)
}

func (b *Blockchain) GetClassRange(ctx context.Context, request *spec.ClassRangeRequest) iter.Seq2[*ClassRangeStreamingResult, error] {
	return func(yield func(*ClassRangeStreamingResult, error) bool) {
		stateRoot := p2p2core.AdaptHash(request.Root)

		snapshot, err := b.findSnapshotMatching(func(record *snapshotRecord) bool {
			return record.stateRoot.Equal(stateRoot)
		})

		if err != nil {
			yield(nil, err)
			return
		}

		s := core.NewState(snapshot.txn)

		contractRoot, classRoot, err := s.StateAndClassRoot()
		if err != nil {
			yield(nil, err)
			return
		}

		// TODO: Verify class trie
		ctrie, classCloser, err := s.ClassTrie()
		if err != nil {
			yield(nil, err)
			return
		}
		defer classCloser()

		response := &spec.Classes{
			Classes: make([]*spec.Class, 0),
		}

		startAddr := p2p2core.AdaptHash(request.Start)
		limitAddr := p2p2core.AdaptHash(request.End)

		// TODO: loop this
		proofs, err := iterateWithLimit(ctrie, startAddr, limitAddr, determineMaxNodes(uint64(request.ChunksPerProof)), func(key, value *felt.Felt) error {
			// response.Classes = append(response) // !!!!!
			panic("not implemented")
			return nil
		})

		if err != nil {
			yield(nil, err)
			return
		}

		yield(&ClassRangeStreamingResult{
			ContractsRoot: contractRoot,
			ClassesRoot:   classRoot,
			Range:         response,
			RangeProof:    Core2P2pProof(proofs),
		}, err)
	}
}

func Core2P2pProof(proofs []trie.ProofNode) *spec.PatriciaRangeProof {
	panic("not implemented")
}

func (b *Blockchain) GetContractRange(ctx context.Context, request *spec.ContractRangeRequest) iter.Seq2[*ContractRangeStreamingResult, error] {
	return func(yield func(*ContractRangeStreamingResult, error) bool) {
		stateRoot := p2p2core.AdaptHash(request.StateRoot)
		snapshot, err := b.findSnapshotMatching(func(record *snapshotRecord) bool {
			return record.stateRoot.Equal(stateRoot)
		})
		if err != nil {
			yield(nil, err)
			return
		}

		s := core.NewState(snapshot.txn)
		contractRoot, classRoot, err := s.StateAndClassRoot()
		if err != nil {
			yield(nil, err)
			return
		}

		// TODO: Verify class trie
		strie, scloser, err := s.StorageTrie()
		if err != nil {
			yield(nil, err)
			return
		}
		defer scloser()

		startAddr := p2p2core.AdaptAddress(request.Start)
		limitAddr := p2p2core.AdaptAddress(request.End)
		states := []*spec.ContractState{}

		proofs, err := iterateWithLimit(strie, startAddr, limitAddr, determineMaxNodes(uint64(request.ChunksPerProof)), func(key, value *felt.Felt) error {
			classHash, err := s.ContractClassHash(key)
			if err != nil {
				return err
			}

			nonce, err := s.ContractNonce(key)
			if err != nil {
				return err
			}

			ctr, err := s.StorageTrieForAddr(key)
			if err != nil {
				return err
			}

			croot, err := ctr.Root()
			if err != nil {
				return err
			}

			states = append(states, &spec.ContractState{
				Address: core2p2p.AdaptAddress(key),
				Class:   core2p2p.AdaptHash(classHash),
				Storage: core2p2p.AdaptHash(croot),
				Nonce:   nonce.Uint64(),
			})
			return nil
		})

		yield(&ContractRangeStreamingResult{
			ContractsRoot: contractRoot,
			ClassesRoot:   classRoot,
			Range:         states,
			RangeProof:    Core2P2pProof(proofs),
		}, nil)
	}
}

func (b *Blockchain) GetClasses(ctx context.Context, felts []*felt.Felt) ([]*spec.Class, error) {
	classes := make([]*spec.Class, len(felts))
	err := b.database.View(func(txn db.Transaction) error {
		state := core.NewState(txn)
		for i, f := range felts {
			d, err := state.Class(f)
			if err != nil {
				return err
			}
			classes[i] = core2p2p.AdaptClass(d.Class)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return classes, nil
}

func (b *Blockchain) GetClassesLoc(felts []*felt.Felt) ([]core.Class, error) {
	classes := make([]core.Class, len(felts))
	err := b.database.View(func(txn db.Transaction) error {
		state := core.NewState(txn)
		for i, f := range felts {
			d, err := state.Class(f)
			if err != nil {
				return err
			}
			classes[i] = d.Class
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return classes, nil
}

func (b *Blockchain) GetStorageRange(ctx context.Context, request *StorageRangeRequest) iter.Seq2[*StorageRangeStreamingResult, error] {
	return func(yield func(*StorageRangeStreamingResult, error) bool) {
		stateRoot := request.StateRoot
		snapshot, err := b.findSnapshotMatching(func(record *snapshotRecord) bool {
			return record.stateRoot.Equal(stateRoot)
		})
		if err != nil {
			yield(nil, err)
			return
		}

		s := core.NewState(snapshot.txn)
		contractRoot, classRoot, err := s.StateAndClassRoot()
		if err != nil {
			yield(nil, err)
			return
		}

		curNodeLimit := int64(1000000)

		for _, query := range request.Queries {
			if ctxerr := ctx.Err(); ctxerr != nil {
				break
			}

			contractLimit := uint64(curNodeLimit)

			strie, err := s.StorageTrieForAddr(p2p2core.AdaptAddress(query.Address))
			if err != nil {
				yield(nil, err)
				return
			}

			handled, err := b.handleStorageRangeRequest(ctx, strie, query, request.ChunkPerProof, contractLimit, func(values []*spec.ContractStoredValue, proofs []trie.ProofNode) {
				yield(&StorageRangeStreamingResult{
					ContractsRoot: contractRoot,
					ClassesRoot:   classRoot,
					Range:         values,
					RangeProof:    Core2P2pProof(proofs),
				}, nil)
			})

			if err != nil {
				yield(nil, err)
				return
			}

			curNodeLimit -= handled

			if curNodeLimit <= 0 {
				break
			}
		}
	}
}

func (b *Blockchain) handleStorageRangeRequest(
	ctx context.Context,
	trie *trie.Trie,
	request *spec.StorageRangeQuery,
	maxChunkPerProof uint64,
	nodeLimit uint64,
	yield func([]*spec.ContractStoredValue, []trie.ProofNode)) (int64, error) {

	totalSent := int64(0)
	finished := false
	startAddr := p2p2core.AdaptFelt(request.Start.Key)
	endAddr := p2p2core.AdaptFelt(request.End.Key)

	for !finished {
		if ctxerr := ctx.Err(); ctxerr != nil {
			break
		}

		response := []*spec.ContractStoredValue{}

		limit := maxChunkPerProof
		if nodeLimit < limit {
			limit = nodeLimit
		}

		proofs, err := iterateWithLimit(trie, startAddr, endAddr, limit, func(key, value *felt.Felt) error {
			response = append(response, &spec.ContractStoredValue{
				Key:   core2p2p.AdaptFelt(key),
				Value: core2p2p.AdaptFelt(key),
			})
			return nil
		})

		if err != nil {
			return 0, err
		}

		if len(response) == 0 {
			finished = true
			return totalSent, nil
		}

		yield(response, proofs)

		totalSent += totalSent
		nodeLimit -= limit
	}

	return totalSent, nil
}

type snapshotRecord struct {
	stateRoot     *felt.Felt
	contractsRoot *felt.Felt
	classRoot     *felt.Felt
	blockHash     *felt.Felt
	txn           db.Transaction
	closer        func() error
}

func (b *Blockchain) seedSnapshot() error {
	headheader, err := b.HeadsHeader()
	if err != nil {
		return err
	}

	state, scloser, err := b.HeadState()
	if err != nil {
		return err
	}

	defer scloser()

	stateS := state.(*core.State)
	contractsRoot, theclassroot, err := stateS.StateAndClassRoot()
	if err != nil {
		return err
	}

	thestateroot, err := stateS.Root()
	if err != nil {
		return err
	}

	txn, closer, err := b.database.PersistedView()
	if err != nil {
		return err
	}

	dbsnap := snapshotRecord{
		stateRoot:     thestateroot,
		contractsRoot: contractsRoot,
		classRoot:     theclassroot,
		blockHash:     headheader.Hash,
		txn:           txn,
		closer:        closer,
	}

	// TODO: Reorgs
	b.snapshots = append(b.snapshots, &dbsnap)
	if len(b.snapshots) > 128 {
		toremove := b.snapshots[0]
		err = toremove.closer()
		if err != nil {
			return err
		}

		// TODO: I think internally, it keep the old array.
		// maybe the append copy it to a new array, who knows...
		b.snapshots = b.snapshots[1:]
	}

	return nil
}

func (b *Blockchain) Close() {
	for _, snapshot := range b.snapshots {
		snapshot.closer()
	}
}
