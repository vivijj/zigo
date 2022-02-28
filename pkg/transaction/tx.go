// Packeage define the tx in the zion network(tx from l2 directly & priority tx from contract)
package transaction

import (
	"crypto/sha256"
	"strconv"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vivijj/zigo/internal/zcontract"
)

type TransactionType string

const (
	Noop         TransactionType = "Noop"
	Deposit      TransactionType = "Deposit"
	Withdraw     TransactionType = "Withdraw"
	Transfer     TransactionType = "Transfer"
	PubKeyUpdate TransactionType = "PubKeyUpdate"
)

func ZionTxHash(tx ZionTx) common.Hash {
	txBytes := tx.GetBytes()
	return sha256.Sum256(txBytes)
}

func PriTxHash(tx PriorityTx) common.Hash {
	return sha256.Sum256(tx.GetBytes())
}

// ZionTx is the L2 transaction(transfer, pubkeyupdate, withdraw) init from user directly.
type ZionTx interface {
	// flag function to indicate that the type is instance of ZionTx with no impl.
	isZionTx()
	// // GetBytes Encode the transaction data as the byte sequence
	GetBytes() []byte
}

// priority tx(deposit) is the transactin init from contract,due to it need confirmation,so when we receieve
// it, we shoule process it "priority".
type ZionPriTx interface {
	// flag function
	isZionPriTx()
}

// Priority transaction description with the metadata required for server to process it.
type PriorityTx struct {
	// Unique ID of this priority transaction.
	SerialId uint64
	// Hash of corresponding l1 transaction.
	L1Hash common.Hash
	// Block in which L1 transaction was included.
	L1Block uint64
	// Transaction index in the L1 Block
	L1BlockIndex uint
	// Priority transaction(only deposit now)
	Data ZionPriTx
}

func (ptx PriorityTx) isZionPriTx() {}

// all the ZionPriorityTx should use the L1 info to identify itself
func (ptx PriorityTx) GetBytes() (out []byte) {
	out = append(out, ptx.L1Hash.Bytes()...)
	out = append(out, IntToBytes(ptx.L1Block)...)
	out = append(out, IntToBytes(ptx.L1BlockIndex)...)
	return
}

func ParsePriorityTxFromLog(logData []byte) PriorityTx {
	eventData := zcontract.ParseNewDepositRequest(logData)

	deposit := DepositTx{
		From:   eventData.Sender,
		To:     eventData.Receiver,
		Amount: eventData.Amount,
		Token:  eventData.TokenId,
	}

	return PriorityTx{
		SerialId:     eventData.PriorityReqId,
		L1Hash:       eventData.Raw.TxHash,
		L1Block:      eventData.Raw.BlockNumber,
		L1BlockIndex: eventData.Raw.Index,
		Data:         deposit,
	}
}

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func IntToBytes[T Integer](i T) []byte {
	return []byte(strconv.Itoa(int(i)))
}
