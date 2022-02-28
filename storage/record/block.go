package record

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

type StorageBlock struct {
	BlockNumber            int         `bson:"block_number"`
	RootHash               common.Hash `bson:"root_hash"`
	OperatorId             int         `bson:"operator_id"`
	UnprocessedPriTxBefore int         `bson:"unprocessed_pri_tx_before"`
	UnprocessedPriTxAfter  int         `bson:"unprocessed_pri_tx_after"`
	BlockSize              int         `bson:"block_size"`
	Commitment             common.Hash `bson:"commitment"`
	Timestamp              int         `bson:"timestamp"`
}

type StoragePendingBlock struct {
	BlockNumber                 int
	ChunksLeft                  int
	UnprocessedPriorityTxBefore int
	PendingBlockIteration       int
	PreviousRootHash            common.Hash
	TimeStamp                   int64
}

type AccTreeCache struct {
	Block     int             `bson:"block"`
	TreeCache json.RawMessage `bson:"tree_cache"`
}
