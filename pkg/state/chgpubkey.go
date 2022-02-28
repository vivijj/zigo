package state

import (
	"github.com/vivijj/zigo/internal/param"
	"github.com/vivijj/zigo/pkg/account"
	"github.com/vivijj/zigo/pkg/operation"
	"github.com/vivijj/zigo/pkg/transaction"
)

// Tx handler deal with pubkey update

func (s *State) applyPubkeyUpdateTx(tx transaction.PubkeyUpdateTx) (OpSuccess, error) {
	op, err := s.createPubkeyUpdateOp(tx)
	if err != nil {
		return OpSuccess{}, err
	}
	fee, updates, err := s.applyPubkeyUpdateOp(op)
	if err != nil {
		return OpSuccess{}, err
	}
	return OpSuccess{
		Fee:        fee,
		Updates:    updates,
		ExecutedTx: op,
	}, nil

}

func (s *State) createPubkeyUpdateOp(tx transaction.PubkeyUpdateTx) (
	op operation.PubkeyUpdateOp, err error,
) {
	accId, acc, ok := s.GetAccountByAddress(tx.Account)
	if !ok {
		err = AccountNotFound
		return
	}
	if err = invariant(tx.Account == acc.Address, InvalidAccountAddress); err != nil {
		return
	}
	if err = invariant(tx.FeeToken <= param.MaxTokenId, InvalidFeeTokenId); err != nil {
		return
	}
	if err = invariant(tx.IsAuthDataValid(), InvalidAuthData); err != nil {
		return
	}
	if err = invariant(accId == tx.AccountId, InvalidAccountId); err != nil {
		return
	}
	op.PubkeyUpdateTx = tx
	// now always from l1 to change the pubkey
	op.ConditionType = 1
	op.MaxFee = tx.Fee
	return
}

func (s *State) applyPubkeyUpdateOp(op operation.PubkeyUpdateOp) (
	fee *CollectedFee, upds account.Updates, err error,
) {
	acc, _ := s.GetAccount(op.AccountId)
	oldBalance := acc.GetBalance(op.FeeToken)
	oldPubkey := acc.PubkeyPair
	oldNonce := acc.Nonce

	if err = invariant(op.Nonce == acc.Nonce, NonceMismatch); err != nil {
		return
	}
	acc.Nonce += 1

	// update pubkey
	acc.PubkeyPair = account.PubKeyPair(op.PubKey)

	if err = invariant(oldBalance.Cmp(op.Fee) != -1, InsufficientBalance); err != nil {
		return
	}
	acc.SubBalance(op.FeeToken, op.Fee)

	newPubkey := acc.PubkeyPair
	newNonce := acc.Nonce
	newBalance := acc.GetBalance(op.FeeToken)

	s.InsertAccount(op.AccountId, acc)

	upds = append(
		upds,
		account.AccUpdateTuple{
			AccId: op.AccountId,
			AccUpdate: account.PubKeyUpdate{
				OldPubkeyPair: oldPubkey,
				NewPubkeyPair: newPubkey,
				OldNonce:      oldNonce,
				NewNonce:      newNonce,
			},
		},
	)

	upds = append(
		upds,
		account.AccUpdateTuple{
			AccId: op.AccountId,
			AccUpdate: account.UpdateBalance{
				OldNonce:   newNonce,
				NewNonce:   newNonce,
				TokenId:    op.FeeToken,
				OldBalance: oldBalance,
				NewBalance: newBalance,
			},
		},
	)
	fee = &CollectedFee{
		Token:  op.FeeToken,
		Amount: op.Fee,
	}
	return
}
