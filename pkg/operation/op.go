// Package operation define the zion operation -- the tx after being processed
package operation

import (
	"math/big"

	"github.com/vivijj/zigo/pkg/transaction"
)

type ZionOp interface {
	isZionOp()
}

type TransferOp struct {
	transaction.TransferTx
	ToId           int
	ConditionType  int
	MaxFee         *big.Int
	PutAddressInDa bool
}

func (op TransferOp) isZionOp() {}

type WithdrawOp struct {
	transaction.WithdrawTx
	MaxFee        *big.Int
	ConditionType uint
}

func (op WithdrawOp) isZionOp() {}

type DepositOp struct {
	transaction.DepositTx
	AccountId int
}

func (op DepositOp) isZionOp() {
}

type PubkeyUpdateOp struct {
	transaction.PubkeyUpdateTx
	ConditionType uint
	MaxFee        *big.Int
}

func (op PubkeyUpdateOp) isZionOp() {}

type NoopOp struct{}

func (op NoopOp) IsZionOp() {}
