package transaction

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vivijj/zigo/crypto/babyjub"
	"github.com/vivijj/zigo/crypto/poseidon"
)

var (
	TransferPoseidonParam = poseidon.NewParams(9, 6, 53)
)

type TransferTx struct {
	AccountId  int
	Nonce      int
	ValidUntil int
	FeeToken   int
	Fee        *big.Int
	From       common.Address
	To         common.Address
	Token      int
	Amount     *big.Int
	// Transaction zion(L2) signature
	Signature babyjub.Signature
}

func (tx TransferTx) isZionTx() {}

func (tx TransferTx) GetBytes() (out []byte) {
	out = append(out, []byte(Transfer)...)
	out = append(out, IntToBytes(tx.AccountId)...)
	out = append(out, IntToBytes(tx.Nonce)...)
	out = append(out, IntToBytes(tx.ValidUntil)...)
	out = append(out, IntToBytes(tx.FeeToken)...)
	out = append(out, tx.Fee.Bytes()...)
	out = append(out, tx.From.Bytes()...)
	out = append(out, tx.To.Bytes()...)
	out = append(out, IntToBytes(tx.Token)...)
	out = append(out, tx.Amount.Bytes()...)

	return
}

// EncodeBi Encode the transaction data as *big.Int by poseidon hash
func (tx TransferTx) EncodeBi() *big.Int {
	var out []*big.Int
	out = append(out, big.NewInt(int64(tx.AccountId)))
	out = append(out, new(big.Int).SetBytes(tx.To.Bytes()))
	out = append(out, big.NewInt(int64(tx.Token)))
	out = append(out, tx.Amount)
	out = append(out, big.NewInt(int64(tx.FeeToken)))
	out = append(out, tx.Fee)
	out = append(out, big.NewInt(int64(tx.ValidUntil)))
	out = append(out, big.NewInt(int64(tx.Nonce)))
	return poseidon.Hash(out, TransferPoseidonParam)
}
