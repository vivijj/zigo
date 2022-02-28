package core

import (
	"fmt"
	"log"

	"github.com/vivijj/zigo/internal/storage"

	"github.com/vivijj/zigo/internal/config"
	"github.com/vivijj/zigo/pkg/transaction"
	"github.com/vivijj/zigo/types/deque"
)

// TODO
type L1WatchRequest interface{}

// mempool is memory buffer for transactions:
// 1. it accept transactions from api,
// 2. when polled return vector of transactions in the queue.
//
// mempool is not persisted, all transactions will be lost on node shutdown.
// so when restart the node, we should restore the mempool from db.

type ProposedBlock struct {
	PriorityTxs []transaction.PriorityTx
	Txs         []transaction.ZionTx
}

func (pb ProposedBlock) IsEmpty() bool {
	return (len(pb.PriorityTxs) == 0) && (len(pb.Txs) == 0)
}

type MempoolBlocksRequest struct {
	LastPriTxNumber int
	ResponseSender  chan<- ProposedBlock
}

// MempoolTransactionRequest ad new tx to mempool.
// ResponseSender is used to receive tx add result.
type MempoolTransactionRequest struct {
	Tx             transaction.ZionTx
	ResponseSender chan<- error
}

type MempoolState struct {
	txQueue deque.Deque[transaction.ZionTx]
}

func restoreStateFromDb(conn storage.ConnectionPool) *MempoolState {
	mempoolSchema := conn.AccessStorage().GetMempoolSchema()
	txs, err := mempoolSchema.LoadTxs()
	if err != nil {
		log.Fatal("Attempt tp restore mempool txs from DB failed: ", err)
	}

	txq := deque.New[transaction.ZionTx]()
	for i := range txs {
		txq.PushBack(txs[i])
	}
	return &MempoolState{txQueue: *txq}

}

func (st *MempoolState) addTx(tx transaction.ZionTx) {
	st.txQueue.PushBack(tx)
}

// MempoolHandler deal with 2 kind of requests: BlocksRequest and TransactionRequest
type MempoolHandler struct {
	conn         storage.ConnectionPool
	mempoolState MempoolState
	blockReq     <-chan MempoolBlocksRequest
	txReq        <-chan MempoolTransactionRequest
	l1WatchReq   chan<- L1WatchRequest
	maxBlockSize int
}

// handle the block

func (mh *MempoolHandler) proposeNewBlock(currentUnprocessedPriTx int) ProposedBlock {
	chunks_left, priorityTx := mh.selectPriorityTxs(currentUnprocessedPriTx)
	txs := mh.prepareTxForBlock(chunks_left)
	return ProposedBlock{
		PriorityTxs: priorityTx,
		Txs:         txs,
	}
}

// return the chunks left and selected priority txs
func (mh *MempoolHandler) selectPriorityTxs(curUnprocessedPriTx int) (
	int, []transaction.PriorityTx,
) {
	respChan := make(chan []transaction.PriorityTx)
	mh.l1WatchReq <- GetPriorityQueueTxs{
		txStartId: curUnprocessedPriTx,
		maxChunk:  mh.maxBlockSize,
		resp:      respChan,
	}
	// wait until received response which contains priority txs
	priorityTxs := <-respChan
	return mh.maxBlockSize - len(priorityTxs), priorityTxs
}

func (mh *MempoolHandler) prepareTxForBlock(chunksLeft int) (txs []transaction.ZionTx) {
	for mh.mempoolState.txQueue.Len() > 0 {
		if chunksLeft <= 0 {
			break
		}
		tx := mh.mempoolState.txQueue.PopFront()
		txs = append(txs, tx)
		chunksLeft -= 1
	}
	return
}

// handle the transaction
func (mh *MempoolHandler) addTx(tx transaction.ZionTx) error {
	mempoolSchema := mh.conn.AccessStorage().GetMempoolSchema()
	err := mempoolSchema.InsertTx(tx)
	if err != nil {
		return fmt.Errorf("Mempool storage access error: %w", err)
	}
	mh.mempoolState.txQueue.PushBack(tx)
	return nil
}

func (mh *MempoolHandler) run() {
	for {
		select {
		// the block request
		case req := <-mh.blockReq:
			proposedBlock := mh.proposeNewBlock(req.LastPriTxNumber)
			req.ResponseSender <- proposedBlock

		// the transaction request
		case req := <-mh.txReq:
			err := mh.addTx(req.Tx)
			req.ResponseSender <- err
		}

	}
}

func RunMempoolTask(
	conn storage.ConnectionPool,
	txRequests <-chan MempoolTransactionRequest,
	blockRequests <-chan MempoolBlocksRequest,
	l1WatchRequests chan<- L1WatchRequest,
	conf config.Config,
) {
	maxBlockSize := conf.Core.BlockChunkSizes
	mempoolState := restoreStateFromDb(conn)
	mempoolHandler := MempoolHandler{
		conn:         conn,
		txReq:        txRequests,
		blockReq:     blockRequests,
		l1WatchReq:   l1WatchRequests,
		maxBlockSize: maxBlockSize,
		mempoolState: *mempoolState,
	}

	// start the goroutine to run the mempool handler
	go mempoolHandler.run()
}
