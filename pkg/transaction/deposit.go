package transaction

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// This file define the Deposit transaction

type DepositTx struct {
	// From address of the tx initiator's L1 account
	From common.Address
	// To address to deposit to
	To     common.Address
	Amount *big.Int
	Token  uint16
}

func (tx DepositTx) isZionPriTx() {}
