package osinput

import (
	"errors"

	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/juno/vm"
)

// GenerateStarknetOSInput generates the starknet OS input, given a block, state and classes.
func GenerateStarknetOSInput(block *core.Block, oldstate core.StateHistoryReader, newstate core.StateHistoryReader, vm vm.VM, vmInput VMParameters) (*StarknetOsInput, error) {
	txnExecInfo, err := TxnExecInfo(vm, &vmInput)
	if err != nil {
		return nil, err
	}

	return calculateOSInput(*block, oldstate, newstate, *txnExecInfo)
}

// Todo: update to use vm.TransactionTrace instead of TransactionExecutionInfo?
func calculateOSInput(block core.Block, oldstate core.StateHistoryReader, newstate core.StateHistoryReader, execInfo []TransactionExecutionInfo) (*StarknetOsInput, error) {
	osinput := &StarknetOsInput{}

	// Todo: complete
	stateSelector, err := get_os_state_selector(osinput.Transactions, execInfo, &osinput.GeneralConfig)
	if err != nil {
		return nil, err
	}

	classHashToCompiledClassHash, err := getClassHashToCompiledClassHash(oldstate, stateSelector.ClassHashes)
	if err != nil {
		return nil, err
	}

	// Todo: commitment info
	contractStateCommitmentInfo := getContractStateCommitmentInfo(oldstate, newstate, stateSelector.ContractAddresses)
	classStateCommitmentInfo := getClassStateCommitmentInfo(oldstate, newstate, classHashToCompiledClassHash)
	osinput.ContractStateCommitmentInfo = contractStateCommitmentInfo
	osinput.ContractClassCommitmentInfo = classStateCommitmentInfo

	for _, tx := range block.Transactions {
		osinput.Transactions = append(osinput.Transactions, tx)

		switch decTxn := tx.(type) {
		case *core.DeclareTransaction:
			decClass, err := oldstate.Class(decTxn.ClassHash)
			if err != nil {
				return nil, err
			}
			switch class := decClass.Class.(type) {
			case *core.Cairo0Class:
				osinput.DeprecatedCompiledClasses[*decTxn.ClassHash] = *class
			case *core.Cairo1Class:
				osinput.CompiledClasses[*decTxn.ClassHash] = *class.Compiled
			default:
				return nil, errors.New("unknown class type")
			}
		}
	}

	osinput.ClassHashToCompiledClassHash, err = getInitialClassHashToCompiledClassHash(oldstate, classHashToCompiledClassHash)

	// Todo: CompiledClassVisitedPcs??

	contractStates, err := getContracts(oldstate, stateSelector.ContractAddresses)
	if err != nil {
		return nil, err
	}
	osinput.Contracts = contractStates

	osinput.GeneralConfig = LoadExampleStarknetOSConfig()

	osinput.BlockHash = *block.Hash
	return osinput, nil
}

func LoadExampleStarknetOSConfig() StarknetGeneralConfig {
	seqAddr, _ := new(felt.Felt).SetString("0xfb3d5e8dec6dfef87f2b48dd3fbe9a455ed188636f6638a8f8bce4555f7938")
	feeTokenAddr, _ := new(felt.Felt).SetString("0x3400a86fdc294a70fac1cf84f81a2127419359096b846be9814786d4fc056b8")
	return StarknetGeneralConfig{
		StarknetOsConfig: StarknetOsConfig{
			ChainID:                   *new(felt.Felt).SetBytes([]byte{0x05, 0x9b, 0x01, 0xba, 0x26, 0x2c, 0x99, 0x9f, 0x26, 0x17, 0x41, 0x2f, 0xfb, 0xba, 0x78, 0x0f, 0x80, 0xb0, 0x10, 0x3d, 0x92, 0x8c, 0xbc, 0xe1, 0xae, 0xcb, 0xaa, 0x50, 0xde, 0x90, 0xab, 0xda}),
			DeprecatedFeeTokenAddress: felt.Felt{},
			FeeTokenAddress:           *feeTokenAddr,
		},
		GasPriceBounds:       GasPriceBounds{},
		InvokeTxMaxNSteps:    3000000,
		ValidateMaxNSteps:    1000000,
		DefaultEthPriceInFri: *new(felt.Felt),
		ConstantGasPrice:     false,
		SequencerAddress:     seqAddr,
		CairoResourceFeeWeights: map[string]float64{
			"n_steps":             1.0,
			"output_builtin":      0.0,
			"pedersen_builtin":    0.0,
			"range_check_builtin": 0.0,
			"ecdsa_builtin":       0.0,
			"bitwise_builtin":     0.0,
			"ec_op_builtin":       0.0,
			"keccak_builtin":      0.0,
			"poseidon_builtin":    0.0,
		},
		EnforceL1HandlerFee: true,
		UseKzgDa:            true,
	}
}

func getClassHashToCompiledClassHash(reader core.StateHistoryReader, classHashes []felt.Felt) (map[felt.Felt]felt.Felt, error) {
	classHashToCompiledClassHash := map[felt.Felt]felt.Felt{}
	for _, classHash := range classHashes {
		decClass, err := reader.Class(&classHash)
		if err != nil {
			return nil, err
		}
		switch t := decClass.Class.(type) {
		case *core.Cairo1Class:
			classHashToCompiledClassHash[classHash] = *t.Compiled.Hash()
		}
	}
	return classHashToCompiledClassHash, nil

}

func getInitialClassHashToCompiledClassHash(oldstate core.StateHistoryReader, classHashToCompiledClassHash map[felt.Felt]felt.Felt) (map[felt.Felt]felt.Felt, error) {
	for range classHashToCompiledClassHash {
		panic("unimplemented getInitialClassHashToCompiledClassHash")
	}
	return nil, nil
}

func getContractStateCommitmentInfo(oldstate core.StateHistoryReader, newstate core.StateHistoryReader, contractAddresses []felt.Felt) CommitmentInfo {
	// Todo: Given the old and new contract Trie, collect all the
	// nodes that were modified
	for _, address := range contractAddresses {
		if address.Equal(&felt.Zero) || address.Equal(new(felt.Felt).SetUint64(1)) { // Todo: Hack to make empty state work for initial tests.
			continue
		}
		panic("unimplemented getContractStateCommitmentInfo")
	}
	return CommitmentInfo{}
}

func getClassStateCommitmentInfo(oldstate core.StateHistoryReader, newstate core.StateHistoryReader, classHashToCompiledClassHash map[felt.Felt]felt.Felt) CommitmentInfo {
	// Todo: Given the old and new class Trie, collect all the
	// nodes that were modified
	for range classHashToCompiledClassHash {
		panic("unimplemented getClassStateCommitmentInfo")
	}
	return CommitmentInfo{}
}

func getContracts(reader core.StateHistoryReader, contractAddresses []felt.Felt) (map[felt.Felt]ContractState, error) {
	contractState := map[felt.Felt]ContractState{}
	for _, addr := range contractAddresses {
		root, err := reader.ContractStorageRoot(&addr)
		if err != nil {
			return nil, err
		}
		nonce, err := reader.ContractNonce(&addr)
		if err != nil {
			return nil, err
		}
		hash, err := reader.ContractClassHash(&addr)
		if err != nil {
			return nil, err
		}
		contractState[addr] = ContractState{
			ContractHash: *hash,
			StorageCommitmentTree: PatriciaTree{
				Root:   *root,
				Height: 251, // Todo: Just leave hardcoded?
			},
			Nonce: *nonce,
		}
	}
	return contractState, nil
}
