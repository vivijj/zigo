package account

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Update interface {
	isAccountUpdate()
}

type Create struct {
	Address common.Address
	Nonce   int
}

func (_ Create) isAccountUpdate() {}

type UpdateBalance struct {
	OldNonce   int
	NewNonce   int
	TokenId    int
	OldBalance *big.Int
	NewBalance *big.Int
}

func (_ UpdateBalance) isAccountUpdate() {}

type PubKeyUpdate struct {
	OldPubkeyPair PubKeyPair
	NewPubkeyPair PubKeyPair
	OldNonce      int
	NewNonce      int
}

func (_ PubKeyUpdate) isAccountUpdate() {}

type AccUpdateTuple struct {
	AccId     int
	AccUpdate Update
}
type Updates = []AccUpdateTuple
