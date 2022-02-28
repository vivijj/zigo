package record

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vivijj/zigo/pkg/account"
)

type StorageAccount struct {
	Id         int            `bson:"id"`
	LastBlock  int            `bson:"last_block"`
	Nonce      int            `bson:"nonce"`
	Address    []byte         `bson:"address"`
	PubkeyComp [32]byte       `bson:"pubkey_comp"`
	Balances   map[int]string `bson:"balances"`
}

func RestoreAccount(sacc StorageAccount) (blockNumber int, acc account.Account) {
	for k := range sacc.Balances {
		bi, ok := new(big.Int).SetString(sacc.Balances[k], 10)
		if !ok {
			panic("fail to SetString for big.Int")
		}
		acc.SetBalance(k, bi)
	}
	acc.Nonce = sacc.Nonce
	acc.PubkeyPair = account.PubkeyFromCompress(sacc.PubkeyComp)
	acc.Address = common.BytesToAddress(sacc.Address)
	return
}

type StorageAccountDiff interface {
	isStorageAccountDiff()
}

type StorageAccountCreation struct {
	AccountId   int
	IsCreate    bool
	BlockNumber int
	Address     []byte
	Nonce       int
}

func (_ StorageAccountCreation) isStorageAccountDiff() {}

type StorageAccountUpdate struct {
	AccountId   int    `bson:"account_id"`
	BlockNumber int    `bson:"block_number"`
	TokenId     int    `bson:"token_id"`
	OldBalance  string `bson:"old_balance"`
	NewBalance  string `bson:"new_balance"`
	OldNonce    int    `bson:"old_nonce"`
	NewNonce    int    `bson:"new_nonce"`
}

func (_ StorageAccountUpdate) isStorageAccountDiff() {}

type StorageAccountPubkeyUpdate struct {
	AccountId     int    `bson:"account_id"`
	BlockNumber   int    `bson:"block_number"`
	OldPubkeyComp []byte `bson:"old_pubkey_comp"`
	NewPubkeyComp []byte `bson:"new_pubkey_comp"`
	OldNonce      int    `bson:"old_nonce"`
	NewNonce      int    `bson:"new_nonce"`
}

func (_ StorageAccountPubkeyUpdate) isStorageAccountDiff() {}
