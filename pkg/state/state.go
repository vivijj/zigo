package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vivijj/zigo/internal/param"
	"github.com/vivijj/zigo/internal/smt"
	"github.com/vivijj/zigo/pkg/account"
	"github.com/vivijj/zigo/pkg/operation"
	"github.com/vivijj/zigo/pkg/transaction"
	"github.com/vivijj/zigo/types/cmp"
)

type CollectedFee struct {
	Token  int
	Amount *big.Int
}

type OpSuccess struct {
	Fee        *CollectedFee
	Updates    account.Updates
	ExecutedTx operation.ZionOp
}

type State struct {
	// accounts stored in a sparse merkle tree
	accountTree        smt.SparseMerkleTree[account.Account]
	accountIdByAddress map[common.Address]int
	nextFreeId         int
	BlockNumber        int
}

func New(tree account.AccountTree, accIdByAddr map[common.Address]int, currentBlock int) *State {
	nextFreeId := 0
	for k := range tree.Items {
		nextFreeId = cmp.Max(nextFreeId, k+1)
	}
	return &State{
		accountTree:        tree,
		accountIdByAddress: accIdByAddr,
		nextFreeId:         nextFreeId,
		BlockNumber:        currentBlock,
	}
}

func FromAccMap(accs map[int]account.Account, currentBlock int) *State {
	s := &State{
		accountTree:        smt.NewSparseMerkleTree[account.Account](param.QuadAccountTreeDepth),
		accountIdByAddress: make(map[common.Address]int),
		nextFreeId:         0,
		BlockNumber:        0,
	}
	nextFreeId := 0

	for i := range accs {
		nextFreeId = cmp.Max(nextFreeId, i+1)
	}
	s.nextFreeId = nextFreeId
	s.BlockNumber = currentBlock
	for id, acc := range accs {
		s.InsertAccount(id, acc)
	}
	return s
}

func (s *State) RootHash() string {
	return s.accountTree.RootHash()
}

func (s *State) GetAccount(accountId int) (account.Account, bool) {
	if acc, ok := s.accountTree.Get(accountId); ok {
		return acc, true
	}
	return account.Account{}, false
}

// ExecutePriorityTx will execute priority transaction which should not fail on L2.
func (s *State) ExecutePriorityTx(tx transaction.ZionPriTx) OpSuccess {
	switch ztx := tx.(type) {
	case transaction.DepositTx:
		opres, err := s.applyDepositTx(ztx)
		if err != nil {
			panic("priority transaction execution failed")
		}
		return opres
	default:
		panic("not invalid priority tx")
	}
}

func (s *State) ExecuteTx(tx transaction.ZionTx) (ops OpSuccess, err error) {
	switch ztx := tx.(type) {
	case transaction.TransferTx:
		ops, err = s.applyTransferTx(ztx)
	case transaction.WithdrawTx:
		ops, err = s.applyWithdrawTx(ztx)
	case transaction.PubkeyUpdateTx:
		ops, err = s.applyPubkeyUpdateTx(ztx)
	}
	return
}

func (s *State) GetFreeAccountId() int {
	return s.nextFreeId
}

// GetAccountByAddress get (accountId, account, ifExist)
func (s *State) GetAccountByAddress(address common.Address) (int, account.Account, bool) {
	if accId, ok := s.accountIdByAddress[address]; ok {
		if acc, hasAcc := s.GetAccount(accId); hasAcc {
			return accId, acc, true
		}
		// This should never happen, if we find account in accountIdByAddress map,
		// we can always get the account from the tree
		panic("failed to get account")
	}
	// the account associated with the address is not exist
	return 0, account.Account{}, false
}

// InsertAccount inserts (id, account) into the account tree,
// while we insert account, the merkle tree hash is being recalculated at same time.
func (s *State) InsertAccount(id int, acc account.Account) {
	s.accountIdByAddress[acc.Address] = id
	s.accountTree.Insert(id, acc)
}

func (s *State) CollectFee(fees []CollectedFee, operatorId int) (upds account.Updates) {
	acc, ok := s.GetAccount(operatorId)
	if !ok {
		panic("operator account should present in the account tree: ")
	}
	for i := range fees {
		if fees[i].Amount.BitLen() == 0 {
			continue
		}

		oldAmount := new(big.Int).Set(acc.GetBalance(fees[i].Token))
		nonce := acc.Nonce
		acc.AddBalance(fees[i].Token, fees[i].Amount)
		newAmount := new(big.Int).Set(acc.GetBalance(fees[i].Token))

		upds = append(
			upds,
			account.AccUpdateTuple{
				AccId: operatorId,
				AccUpdate: account.UpdateBalance{
					OldNonce:   nonce,
					NewNonce:   nonce,
					TokenId:    fees[i].Token,
					OldBalance: oldAmount,
					NewBalance: newAmount,
				},
			},
		)
	}
	s.InsertAccount(operatorId, acc)
	return
}
