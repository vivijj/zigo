package state

import (
	"github.com/vivijj/zigo/internal/param"
	"github.com/vivijj/zigo/pkg/account"
	"github.com/vivijj/zigo/pkg/operation"
	"github.com/vivijj/zigo/pkg/transaction"
)

func (s *State) applyDepositTx(tx transaction.DepositTx) (OpSuccess, error) {
	op, err := s.createDepositOp(tx)
	if err != nil {
		return OpSuccess{}, err
	}
	updates, err := s.applyDepositOp(op)
	if err != nil {
		return OpSuccess{}, err
	}
	return OpSuccess{
		Fee:        nil,
		Updates:    updates,
		ExecutedTx: op,
	}, nil
}

func (s *State) createDepositOp(tx transaction.DepositTx) (operation.DepositOp, error) {
	if err := invariant(tx.Token <= param.TotalTokens, InvalidToken); err != nil {
		return operation.DepositOp{}, err
	}
	var accountId int
	if id, _, ok := s.GetAccountByAddress(tx.To); ok {
		accountId = id
	} else {
		accountId = s.GetFreeAccountId()
	}
	return operation.DepositOp{
		DepositTx: tx,
		AccountId: accountId,
	}, nil
}

func (s *State) applyDepositOp(op operation.DepositOp) (account.Updates, error) {
	var acc account.Account
	var upds account.Updates
	if tacc, ok := s.GetAccount(op.AccountId); ok {
		acc = tacc
	} else {
		tacc, upd := account.CreateAccount(op.AccountId, op.To)
		upds = append(upds, upd...)
		acc = tacc
	}

	oldAmount := acc.GetBalance(int(op.DepositTx.Token))
	oldNonce := acc.Nonce
	acc.AddOrSubBalance(int(op.DepositTx.Token), op.DepositTx.Amount)
	newAmount := acc.GetBalance(int(op.DepositTx.Token))
	s.InsertAccount(op.AccountId, acc)
	upds = append(
		upds, account.AccUpdateTuple{
			AccId: op.AccountId,
			AccUpdate: account.UpdateBalance{
				OldNonce:   oldNonce,
				NewNonce:   oldNonce,
				TokenId:    int(op.Token),
				OldBalance: oldAmount,
				NewBalance: newAmount,
			},
		},
	)
	return upds, nil
}
