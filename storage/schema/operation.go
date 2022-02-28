package schema

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CollTx    = "executed_transactions"
	CollPriTx = "executed_priority_transactions"
)

// OperationSchema is capable of storing and loading transactions.
// Every non-executed,executed tx and priority tx can be saved or loaded from this schema.
type OperationSchema struct {
	StorageCore Processor
}

func (o OperationSchema) StoreExecutedTx(operation StoredExecutedTx) (err error) {
	// mempoolSchema := MempoolSchema{StorageCore: o.StorageCore}
	// mempoolSchema.RemoveTx(operation.TxHash)
	if operation.Success {
		// If transaction succeed, it should replace the stored tx with the same hash.
		// This may happen only when has failed previously.
		opts := options.Replace().SetUpsert(true)
		filter := bson.D{{"tx_hash", operation.TxHash}}
		coll := o.StorageCore.AccessCollection(CollTx)
		_, err = coll.ReplaceOne(context.TODO(), filter, operation, opts)
		if err != nil {
			return
		}
	}
	// If transaction failed,just insert the doc,if exists,ignore it.
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"tx_hash", operation.TxHash}}
	coll := o.StorageCore.AccessCollection(CollTx)
	_, err = coll.UpdateOne(context.TODO(), filter, bson.M{"$setOnInsert": operation}, opts)
	return
}

func (o OperationSchema) StoreExecutedPriorityTx(operation StoredExecutedPriTx) (err error) {
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"tx_hash", operation.TxHash}}
	coll := o.StorageCore.AccessCollection(CollPriTx)
	_, err = coll.UpdateOne(context.TODO(), filter, bson.M{"$setOnInsert": operation}, opts)
	return
}

// GetExecutedPriorityTransaction retrieves priority transactions given its ID.
func (o OperationSchema) GetExecutedPriorityTransaction(
	priorityTxId int,
) (StoredExecutedPriTx, error) {
	coll := o.StorageCore.AccessCollection(CollPriTx)
	var op StoredExecutedPriTx
	err := coll.FindOne(context.TODO(), bson.D{{"priority_op_serialid", priorityTxId}}).Decode(&op)
	return op, err
}

// GetExecutedPriorityTxByL1Hash retrieves priority transaction by its L1 hash(transaction hash on cortex mainnet).
func (o OperationSchema) GetExecutedPriorityTxByL1Hash(l1Hash common.Hash) (
	StoredExecutedPriTx, error,
) {
	coll := o.StorageCore.AccessCollection(CollPriTx)
	var op StoredExecutedPriTx
	err := coll.FindOne(context.TODO(), bson.D{{"l1_hash", l1Hash}}).Decode(&op)
	return op, err
}

//
// RECORD
//

// StoredExecutedTx some record of the operation
