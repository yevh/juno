package osinput

import (
	"errors"

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

func TxnExecInfo(vmParams *VMParameters) (*[]vm.TransactionTrace, error) {
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

	return &traces, nil
}
