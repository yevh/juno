package core2pyrunner

import (
	osinputs "github.com/NethermindEth/juno/starknetOSInput/core"
	pyosinputs "github.com/NethermindEth/juno/starknetOSInput/python"
)

// Todo
func ConvertToStarknetOsInputInt(input *osinputs.StarknetOsInput) *pyosinputs.StarknetOsInput {
	output := &pyosinputs.StarknetOsInput{
		// ContractStateCommitmentInfo: input.ContractStateCommitmentInfo,
		// ContractClassCommitmentInfo: input.ContractClassCommitmentInfo,
		// DeprecatedCompiledClasses:    convertMapKeysToInt(input.DeprecatedCompiledClasses),
		// CompiledClasses:              convertMapKeysToInt(input.CompiledClasses),
		// CompiledClassVisitedPcs:      convertMapKeysAndValuesToInt(input.CompiledClassVisitedPcs),
		// Contracts:                    convertMapKeysToInt(input.Contracts),
		// ClassHashToCompiledClassHash: convertMapKeysAndValuesToInt(input.ClassHashToCompiledClassHash),
		// GeneralConfig: input.GeneralConfig,
		// Transactions:  input.Transactions,
		// BlockHash:     int(input.BlockHash),
	}

	return output
}
