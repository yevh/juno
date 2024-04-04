package osinput

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGenerateStarknetOSInput(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	// stateReader := mocks.NewMockStateHistoryReader(mockCtrl)
	// block := &core.Block{}
	// vmInput := VMParameters{}

	t.Run("empty test", func(t *testing.T) {
		result, err := TxnExecInfo(nil)
		assert.Error(t, err, errors.New("vmParameters can not be nil"))
		assert.Nil(t, result)
	})

}
