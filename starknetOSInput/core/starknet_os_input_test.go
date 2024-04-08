package osinput

import (
	"testing"

	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/juno/db/pebble"
	"github.com/NethermindEth/juno/mocks"
	"github.com/NethermindEth/juno/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// Note: given the depreceated classes (dummy_token,dummy_account,token_for_testing)
// and initial transactions (deploy_token_tx, fund_account_tx, deploy_account_tx)
// we should be able to build the same initial and new state as run_os.py
// Given these two states+txns etc, we should be able to test GenerateStarknetOSInput (==get_os_hints).
// The final goal is to compute a StarknetOsInput equivalent to testdata/os_input.json
func TestGenerateStarknetOSInput(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)

	network := utils.Sepolia

	testDB := pebble.NewMemTest(t)
	txn, err := testDB.NewTransaction(true)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, txn.Discard())
	})

	state := core.NewState(txn)
	mockVM := mocks.NewMockVM(mockCtrl)
	t.Run("empty inputs", func(t *testing.T) {

		// Update the contract state trie to inc the contracts "0" and "1"
		// These will always be checked for changes etc
		zeroHash := utils.HexToFelt(t, "0x0")
		oneHash := utils.HexToFelt(t, "0x1")
		newClasses := map[felt.Felt]core.Class{
			*zeroHash: &core.Cairo0Class{},
			*oneHash:  &core.Cairo0Class{},
		}
		su := &core.StateUpdate{
			OldRoot:   &felt.Zero,
			NewRoot:   utils.HexToFelt(t, "0x0"),
			BlockHash: &felt.Zero,
			StateDiff: &core.StateDiff{
				Nonces: map[felt.Felt]*felt.Felt{
					*zeroHash: &felt.Zero,
					*oneHash:  &felt.Zero,
				},
				DeclaredV0Classes: []*felt.Felt{zeroHash, oneHash},
				DeclaredV1Classes: nil,
			},
		}
		err := state.Update(0, su, newClasses)
		require.NoError(t, err)

		vmParas := VMParameters{
			Txns:            nil,
			DeclaredClasses: nil,
			PaidFeesOnL1:    nil,
			State:           state,
			Network:         &network,
			SkipChargeFee:   false,
			SkipValidate:    false,
			ErrOnRevert:     false,
			UseBlobData:     false,
		}
		block := core.Block{
			Header: &core.Header{
				Hash: utils.HexToFelt(t, "0xdeadbeef"),
			},
		}
		mockVM.EXPECT().Execute(vmParas.Txns, vmParas.DeclaredClasses, vmParas.PaidFeesOnL1,
			vmParas.BlockInfo, state, vmParas.Network, vmParas.SkipChargeFee, vmParas.SkipValidate, vmParas.ErrOnRevert, vmParas.UseBlobData).Return(nil, nil, nil, nil)
		_, err = GenerateStarknetOSInput(&block, state, state, mockVM, vmParas)
		require.NoError(t, err)
	})
}

// func loadInitClasses() ([]core.Cairo0Class, []core.Cairo1Class, error) {
// 	loadJSON := func(filePath string, target interface{}) {
// 		data, err := os.ReadFile(filePath)
// 		if err != nil {
// 			log.Fatalf("unable to read json file: %v", err)
// 		}
// 		if err := json.Unmarshal(data, target); err != nil {
// 			panic(fmt.Errorf("unable to unmarshal json compiled class: %v", err))
// 		}
// 	}

// 	testContractClass := new(core.Cairo0Class)
// 	loadJSON("testdata/test_contract.json", testContractClass)

// 	dummyAccountClass := new(core.Cairo0Class)
// 	loadJSON("testdata/dummy_accout.json", dummyAccountClass)

// 	dummyTokenClass := new(core.Cairo0Class)
// 	loadJSON("testdata/dummy_token.json", dummyTokenClass)

// 	return []core.Cairo0Class{*testContractClass, *dummyAccountClass, *dummyTokenClass}, nil, nil
// }

// func stringToFelt(s string) *felt.Felt {
// 	f, err := new(felt.Felt).SetString(s)
// 	if err != nil {
// 		panic(err)

// 	}
// 	return f
// }

// func getInitTxns(depClasses []core.Cairo0Class) ([]core.Transaction, error) {

// 	deployDummyTokenTransaction := core.DeployAccountTransaction{
// 		DeployTransaction: core.DeployTransaction{
// 			TransactionHash:     stringToFelt("22688876470218804543887986415455328819098091743988100398197353790124740200"),
// 			Version:             new(core.TransactionVersion).SetUint64(1),
// 			ContractAddressSalt: new(felt.Felt).SetUint64(0),
// 			ClassHash:           stringToFelt("3531298130119845387864440863187980726515137569165069484670944625223023734186"),
// 			ConstructorCalldata: []byte{},
// 		},
// 		MaxFee:               new(felt.Felt).SetBigInt(big.NewInt(1267650600228229401496703205376)),
// 		TransactionSignature: []*felt.Felt{},
// 		Nonce:                new(felt.Felt).SetUint64(0),
// 	}

// 	deployDummyAccountTransaction := core.DeployAccountTransaction{
// 		DeployTransaction: core.DeployTransaction{
// 			TransactionHash:     new(felt.Felt).SetBigInt(big.NewInt(96374521715508826444566467091393680183010464453336720810014746622481735761)),
// 			Version:             new(core.TransactionVersion).SetUint64(1),
// 			SenderAddress:       new(felt.Felt).SetBigInt(big.NewInt(2618767603815038378512366346550627731109766804643583016834052353912473402832)),
// 			ContractAddressSalt: new(felt.Felt).SetUint64(0),
// 			ClassHash:           big.NewInt(646245114977324210659279014519951538684823368221946044944492064370769527799),
// 			ConstructorCalldata: []byte{},
// 		},
// 		MaxFee:               new(felt.Felt).SetBigInt(big.NewInt(1267650600228229401496703205376)),
// 		TransactionSignature: []*felt.Felt{},
// 		Nonce:                new(felt.Felt).SetUint64(0),
// 	}

// 	fundInvokeTransaction := core.InvokeTransaction{
// 		TransactionHash:      new(felt.Felt).SetBigInt(big.NewInt(2852915394592604060963909836770256627436576776991723431020681987492769528640)),
// 		Version:              new(core.TransactionVersion).SetUint64(1),
// 		MaxFee:               new(felt.Felt).SetBigInt(big.NewInt(1267650600228229401496703205376)),
// 		TransactionSignature: []*felt.Felt{},
// 		Nonce:                new(felt.Felt).SetUint64(1),
// 		SenderAddress:        new(felt.Felt).SetBigInt(big.NewInt(1470089414715992704702781317133162679047468004062084455026636858461958198968)),
// 		EntryPointSelector:   new(felt.Felt).SetBigInt(big.NewInt(617075754465154585683856897856256838130216341506379215893724690153393808813)),
// 		EntryPointType:       core.EntryPointType.EXTERNAL,
// 		CallData: []*felt.Felt{
// 			new(felt.Felt).SetBigInt(big.NewInt(1470089414715992704702781317133162679047468004062084455026636858461958198968)),
// 			new(felt.Felt).SetBigInt(big.NewInt(232670485425082704932579856502088130646006032362877466777181098476241604910)),
// 			new(felt.Felt).SetUint64(3),
// 			new(felt.Felt).SetBigInt(big.NewInt(2618767603815038378512366346550627731109766804643583016834052353912473402832)),
// 			new(felt.Felt).SetBigInt(big.NewInt(1329227995784915872903807060280344576)),
// 			new(felt.Felt).SetUint64(0),
// 		},
// 	}

// 	return []core.Transaction{&deployDummyAccountTransaction, &deployDummyTokenTransaction, &fundInvokeTransaction}, nil
// }
