package schema

import (
	"context"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vivijj/zigo/pkg/block"
	"github.com/vivijj/zigo/storage"
	"github.com/vivijj/zigo/storage/record"
)

const (
	CollBlock        = "blocks"
	CollPending      = "pending_block"
	CollAccTreeCache = "account_tree_cache"
)

// BlockSchema is primary L2 chain storage controller.
type BlockSchema struct {
	Conn storage.Processor
}

// SaveBlockTransactions Given a block, stores its transactions.
func (bs BlockSchema) SaveBlockTransactions(
	blockNumber int,
	operations []block.ExecutedOperation,
) (err error) {
	for i := range operations {
		switch tx := (operations[i]).(type) {
		case block.ExecutedPriorityTx:
			storeTx := PrepareStoredPriorityTx(tx, blockNumber)
			err = bs.StorageCore.OperationSchema().StoreExecutedPriorityTx(storeTx)
		case block.ExecutedTx:
			storeTx := PrepareStoredTx(tx, blockNumber)
			err = bs.StorageCore.OperationSchema().StoreExecutedTx(storeTx)
		}
		if err != nil {
			return
		}
	}
	return
}

// SaveBlock contains 2 topic:
// - save all the transactions in the block
// - save the block as `StorageBlock` and remove associative InterPendingBlock
func (bs BlockSchema) SaveBlock(block block.Block) (err error) {
	err = bs.SaveBlockTransactions(block.BlockNumber, block.BlockTransactions)
	if err != nil {
		return
	}
	newBlock := StoredBlock{
		BlockNumber:            block.BlockNumber,
		RootHash:               block.NewRootHash,
		OperatorId:             block.Operator,
		UnprocessedPriTxBefore: block.ProcessedPriTxBefore,
		UnprocessedPriTxAfter:  block.ProcessedPriTxAfter,
		blockSize:              block.BlockSize,
		commitment:             block.BlockCommitment,
		timestamp:              block.TimeStamp,
	}
	collPending := bs.AccessCollection(CollPending)
	_, err = collPending.DeleteOne(context.TODO(), bson.D{{"block_number", newBlock.BlockNumber}})
	if err != nil {
		return
	}
	collBlock := bs.AccessCollection(CollBlock)
	_, err = collBlock.InsertOne(context.TODO(), newBlock)
	return
}

// GetAccTreeCache get the latest stored account tree cache.
// if not exist, the ErrNoDocuments error will be return.
func (bs BlockSchema) GetAccTreeCache() (
	blockNumber int,
	jvalue json.RawMessage,
	ok bool,
	err error,
) {
	coll := bs.Conn.AccessCollection(CollAccTreeCache)

	var accTreeCache record.AccTreeCache
	opts := options.FindOne().SetSort(bson.D{{"block", -1}})
	err = coll.FindOne(context.TODO(), bson.D{}, opts).Decode(&accTreeCache)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// no match doc
			return 0, nil, false, nil
		}
		return
	}
	return accTreeCache.Block, accTreeCache.TreeCache, true, nil
}

// GetAccTreeCacheBlock gets stored account tree cache for a block
func (bs BlockSchema) GetAccTreeCacheBlock(blockNumber int) (json.RawMessage, bool, error) {
	var accTreeCache record.AccTreeCache
	err := bs.Conn.AccessCollection(CollAccTreeCache).FindOne(
		context.TODO(),
		bson.D{
			{
				"block",
				blockNumber,
			},
		},
	).Decode(accTreeCache)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, false, nil
		}
		return nil, false, err
	}
	return accTreeCache.TreeCache, true, nil
}
