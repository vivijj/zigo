package core

import (
	"github.com/vivijj/zigo/internal/config"
	"github.com/vivijj/zigo/internal/storage"
	"github.com/vivijj/zigo/pkg/account"
	"github.com/vivijj/zigo/pkg/block"
)

type CommitRequest interface {
	isCommitRequest()
}
type CommitBlockRequest struct {
	Block          block.Block
	FirstUpdateId  int
	AccountUpdates account.Updates
}

func (_ CommitBlockRequest) isCommitRequest() {}

type CommitPendingBlockRequest struct {
	PendingBlock   block.InterPendingBlock
	FirstUpdateId  int
	AccountUpdates account.Updates
}

func (_ CommitPendingBlockRequest) isCommitRequest() {}

func handleNewCommitTask(
	rxForOps <-chan CommitRequest,
	conn storage.ConnectionPool,
) {
	for {
		req := <-rxForOps
		switch req.(type) {
		case CommitBlockRequest:
			commitBlock()
		case CommitPendingBlockRequest:
			savePendingBlock()
		}
	}

}

func savePendingBlock(
	pendingBlock block.InterPendingBlock,
	conn storage.ConnectionPool,
) {
	storageCore := conn.AccessStorage()
	storageCore.GetBlockSchema().SavePendingBlock(pendingBlock)
	blockNumber := pendingBlock.BlockNumber
	storageCore.GetStateSchema().CommitStateUpdate(
		blockNumber, appliedUpdatesRequest.AccountUpdates,
	)

}

func commitBlock(
	blockCommitReq CommitBlockRequest,
	conn storage.ConnectionPool,
) {
	stateSchema := conn.AccessStorage().GetStateSchema()

	err := stateSchema.CommitStateUpdate(
		blockCommitReq.Block.BlockNumber,
		blockCommitReq.AccountUpdates,
	)
	if err != nil {
		panic("committer must commit the pending block into db")
	}

	blockSchema := conn.AccessStorage().GetBlockSchema()
	err = blockSchema.SaveBlock(blockCommitReq.Block)

}

// RunCommitter start the standalone go routine.
func RunCommitter(
	rxForOps <-chan CommitRequest,
	mempoolReqSender chan<- MempoolBlocksRequest,
	conn storage.ConnectionPool,
	conf config.Config,
) {
	go handleNewCommitTask(rxForOps, mempoolReqSender, conn)
}
