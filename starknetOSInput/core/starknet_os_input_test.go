package osinput

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/stretchr/testify/assert"
)

// Note: given the depreceated classes (dummy_token,dummy_account,token_for_testing)
// and initial transactions (deploy_token_tx, fund_account_tx, deploy_account_tx)
// we should be able to build the same initial and new state as run_os.py
// Given these two states+txns etc, we should be able to test GenerateStarknetOSInput (==get_os_hints).
// The final goal is to compute a StarknetOsInput equivalent to testdata/os_input.json
func TestGenerateStarknetOSInput(t *testing.T) {
	// mockCtrl := gomock.NewController(t)

	// bc := blockchain.New(pebble.NewMemTest(t), utils.Sepolia)
	// // require.NoError(t, bc.StoreGenesis(core.EmptyStateDiff(), nil))
	// mockVM := mocks.NewMockVM(mockCtrl)

	t.Run("empty test", func(t *testing.T) {
		result, err := TxnExecInfo(nil)
		assert.Error(t, err, errors.New("vmParameters can not be nil"))
		assert.Nil(t, result)
	})

}

func loadInitClasses() ([]core.Cairo0Class, []core.Cairo1Class, error) {
	loadJSON := func(filePath string, target interface{}) {
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("unable to read json file: %v", err)
		}
		if err := json.Unmarshal(data, target); err != nil {
			panic(fmt.Errorf("unable to unmarshal json compiled class: %v", err))
		}
	}

	testContractClass := new(core.Cairo0Class)
	loadJSON("testdata/test_contract.json", testContractClass)

	dummyAccountClass := new(core.Cairo0Class)
	loadJSON("testdata/dummy_accout.json", dummyAccountClass)

	dummyTokenClass := new(core.Cairo0Class)
	loadJSON("testdata/dummy_token.json", dummyTokenClass)

	return []core.Cairo0Class{*testContractClass, *dummyAccountClass, *dummyTokenClass}, nil, nil
}

func stringToFelt(s string) *felt.Felt {
	f, err := new(felt.Felt).SetString(s)
	if err != nil {
		panic(err)

	}
	return f
}

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
