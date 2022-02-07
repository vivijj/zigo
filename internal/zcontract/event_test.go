package zcontract

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestEvent(t *testing.T) {
	a := NewContracTopics()

	// The default value
	assert.Equal(
		t,
		common.HexToHash("0x53a959842d0c8b0e8e28011927f401d5d1258c0db5ebb6043fd6949220fdc778"),
		a.DepositRequested,
	)
}
