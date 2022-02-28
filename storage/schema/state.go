package schema

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vivijj/zigo/pkg/account"
	"github.com/vivijj/zigo/storage"
	"github.com/vivijj/zigo/storage/record"
)

// State schema is capable of managing the state of the zion "chain" state:
// - Account management: applying diff to the account map
//
// Saving state is done in 2 steps:
// 1. When block is committed,save all state updates.
// 2. Once block is verified, apply this updates to stored state snapshot.

// some collection StateSchema include:
// account_creates: include account create
// account_balance_updates: the update balance of specific account
// account_pubkey_updates: the update pubkey of specific account
const (
	CollAcc     = "accounts"
	CollCreate  = "account_create"
	CollBalance = "account_balance_update"
	CollPubkey  = "account_pubkey_update"
)

type StateSchema struct {
	StorageCore storage.Processor
}

// CommitStateUpdate Stores the list of updates to the account map in the database.
// At this step, the changes are not verified yet, and thus are not applied.
func (ss StateSchema) CommitStateUpdate(
	blockNumber int,
	accountsUpdated account.Updates,
) (err error) {
	for i := range accountsUpdated {
		upd := accountsUpdated[i].AccUpdate

		switch upd := upd.(type) {
		case account.Create:
			coll := ss.AccessCollection(CollCreate)
			_, err = coll.InsertOne(
				context.TODO(),
				bson.D{
					{"account_id", accountsUpdated[i].AccId},
					{"is_create", true},
					{"block_number", blockNumber},
					{"address", upd.Address},
					{"nonce", upd.Nonce},
				},
			)

		case account.UpdateBalance:
			coll := ss.AccessCollection(CollBalance)
			_, err = coll.InsertOne(
				context.TODO(),
				bson.D{
					{"account_id", accountsUpdated[i].AccId},
					{"block_number", blockNumber},
					{"token_id", upd.TokenId},
					{"old_balance", upd.OldBalance},
					{"new_balance", upd.NewBalance},
					{"old_nonce", upd.OldNonce},
					{"new_nonce", upd.NewNonce},
				},
			)

		case account.PubKeyUpdate:
			coll := ss.AccessCollection(CollPubkey)
			_, err = coll.InsertOne(
				context.TODO(),
				bson.D{
					{"account_id", accountsUpdated[i].AccId},
					{"block_number", blockNumber},
					{"old_pubkey", upd.OldPubkeyPair},
					{"new_pubkey", upd.NewPubkeyPair},
					{"old_nonce", upd.OldNonce},
					{"new_nonce", upd.NewNonce},
				},
			)
		}
	}
	return
}

// LoadCommittedState Loads the committed (not necessarily verified) account map state along with a
// block number to which this state applies. If the provided block number is -1, then the latest
// committed state will be loaded.
func (ss StateSchema) LoadCommittedState(block int) (
	blockNumber int,
	accmap map[int]account.Account,
	err error,
) {
	verifyBlock, accs, err := ss.LoadVerifiedState()
	if err != nil {
		return
	}

	stateDiff, err := ss.LoadStateDiff(verifyBlock, block)
	if err != nil {
		return
	}
	if stateDiff != nil {

	}

}

// LoadVerifiedState Loads the verified account map state along with a block number to which
// this state applies.
func (ss StateSchema) LoadVerifiedState() (
	lastBlock int, accmap map[int]account.Account, err error,
) {
	coll := ss.StorageCore.AccessCollection(CollAcc)
	lastBlock, err = ss.StorageCore.BlockSchema().GetLastVerifiedConfirmedBlock()
	cursor, _ := coll.Find(context.TODO(), bson.D{})
	var saccs []record.StorageAccount
	err = cursor.All(context.TODO(), &saccs)
	if err != nil {
		return
	}
	accmap = make(map[int]account.Account)
	for k := range saccs {
		_, acc := record.RestoreAccount(saccs[k])
		accmap[k] = acc
	}
	return
}

// LoadStateDiff Returns the list of updates, and the block number such that if we apply
// these updates to the state of the block #(from_block), we will obtain state of the block
// #(returned block number).
// Returned block number is either `to_block`, the latest committed block before `to_block`.
// If `to_block` is -1, then it will be assumed to be the number of the latest committed
// block.
func (ss StateSchema) LoadStateDiff(fromBlock int, toBlock int) (
	blockNumber int, upds account.Updates, err error,
) {
	if toBlock == -1 {
		coll := ss.StorageCore.AccessCollection("blocks")
		var res record.StorageBlock
		opt := options.FindOne().
			SetSort(bson.D{{"block_number", -1}})
		_ = coll.FindOne(context.TODO(), bson.D{}, opt).Decode(&res)
		toBlock = res.BlockNumber
	}

	filter := bson.D{
		{
			"$and",
			bson.A{
				bson.D{{"block_number", bson.D{{"$gt", fromBlock}}}},
				bson.D{{"block_number", bson.D{{"$lte", toBlock}}}},
			},
		},
	}
	// get the `account_balance_update`
	cursor, _ := ss.StorageCore.AccessCollection(CollBalance).Find(context.TODO(), filter)
	var accBalanceDiffs []record.StorageAccountUpdate
	_ = cursor.All(context.TODO(), &accBalanceDiffs)

	// get the `account_create`
	cursor, _ = ss.StorageCore.AccessCollection(CollCreate).Find(context.TODO(), filter)
	var accCreationDiff []record.StorageAccountCreation
	_ = cursor.All(context.TODO(), &accCreationDiff)

	cursor, _ = ss.StorageCore.AccessCollection(CollPubkey).Find(context.TODO(), filter)
	var accPubkeyDiff []record.StorageAccountPubkeyUpdate
	_ = cursor.All(context.TODO(), &accPubkeyDiff)

	numDiff := len(accBalanceDiffs) + len(accCreationDiff) + len(accPubkeyDiff)
	accDiff := make([]record.StorageAccountDiff, 0, numDiff)
	for i := range accBalanceDiffs {
		accDiff = append(accDiff, accBalanceDiffs[i])
	}
	for i := range accCreationDiff {
		accDiff = append(accDiff, accCreationDiff[i])
	}
	for i := range accPubkeyDiff {
		accDiff = append(accDiff, accPubkeyDiff[i])
	}

}
