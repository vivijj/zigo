// Package block include some zion network block definitions
package block

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/vivijj/zigo/pkg/operation"
	"github.com/vivijj/zigo/pkg/transaction"
)

type InterPendingBlock struct {
	BlockNumber            int
	ChunksLeft             int
	UnprocessedPriTxBefore int
	// Amount of processing iterations applied to the pending block.
	// block will be sealed even if it's not full when this amount exceeds limit.
	PendingBlockIteration int
	SuccessOperations     []ExecutedOperation
	FailedTxs             []ExecutedOperation
	PreviousBlockRootHash common.Hash
	TimeStamp             int64
}

// Block zion network block
type Block struct {
	BlockNumber int
	// state of chain root hash after execute this block
	NewRootHash string
	// ID of zion network operator
	Operator int
	// List of operation executed in the block(L1 & L2).
	BlockTransactions    []ExecutedOperation
	ProcessedPriTxBefore int
	ProcessedPriTxAfter  int
	BlockSize            int
	TimeStamp            int
}

// ExecutedTx Executed L2 transactions.
type ExecutedTx struct {
	Tx      transaction.ZionTx
	Success bool
	// if fail, Op is nil
	Op         operation.ZionOp
	FailReason string
	BlockIndex int
	CreatedAt  int64
}

func (_ ExecutedTx) isExecutedOperation() {}

// ExecutedPriorityTx L1 priority transactions, can't fail in L2.
type ExecutedPriorityTx struct {
	PriTx      transaction.PriorityTx
	Op         operation.ZionOp
	BlockIndex int
	CreatedAt  int64
}

func (_ ExecutedPriorityTx) isExecutedOperation() {}

// ExecutedOperation Representation of executed operations, which can be either L1 or L2.
type ExecutedOperation interface {
	isExecutedOperation()
}
