package osinput

import (
	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/juno/utils"
	"github.com/NethermindEth/juno/vm"
)

// Given a state and a set of transactions, execute the transactions in blockifier to return the traces
// These traces make up the TransactionExecutionInfo that we need to populate the OSInputs

type VMParameters struct {
	Txns            []core.Transaction
	DeclaredClasses []core.Class
	PaidFeesOnL1    []*felt.Felt
	BlockInfo       *vm.BlockInfo
	State           core.StateReader
	Network         *utils.Network
	SkipChargeFee   bool
	SkipValidate    bool
	ErrOnRevert     bool
	UseBlobData     bool
}

func executeTransactions(vmParams *VMParameters) (*[]TransactionExecutionInfo, error) {
	_, _, traces, err := vm.New(nil).Execute(
		vmParams.Txns,
		vmParams.DeclaredClasses,
		vmParams.PaidFeesOnL1,
		vmParams.BlockInfo,
		vmParams.State,
		vmParams.Network,
		vmParams.SkipChargeFee,
		vmParams.SkipValidate,
		vmParams.ErrOnRevert,
		vmParams.UseBlobData,
	)
	if err != nil {
		return nil, err
	}

	execInfos := []TransactionExecutionInfo{}
	for _, trace := range traces {
		execInfo := convertVMTraceToTxnExecInfo(trace)
		execInfos = append(execInfos, execInfo)
	}

	return &execInfos, nil
}

func convertVMTraceToTxnExecInfo(trace vm.TransactionTrace) TransactionExecutionInfo {
	panic("todo")
}
