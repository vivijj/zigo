package state

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vivijj/zigo/internal/param"
	"github.com/vivijj/zigo/internal/smt"
	"github.com/vivijj/zigo/pkg/account"
	"github.com/vivijj/zigo/pkg/block"
	"github.com/vivijj/zigo/storage"
)

// ************* State Init Param ****************

// InitParam is used to init the zion state,we put it as separate struct due to the logic of
// init is a little tricky
type InitParam struct {
	Tree               smt.SparseMerkleTree[account.Account]
	AccIdByAddr        map[common.Address]int
	LastBlockNumber    int
	UnprocessedPriorTx int
}

func NewInitParam() *InitParam {
	return &InitParam{
		Tree:               smt.NewSparseMerkleTree[account.Account](param.QuadAccountTreeDepth),
		AccIdByAddr:        make(map[common.Address]int),
		LastBlockNumber:    0,
		UnprocessedPriorTx: 0,
	}
}

func (p *InitParam) GetPendingBlock(sp *storage.Processor) *block.InterPendingBlock {
	pendingBlock, err := sp.BlockSchema().LoadPendingBlock()
	if err != nil {
		return nil
	}
	if pendingBlock.BlockNumber <= p.LastBlockNumber {
		return nil
	}

	// pending block must be greater than the last committed block exactly by 1.
	if pendingBlock.BlockNumber != p.LastBlockNumber+1 {
		panic("does not match pending block number")
	}
	return pendingBlock
}

// RestoreFromDB will build the state init param
func RestoreFromDB(sp *storage.Processor) (*InitParam, error) {
	initParam := NewInitParam()
	if err := initParam.loadFromDb(sp); err != nil {
		return nil, err
	}
	return initParam, nil
}

func (p *InitParam) loadFromDb(sp *storage.Processor) error {
	blockNumber, err := p.loadAccountTree(sp)
	if err != nil {
		return err
	}
	p.LastBlockNumber = blockNumber
	if p.UnprocessedPriorTx, err = unprocessedPriorityTxId(sp, blockNumber); err != nil {
		return err
	}
	return nil
}

func (p *InitParam) loadAccountTree(sp *storage.Processor) (int, error) {
	var lastCachedBlockNumber int
	var accs map[int]account.Account

	blkNum, _, ok, err := sp.BlockSchema().GetAccTreeCache()
	if err != nil {
		return 0, err
	}
	if ok {
		lastCachedBlockNumber, accs, err = sp.StateSchema().LoadCommittedState(blkNum)
	} else {
		lastCachedBlockNumber, accs, err = sp.StateSchema().LoadVerifiedState()
	}
	if err != nil {
		return 0, err
	}

	for i := range accs {
		p.InsertAccount(i, accs[i])
	}

	treeCache, ok, err := sp.BlockSchema().GetAccTreeCacheBlock(lastCachedBlockNumber)
	if err != nil {
		return 0, err
	}
	if ok {
		if err = json.Unmarshal(treeCache, &(p.Tree)); err != nil {
			return 0, err
		}
	} else {
		treeCache, _ = json.Marshal(p.Tree)
		sp.BlockSchema().StoreAccTreeCache(lastCachedBlockNumber, treeCache)
	}

	blockNumber, accs, err := sp.StateSchema().LoadCommittedState(-1)
	if err != nil {
		return 0, fmt.Errorf("couldn't load committed state: %w", err)
	}
	if blockNumber != lastCachedBlockNumber {

	}

}

func (p *InitParam) InsertAccount(accId int, acc account.Account) {
	p.AccIdByAddr[acc.Address] = accId
	p.Tree.Insert(accId, acc)
}
func unprocessedPriorityTxId(sp *storage.Processor, blockNumber int) (int, error) {
	blk, err := sp.BlockSchema().GetBlock(blockNumber)
	if err != nil {
		return -1, err
	}
	if blk == -1 {
		// there is no specific block yet, so 0 is return.
		return 0, nil
	}
	return blk, nil
}
