package osinput

import (
	"errors"

	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/juno/rpc"
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

func TxnExecInfo(vmParams *VMParameters) (*[]TransactionExecutionInfo, error) {
	if vmParams == nil {
		return nil, errors.New("vmParameters can not be nil")
	}
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

	// convert traces to TransactionExecutionInfo
	var txnExecInfo []TransactionExecutionInfo
	for _, trace := range traces {
		txnExecInfo = append(txnExecInfo, adaptVMTraceToTxnExecInfo(trace))
	}

	return &txnExecInfo, nil
}

// Todo: finish fields that we need here
func adaptVMTraceToTxnExecInfo(trace vm.TransactionTrace) TransactionExecutionInfo {
	return TransactionExecutionInfo{
		ValidateInfo:    adaptFnInvocationToCallInfo(*trace.ValidateInvocation),
		FeeTransferInfo: adaptFnInvocationToCallInfo(*trace.FeeTransferInvocation),
		TxType:          (*rpc.TransactionType)(&trace.Type),
	}
}

func adaptFnInvocationToCallInfo(fnInvoc vm.FunctionInvocation) *CallInfo {
	var epType EntryPointType
	switch fnInvoc.EntryPointType {
	case "External":
		epType = External
	case "L1Handler":
		epType = L1Handler
	case "Constructor":
		epType = Constructor
	default:
		panic("unknown EntryPointType")
	}

	var callType CallType
	switch fnInvoc.CallType {
	case "Call":
		callType = Call
	case "Delegate":
		callType = Delegate
	default:
		panic("unknown CallType")
	}

	return &CallInfo{
		CallerAddress:      fnInvoc.CallerAddress,
		CallType:           &callType,
		ContractAddress:    fnInvoc.ContractAddress,
		ClassHash:          fnInvoc.ClassHash,
		EntryPointSelector: *fnInvoc.EntryPointSelector,
		EntryPointType:     &epType,
		Calldata:           fnInvoc.Calldata,
		// GasConsumed:         0,   // Not present in FunctionInvocation
		// FailureFlag:         0,   // Not present in FunctionInvocation
		// Retdata:             nil, // Not present in FunctionInvocation
		// ExecutionResources: fi.ExecutionResources,	// Todo: Do we need this?
		// Events:             fi.Events,				// Todo: Do we need this?
		// L2ToL1Messages:     fi.Messages,				// Todo: Do we need this?
		// StorageReadValues:   nil, // Not present in FunctionInvocation
		// AccessedStorageKeys: nil, // Not present in FunctionInvocation
		// InternalCalls: internalCalls,				// Todo: Do we need this?
		// CodeAddress:         nil, // Not present in FunctionInvocation
	}
}
