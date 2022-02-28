package core

import "github.com/vivijj/zigo/pkg/transaction"

type GetPriorityQueueTxs struct {
	txStartId int
	maxChunk  int
	resp      chan<- []transaction.PriorityTx
}
