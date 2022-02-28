package record

import (
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vivijj/zigo/pkg/block"
	"github.com/vivijj/zigo/pkg/transaction"
)

type StoredExecutedTx struct {
	BlockNumber           int             `bson:"block_number"`
	BlockIndex            int             `bson:"BlockIndex"`
	Tx                    json.RawMessage `bson:"tx"`
	Operation             json.RawMessage `bson:"operation"`
	TxHash                common.Hash     `bson:"tx_hash"`
	FromAccount           common.Address  `bson:"from_account"`
	ToAccount             common.Address  `bson:"to_account"`
	Success               bool            `bson:"success"`
	FailReason            string          `bson:"fail_reason"`
	PrimaryAccountAddress common.Address  `bson:"primary_account_address"`
	Nonce                 int             `bson:"nonce"`
	CreatedAt             time.Time       `bson:"created_at"`
}

func PrepareStoredTx(execTx block.ExecutedTx, blockNumber int) StoredExecutedTx {
	jtx, err := json.Marshal(execTx.Tx)
	if err != nil {
		panic("can't serialize tx")
	}
	jop, err := json.Marshal(execTx.Op)
	if err != nil {
		panic("can't serialize operation")
	}

	return StoredExecutedTx{
		BlockNumber: blockNumber,
		TxHash:      transaction.ZionTxHash(execTx.Tx),
		Tx:          jtx,
		Operation:   jop,
		Success:     execTx.Success,
	}
}

type StoredExecutedPriTx struct {
	BlockNumber        int             `bson:"block_number"`
	BlockIndex         int             `bson:"block_index"`
	Operation          json.RawMessage `bson:"operation"`
	FromAccount        common.Address  `bson:"from_account"`
	ToAccount          common.Address  `bson:"to_account"`
	PriorityOpSerialid int             `bson:"priority_op_serialid"`
	L1Hash             common.Hash     `bson:"l1_hash"`
	L1Block            int             `bson:"l1_block"`
	L1BlockIndex       int             `bson:"l1_block_index"`
	TxHash             common.Hash     `bson:"tx_hash"`
	CreatedAt          time.Time       `bson:"created_at"`
}

func PrepareStoredPriorityTx(
	execPriTx block.ExecutedPriorityTx,
	blockNumber int,
) StoredExecutedPriTx {
	jop, err := json.Marshal(execPriTx.Op)
	if err != nil {
		panic("can't serialize priority tx")
	}
	if deposit, ok := execPriTx.PriTx.Data.(transaction.DepositTx); ok {
		return StoredExecutedPriTx{
			BlockNumber:        blockNumber,
			BlockIndex:         execPriTx.BlockIndex,
			Operation:          jop,
			FromAccount:        deposit.From,
			ToAccount:          deposit.To,
			PriorityOpSerialid: int(execPriTx.PriTx.SerialId),
			L1Hash:             execPriTx.PriTx.L1Hash,
			L1Block:            int(execPriTx.PriTx.L1Block),
			L1BlockIndex:       int(execPriTx.PriTx.L1BlockIndex),
			TxHash:             transaction.PriTxHash(execPriTx.PriTx),
			CreatedAt:          time.Unix(int64(execPriTx.CreatedAt), 0),
		}
	}
	panic("incorrect type of priority tx")
}
