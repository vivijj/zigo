package state

import (
	"math/big"

	"github.com/vivijj/zigo/crypto/babyjub"
	"github.com/vivijj/zigo/internal/param"
	"github.com/vivijj/zigo/pkg/account"
	"github.com/vivijj/zigo/pkg/operation"
	"github.com/vivijj/zigo/pkg/transaction"
)

func (s *State) applyWithdrawTx(tx transaction.WithdrawTx) (OpSuccess, error) {
	op, err := s.CreateWithdrawOp(tx)
	if err != nil {
		return OpSuccess{}, err
	}
	fee, updates, err := s.applyWithdrawOp(op)
	if err != nil {
		return OpSuccess{}, err
	}
	return OpSuccess{
		Fee:        fee,
		Updates:    updates,
		ExecutedTx: op,
	}, nil
}

func (s *State) CreateWithdrawOp(tx transaction.WithdrawTx) (op operation.WithdrawOp, err error) {
	// do the check at first
	if err = invariant(tx.Token <= param.MaxTokenId, InvalidTokenId); err != nil {
		return
	}
	if tx.Fee.BitLen() != 0 {
		if err = invariant(tx.FeeToken <= param.MaxTokenId, InvalidFeeTokenId); err != nil {
			return
		}
	}
	accId, acc, ok := s.GetAccountByAddress(tx.From)
	if !ok {
		err = FromAccountNotFound
		return
	}
	if err = invariant(!acc.PubkeyPair.IsEmpty(), FromAccountLocked); err != nil {
		return
	}
	pubkey := babyjub.PublicKey(acc.PubkeyPair)
	msg := tx.EncodeBi()
	if !pubkey.VerifyPoseidon(msg, &tx.Signature) {
		err = InvalidSignature
		return
	}
	if err = invariant(accId == tx.AccountId, FromAccountIncorrect); err != nil {
		return
	}

	op.WithdrawTx = tx
	op.MaxFee = tx.Fee
	op.ConditionType = 0
	return
}

func (s *State) applyWithdrawOp(op operation.WithdrawOp) (
	fee *CollectedFee,
	upds account.Updates,
	err error,
) {
	accfrom, _ := s.GetAccount(op.AccountId)

	fromOldBalance := accfrom.GetBalance(op.Token)
	fromOldNonce := accfrom.Nonce

	if err = invariant(op.Nonce == fromOldNonce, NonceMismatch); err != nil {
		return
	}
	if err = invariant(
		fromOldBalance.Cmp(big.NewInt(0).Add(op.Amount, op.Fee)) != -1,
		InsufficientBalance,
	); err != nil {
		return
	}

	accfrom.SubBalance(op.Token, big.NewInt(0).Add(op.Amount, op.Fee))
	accfrom.Nonce += 1

	fromNewBalance := accfrom.GetBalance(op.Token)
	fromNewNonce := accfrom.Nonce
	s.InsertAccount(op.AccountId, accfrom)

	upds = append(
		upds, account.AccUpdateTuple{
			AccId: op.AccountId,
			AccUpdate: account.UpdateBalance{
				OldNonce:   fromOldNonce,
				NewNonce:   fromNewNonce,
				TokenId:    op.Token,
				OldBalance: fromOldBalance,
				NewBalance: fromNewBalance,
			},
		},
	)

	fee = &CollectedFee{
		Token:  op.FeeToken,
		Amount: op.Fee,
	}

	return
}
