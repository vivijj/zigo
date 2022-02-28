package core

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vivijj/zigo/pkg/account"
	"github.com/vivijj/zigo/pkg/block"
	"github.com/vivijj/zigo/pkg/state"
	"github.com/vivijj/zigo/pkg/transaction"
	"github.com/vivijj/zigo/types/cmp"
)

type StateKeeperRequest interface {
	isStateKeeperRequest()
}

type PendingBlock struct {
	chunksLeft             int
	pendingOpBlockIndex    int
	unprocessedPriTxBefore int
	pendingBlockIteration  int
	// Number of stored account updates in the db
	storedAccountUpdates  int
	successOperation      []block.ExecutedOperation
	failedTxs             []block.ExecutedOperation
	accountUpdates        []account.AccUpdateTuple
	collectedFees         []state.CollectedFee
	previousBlockRootHash common.Hash
	timestamp             int64
}

func NewPendingBlock(
	unprocessedPriorTxBefore int,
	availableSize []int,
	previousBlockRootHash common.Hash,
	timestamp int64,
) PendingBlock {
	chunksLeft := cmp.MaxInSeq(availableSize)

	return PendingBlock{
		chunksLeft:             chunksLeft,
		pendingOpBlockIndex:    0,
		unprocessedPriTxBefore: unprocessedPriorTxBefore,
		pendingBlockIteration:  0,
		storedAccountUpdates:   0,
		successOperation:       []block.ExecutedOperation{},
		failedTxs:              []block.ExecutedOperation{},
		accountUpdates:         []account.AccUpdateTuple{},
		collectedFees:          []state.CollectedFee{},
		previousBlockRootHash:  previousBlockRootHash,
		timestamp:              timestamp,
	}
}

// StateKeeper is the critical part that responsible for tx processing and block forming
// state keeper has no access to the storage.
type StateKeeper struct {
	maxMiniblockIterations  int
	operatorId              int
	currentUnprocessedPriTx int
	successTxsPendingLen    int
	failedTxsPendingLen     int
	availableBlockSizes     []int
	rxForBlocks             <-chan StateKeeperRequest
	txForCommitment         chan<- CommitRequest
	state                   state.State
	pendingBlock            PendingBlock
}

func NewStateKeeper(
	initState state.InitParam,
	operatorAddr common.Address,
	rxForBlocks <-chan StateKeeperRequest,
	txForCommitments chan<- CommitRequest,
	availableBlockSize []int,
	maxMiniblockIteration int,
) *StateKeeper {
	state := state.New(initState.Tree, initState.AccIdByAddr, initState.LastBlockNumber+1)
	operatorId, _, ok := state.GetAccountByAddress(operatorAddr)
	if !ok {
		panic("operator account should be present in the account tree")
	}
	previousRootHash := common.HexToHash(state.RootHash())
	keeper := StateKeeper{
		maxMiniblockIterations:  maxMiniblockIteration,
		operatorId:              operatorId,
		currentUnprocessedPriTx: initState.UnprocessedPriorTx,
		successTxsPendingLen:    0,
		failedTxsPendingLen:     0,
		availableBlockSizes:     availableBlockSize,
		rxForBlocks:             rxForBlocks,
		txForCommitment:         txForCommitments,
		state:                   *state,
		pendingBlock: NewPendingBlock(
			initState.UnprocessedPriorTx, availableBlockSize,
			previousRootHash, time.Now().Unix(),
		),
	}
	return &keeper
}

// return false when no enough space in block
func (sk *StateKeeper) applyPriorityTx(tx transaction.PriorityTx) {
	ops := sk.state.ExecutePriorityTx(tx)
	sk.pendingBlock.chunksLeft -= 1
	sk.pendingBlock.accountUpdates = append(sk.pendingBlock.accountUpdates, ops.Updates...)
	if ops.Fee != nil {
		sk.pendingBlock.collectedFees = append(sk.pendingBlock.collectedFees, *ops.Fee)
	}
	blockIndex := sk.pendingBlock.pendingOpBlockIndex
	sk.pendingBlock.pendingOpBlockIndex += 1
	execResult := block.ExecutedPriorityTx{
		PriTx:      tx,
		Op:         ops.ExecutedTx,
		BlockIndex: blockIndex,
		CreatedAt:  time.Now().Unix(),
	}
	sk.pendingBlock.successOperation = append(sk.pendingBlock.successOperation, execResult)
	sk.currentUnprocessedPriTx += 1
}

// if apply tx fail, return false, otherwise true
func (sk *StateKeeper) applyTx(tx transaction.ZionTx) {
	ops, err := sk.state.ExecuteTx(tx)
	if err != nil {
		failedTX := block.ExecutedTx{
			Tx:         tx,
			Success:    false,
			Op:         nil,
			FailReason: err.Error(),
			BlockIndex: 0,
			CreatedAt:  time.Now().Unix(),
		}
		sk.pendingBlock.failedTxs = append(sk.pendingBlock.failedTxs, failedTX)
	}

	sk.pendingBlock.chunksLeft -= 1
	sk.pendingBlock.accountUpdates = append(sk.pendingBlock.accountUpdates)
	if ops.Fee != nil {
		sk.pendingBlock.collectedFees = append(sk.pendingBlock.collectedFees, *ops.Fee)
	}
	blockIndex := sk.pendingBlock.pendingOpBlockIndex
	sk.pendingBlock.pendingOpBlockIndex += 1

	execResult := block.ExecutedTx{
		Tx:         tx,
		Success:    true,
		Op:         ops.ExecutedTx,
		FailReason: "",
		BlockIndex: blockIndex,
		CreatedAt:  time.Now().Unix(),
	}
	sk.pendingBlock.successOperation = append(sk.pendingBlock.successOperation, execResult)
}

func (sk *StateKeeper) executeProposedBlock(proposedBlock ProposedBlock) {
	if len(sk.pendingBlock.successOperation) == 0 {
		sk.pendingBlock.timestamp = time.Now().Unix()
	}
	isEmptyProposedBlock := proposedBlock.IsEmpty()

	numPritx := len(proposedBlock.PriorityTxs)
	for i := 0; i < numPritx; i++ {
		if sk.pendingBlock.chunksLeft < 1 {
			// when no enough space in block, we should seal old one.
			sk.sealPendingBlock()
			// after seal the old block, we should deal with this priority tx again.
			i--
		}
		sk.applyPriorityTx(proposedBlock.PriorityTxs[i])
	}

	numTx := len(proposedBlock.Txs)
	for i := 0; i < numTx; i++ {
		if sk.pendingBlock.chunksLeft < 1 {
			// when no enough space in block, we should seal old one.
			sk.sealPendingBlock()
			// after seal the old block, we should deal with this priority tx again.
			i--
		}
		sk.applyTx(proposedBlock.Txs[i])
	}

	if len(sk.pendingBlock.successOperation) != 0 {
		sk.pendingBlock.pendingBlockIteration += 1
	}

	if sk.pendingBlock.chunksLeft == 0 || sk.pendingBlock.pendingBlockIteration > sk.maxMiniblockIterations {
		sk.sealPendingBlock()
	} else {
		if !isEmptyProposedBlock {
			sk.storePendingBlock()
		}
	}
}

// Finalizes the pending block,transforming it into a full block.
func (sk *StateKeeper) sealPendingBlock() {
	feeUpdates := sk.state.CollectFee(sk.pendingBlock.collectedFees, sk.operatorId)
	sk.pendingBlock.accountUpdates = append(sk.pendingBlock.accountUpdates, feeUpdates...)
	pendingBlock := sk.pendingBlock
	sk.pendingBlock = PendingBlock{
		unprocessedPriTxBefore: sk.currentUnprocessedPriTx,
		chunksLeft:             sk.blockSize,
		timestamp:              time.Now().Unix(),
	}
	sk.successTxsPendingLen = 0
	sk.failedTxsPendingLen = 0

	blockTransactions := append(pendingBlock.successOperation, pendingBlock.failedTxs...)

	fullBlock := block.Block{
		BlockNumber:          sk.state.BlockNumber,
		NewRootHash:          sk.state.RootHash(),
		Operator:             sk.operatorId,
		BlockTransactions:    blockTransactions,
		ProcessedPriTxBefore: pendingBlock.unprocessedPriTxBefore,
		ProcessedPriTxAfter:  sk.currentUnprocessedPriTx,
		BlockSize:            sk.blockSize,
		TimeStamp:            0,
	}
	firstUpdateId := pendingBlock.storedAccountUpdates
	accUpds := pendingBlock.accountUpdates[firstUpdateId:]
	blockCommitReq := CommitBlockRequest{
		Block:          fullBlock,
		FirstUpdateId:  firstUpdateId,
		AccountUpdates: accUpds,
	}
	pendingBlock.storedAccountUpdates = len(pendingBlock.accountUpdates)

	sk.state.BlockNumber += 1
	sk.txForCommitment <- blockCommitReq
}

// store intermediate representation of a pending block in database.
// so the executed transaction are persisted
func (sk *StateKeeper) storePendingBlock() {
	newSuccessOps := CloneSlice(sk.pendingBlock.successOperation[sk.successTxsPendingLen:])
	newFailedOps := CloneSlice(sk.pendingBlock.failedTxs[sk.failedTxsPendingLen:])

	sk.successTxsPendingLen = len(sk.pendingBlock.successOperation)
	sk.failedTxsPendingLen = len(sk.pendingBlock.failedTxs)

	pendingBlock := block.InterPendingBlock{
		BlockNumber:            sk.state.BlockNumber,
		ChunksLeft:             sk.pendingBlock.chunksLeft,
		UnprocessedPriTxBefore: sk.pendingBlock.unprocessedPriTxBefore,
		PendingBlockIteration:  sk.pendingBlock.pendingBlockIteration,
		SuccessOperations:      newSuccessOps,
		FailedTxs:              newFailedOps,
		PreviousBlockRootHash:  sk.pendingBlock.previousBlockRootHash,
		TimeStamp:              sk.pendingBlock.timestamp,
	}
	firstUpdateId := sk.pendingBlock.storedAccountUpdates

	// copy the account updates because we will transfer it to committer by chan.
	accUpds := CloneSlice(sk.pendingBlock.accountUpdates)
	sk.pendingBlock.storedAccountUpdates = len(sk.pendingBlock.accountUpdates)
	commitReq := CommitPendingBlockRequest{
		PendingBlock:   pendingBlock,
		FirstUpdateId:  firstUpdateId,
		AccountUpdates: accUpds,
	}
	sk.txForCommitment <- commitReq
}
