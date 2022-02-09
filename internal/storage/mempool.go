package storage

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/steinselite/zigo/pkg/tx"
)

// Schema for persisting transactions awaiting for the execution. This schema holds the transactions
// that are received by the `mempool`, but not yet have been included into some block.

const CollMempool = "mempool"

type MempoolTxRecord struct {
	TxHash   string        `bson:"tx_hash"`
	Tx       tx.ZionTxJson `bson:"tx"`
	CreateAt int           `bson:"created_at"`
}

// ToZionTx parse the mempoolTxRecord into ZionTx
func (rec *MempoolTxRecord) ToZionTx() tx.ZionTx {
	return rec.Tx.ParseZionTx()
}

type MempoolSchema struct {
	StorageP *StorageProcessor
}

// AccessCollection try to get the the specific schema collection in the database.
func (ms MempoolSchema) AccessCollection() *mongo.Collection {
	return ms.StorageP.conn.Collection(CollMempool)
}

// LoadTxs load all the transactions stored in the mempool schema.
func (ms MempoolSchema) LoadTxs() (txs []tx.ZionTx, err error) {
	coll := ms.AccessCollection()
	opts := options.Find().SetSort(bson.D{{"created_at", 1}})
	cursor, err := coll.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		return
	}
	var queryResults []MempoolTxRecord
	err = cursor.All(context.TODO(), &queryResults)
	if err != nil {
		return
	}
	txs = make([]tx.ZionTx, 0, len(queryResults))
	for i := range queryResults {
		txs = append(txs, queryResults[i].ToZionTx())
	}
	return
}

// InsertTx add a new transaction into the mempool schema
func (ms MempoolSchema) InsertTx(txData tx.ZionTx) error {
	txHash := tx.ZionTxHash(txData)
	txHashStr := hex.EncodeToString(txHash[:])

	jtx := tx.FromZionTxToJson(txData)

	coll := ms.AccessCollection()
	_, err := coll.InsertOne(context.TODO(), MempoolTxRecord{
		TxHash:   txHashStr,
		Tx:       jtx,
		CreateAt: int(time.Now().Unix()),
	})
	return err
}

func (ms MempoolSchema) RemoveTx(txHash []byte) (err error) {
	txHashStr := hex.EncodeToString(txHash)
	coll := ms.AccessCollection()

	filter := bson.D{{"tx_hash", txHashStr}}
	_, err = coll.DeleteOne(context.TODO(), filter)
	return
}

// ContainTx check if memory pool contains transaction with the given hash.
func (ms MempoolSchema) ContainTx(txHash []byte) bool {
	coll := ms.AccessCollection()
	txHashStr := hex.EncodeToString(txHash)
	filter := bson.D{{"tx_hash", txHashStr}}
	cursor, _ := coll.Find(context.TODO(), filter)
	if cursor.Next(context.TODO()) {
		return true
	}
	return false
}

// GetTx return zion transaction with the given hash
// if tx not exist, return ErrNoDocuments
func (ms *MempoolSchema) GetTx(txHash []byte) (tx tx.ZionTx, err error) {
	mempoolTx, err := ms.GetMempoolTx(txHash)
	if err != nil {
		err = fmt.Errorf("error GetTx: %w", err)
		return
	}
	tx = mempoolTx.ToZionTx()
	return
}

// GetMempoolTx returns mempool transactions as it is stored in the database
func (ms *MempoolSchema) GetMempoolTx(txHash []byte) (mempoolTx MempoolTxRecord, err error) {
	coll := ms.AccessCollection()
	txHashStr := hex.EncodeToString(txHash)
	filter := bson.D{{"tx_hash", txHashStr}}
	err = coll.FindOne(context.TODO(), filter).Decode(&mempoolTx)
	if err != nil {
		err = fmt.Errorf("fail GetMempoolTx: %w", err)
	}
	return
}
