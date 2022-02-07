package zcontract

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// some helper func of event
const (
	DepositEvent     = "NewDepositRequest"
	DepositRequested = "NewDepositRequest(address,address,uint64,uint16,uint96)"
)

type ContractTopics struct {
	DepositRequested common.Hash
}

func NewContracTopics() *ContractTopics {
	depositRequested := crypto.Keccak256Hash([]byte(DepositRequested))
	return &ContractTopics{
		DepositRequested: depositRequested,
	}
}

func ParseNewDepositRequest(logData []byte) ZionCNewDepositRequest {
	contracAbi, _ := abi.JSON(strings.NewReader(string(ZionCABI)))

	elog := ZionCNewDepositRequest{}
	contracAbi.UnpackIntoInterface(&elog, DepositEvent, logData)
	return elog
}
