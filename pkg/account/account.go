// Package account impl the zion network account(L2 account)
package account

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vivijj/zigo/crypto/babyjub"
	"github.com/vivijj/zigo/types/fr"
)

type (
// BalanceTree = smt.SparseMerkleTree[Balance]
// AccountTree = smt.SparseMerkleTree[Account]
)

type PubKeyPair babyjub.PublicKey

func NewPubKey() PubKeyPair {
	return PubKeyPair{
		X: big.NewInt(0),
		Y: big.NewInt(0),
	}
}
func PubkeyFromCompress(compk [32]byte) PubKeyPair {
	pubkeyComp := babyjub.PublicKeyComp(compk)
	pubkey, _ := pubkeyComp.Decompress()
	return PubKeyPair(*pubkey)
}

func (p PubKeyPair) IsEmpty() bool {
	return p.X.BitLen() == 0 && p.Y.BitLen() == 0
}

// Account represents zion network account.
type Account struct {
	// Address is corresponding to the L1 account address
	Address common.Address
	// Public key of the account used to authorize operation for this account.
	PubkeyPair PubKeyPair
	// Nonce to avoid double spend in L2.
	Nonce    int
	Balances map[int]*big.Int
}

// Clone return a deep copy of the account.
func (acc *Account) Clone() (a Account) {
	a.Address = acc.Address
	a.PubkeyPair = NewPubKey()
	a.PubkeyPair.X.Set(acc.PubkeyPair.X)
	a.PubkeyPair.Y.Set(acc.PubkeyPair.Y)
	a.Nonce = acc.Nonce

	a.Balances = make(map[int]*big.Int)
	for i := range acc.Balances {
		a.Balances[i] = new(big.Int).Set(acc.Balances[i])
	}
	return
}

func (acc Account) ToCircuitAcc() CircuitAccount {
	ca := NewCircuitAccount()
	fmt.Println("here")
	for i := range acc.Balances {
		cb := CircuitBalance{Value: fr.FromBigInt(acc.Balances[i])}
		ca.SubTree.InsertJust(i, cb)
	}
	ca.Nonce = fr.FromInt(acc.Nonce)
	ca.PubkeyX = fr.FromBigInt(acc.PubkeyPair.X)
	ca.PubkeyY = fr.FromBigInt(acc.PubkeyPair.Y)
	ca.Address = fr.FromAddress(acc.Address)
	return ca
}

// ItemHasher return the poseidon hash of the account.
// NOTE: it will convert the account to circuit account and calculate the
// sub-balance-tree merkle root, so it's maybe a time-consuming process.
func (acc Account) HashFrContent() fr.Fr {
	return ""
}

// NewAccountAddress return a new account with address set.
func NewAccountAddress(address common.Address) *Account {
	return &Account{
		Address:    address,
		PubkeyPair: NewPubKey(),
		Nonce:      0,
		Balances:   make(map[int]*big.Int),
	}
}

// CreateAccount creates a new account object and list of updates has to be applied on the state
// in order to get this account created within the network.
func CreateAccount(id int, addr common.Address) (acc *Account, upds Updates) {
	acc = NewAccountAddress(addr)
	upds = []AccUpdateTuple{
		{
			AccId: id,
			AccUpdate: Create{
				Address: addr,
				Nonce:   acc.Nonce,
			},
		},
	}
	return
}

// GetBalance return the balance, if not exist return 0, notice that the balance is a new one
func (acc *Account) GetBalance(tokenId int) *big.Int {
	if v, ok := acc.Balances[tokenId]; ok {
		return new(big.Int).Set(v)
	}
	return big.NewInt(0)
}

// SetBalance Overrides the token balance value
func (acc *Account) SetBalance(tokenId int, amount *big.Int) {
	acc.Balances[tokenId] = amount
}

// SubBalance subtracts the provided amount from the token balance
// panic if the amount to subtract is greater than the existing token balance.
func (acc *Account) SubBalance(tokenId int, amount *big.Int) {
	balance := acc.GetBalance(tokenId)
	if balance.Cmp(amount) == -1 {
		panic("account balance can't be less than 0")
	}
	balance.Sub(balance, amount)
}

// AddBalance add provided amount to the token balance
func (acc *Account) AddBalance(tokenId int, amount *big.Int) {
	balance := acc.GetBalance(tokenId)
	balance = balance.Add(balance, amount)
}

// ApplyUpdates apply the list of update, change the account state.
func ApplyUpdates(acc *Account, upds []Update) *Account {
	for i := range upds {
		acc = ApplyUpdate(acc, upds[i])
	}
	return acc
}

// ApplyUpdate applies an update to the account state.
func ApplyUpdate(acc *Account, upd Update) *Account {
	if acc == nil {
		// none account is only allow to apply to create update
		switch upd := upd.(type) {
		case Create:
			return &Account{
				Address:    upd.Address,
				PubkeyPair: NewPubKey(),
				Nonce:      upd.Nonce,
				Balances:   make(map[int]*big.Int),
			}
		default:
			return nil
		}
	}
	switch upd := upd.(type) {
	case UpdateBalance:
		acc.SetBalance(upd.TokenId, upd.NewBalance)
		acc.Nonce = upd.NewNonce
	case PubKeyUpdate:
		acc.PubkeyPair = upd.NewPubkeyPair
		acc.Nonce = upd.NewNonce
	default:
		fmt.Println("incorrect update received for account")
	}
	return acc

}
