package osinput

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core"
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
