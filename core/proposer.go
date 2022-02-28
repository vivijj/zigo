package core

import (
	"time"

	"github.com/vivijj/zigo/internal/config"
)

// block proposer is main driver of the application, it polls transactions from mempool and send
// them to state keeper in small batch call mini block.

type BlockProposer struct {
	curPriTxNumber  int
	mempoolReqs     chan<- MempoolBlocksRequest
	statekeeperReqs chan<- StateKeeperRequest
}

func (bp *BlockProposer) proposeNewBlock() ProposedBlock {
	respChan := make(chan ProposedBlock)
	req := MempoolBlocksRequest{
		LastPriTxNumber: bp.curPriTxNumber,
		ResponseSender:  respChan,
	}
	bp.mempoolReqs <- req
	res := <-respChan
	return res
}

func (bp *BlockProposer) commitNewTxMiniBatch() {
	proposedBlock := bp.proposeNewBlock()
	bp.curPriTxNumber += len(proposedBlock.PriorityTxs)
	bp.statekeeperReqs <- ExecuteMiniBlockReq{proposedBlock: proposedBlock}

}

func RunBlockProposerTask(
	conf config.Config,
	mempoolRequests chan<- MempoolBlocksRequest,
	statekeeperRequests chan<- StateKeeperRequest,
) {
	miniblockInterval := conf.Core.MiniBlockInterations
	lastUnprocessedPriTxChan := make(chan int)
	statekeeperRequests <- GetLastUnprocessedPriTxReq{resp: lastUnprocessedPriTxChan}

	curPriTxNumber := <-lastUnprocessedPriTxChan

	blockProposer := BlockProposer{
		curPriTxNumber:  curPriTxNumber,
		mempoolReqs:     mempoolRequests,
		statekeeperReqs: statekeeperRequests,
	}

	t := time.NewTicker(time.Duration(miniblockInterval))
	for {
		<-t.C
		blockProposer.commitNewTxMiniBatch()
	}
}
