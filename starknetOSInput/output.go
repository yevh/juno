package osinput

import (
	"errors"

	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/core/felt"
)

// Note: assumes the block and state have already been stored
func PopulateOSInput(block core.Block, reader core.StateHistoryReader, execInfo []TransactionExecutionInfo) (*StarknetOsInput, error) {
	osinput := &StarknetOsInput{}

	// Todo
	stateSelector, err := get_os_state_selector(osinput.Transactions, execInfo, &osinput.GeneralConfig)
	if err != nil {
		return nil, err
	}

	// Todo: commitment info
	contractStateCommitmentInfo := getContractStateCommitmentInfo()
	classStateCommitmentInfo := getClassStateCommitmentInfo()
	osinput.ContractStateCommitmentInfo = contractStateCommitmentInfo
	osinput.ContractClassCommitmentInfo = classStateCommitmentInfo

	// populate transactions, deprecrated classes, and classes
	for _, tx := range block.Transactions {
		osinput.Transactions = append(osinput.Transactions, tx)

		switch decTxn := tx.(type) {
		case *core.DeclareTransaction:
			decClass, err := reader.Class(decTxn.ClassHash)
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

	osinput.ClassHashToCompiledClassHash, err = getClassHashToCompiledClassHash(reader, stateSelector.ClassHashes)

	// Todo: CompiledClassVisitedPcs??

	// Todo: contracts
	contracts := getContracts(stateSelector.ContractAddresses)
	osinput.Contracts = contracts

	osinput.GeneralConfig = loadExampleStarknetOSConfig()

	osinput.BlockHash = *block.Hash
	return osinput, nil
}

func loadExampleStarknetOSConfig() StarknetGeneralConfig {
	seqAddr, _ := new(felt.Felt).SetString("0x6c95526293b61fa708c6cba66fd015afee89309666246952456ab970e9650aa")
	feeTokenAddr, _ := new(felt.Felt).SetString("0x6c95526293b61fa708c6cba66fd015afee89309666246952456ab970e9650aa")
	chainID := new(felt.Felt).SetBytes([]byte("TODO"))
	return StarknetGeneralConfig{
		StarknetOsConfig: StarknetOsConfig{
			ChainID:         *chainID,
			FeeTokenAddress: *feeTokenAddr,
		},
		InvokeTxMaxNSteps: 3000000,
		ValidateMaxNSteps: 1000000,
		ConstantGasPrice:  false,
		SequencerAddress:  *seqAddr,
		CairoResourceFeeWeights: map[string]float64{
			"poseidon_builtin":    0.0,
			"n_steps":             1.0,
			"ecdsa_builtin":       0.0,
			"keccak_builtin":      0.0,
			"range_check_builtin": 0.0,
			"pedersen_builtin":    0.0,
			"output_builtin":      0.0,
			"bitwise_builtin":     0.0,
			"ec_op_builtin":       0.0,
		},
		EnforceL1HandlerFee: true,
		UseKzgDa:            false,
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

func getContractStateCommitmentInfo() CommitmentInfo {

	// Todo: Given the old and new contract Trie, collect all the
	// nodes that were modified

	panic("unimplemented")
}

func getClassStateCommitmentInfo() CommitmentInfo {

	// Todo: Given the old and new class Trie, collect all the
	// nodes that were modified

	panic("unimplemented")
}

func getContracts(contractAddresses []felt.Felt) map[felt.Felt]ContractState {
	// Todo: given a batch of transactions, collect the set of ContractState's
	// for every contract that the contract touched
	panic("unimplemented")
}
