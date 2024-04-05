package osinput

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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

	t.Run("empty test", func(t *testing.T) {
		result, err := TxnExecInfo(nil)
		assert.Error(t, err, errors.New("vmParameters can not be nil"))
		assert.Nil(t, result)
	})

}
