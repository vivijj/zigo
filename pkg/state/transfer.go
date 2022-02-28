package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vivijj/zigo/crypto/babyjub"
	"github.com/vivijj/zigo/internal/param"
	"github.com/vivijj/zigo/pkg/account"
	"github.com/vivijj/zigo/pkg/operation"
	"github.com/vivijj/zigo/pkg/transaction"
)

func (s *State) applyTransferTx(tx transaction.TransferTx) (OpSuccess, error) {
	op, err := s.createTransferOp(tx)
	if err != nil {
		return OpSuccess{}, err
	}
	fee, updates, err := s.applyTransferOp(op)
	if err != nil {
		return OpSuccess{}, err
	}
	return OpSuccess{
		Fee:        fee,
		Updates:    updates,
		ExecutedTx: op,
	}, nil
}

func (s *State) createTransferOp(tx transaction.TransferTx) (op operation.TransferOp, err error) {

	if err = invariant(tx.Token <= param.MaxTokenId, InvalidTokenId); err != nil {
		return
	}
	if tx.Fee.BitLen() != 0 {
		if err = invariant(tx.FeeToken <= param.MaxTokenId, InvalidFeeTokenId); err != nil {
			return
		}
	}
	if err = invariant(tx.To != common.Address{}, TargetAccountZero); err != nil {
		return
	}
	fromId, fromAcc, ok := s.GetAccountByAddress(tx.From)
	if !ok {
		err = FromAccountNotFound
		return
	}

	if fromAcc.PubkeyPair.IsEmpty() {
		err = FromAccountLocked
		return
	}
	pubkey := babyjub.PublicKey(fromAcc.PubkeyPair)
	msg := tx.EncodeBi()
	if !pubkey.VerifyPoseidon(msg, &tx.Signature) {
		err = InvalidSignature
		return
	}
	if err = invariant(fromId == tx.AccountId, TransferAccountIncorrect); err != nil {
		return
	}

	if toId, _, ok := s.GetAccountByAddress(tx.To); ok {
		return operation.TransferOp{
			TransferTx:     tx,
			ToId:           toId,
			ConditionType:  0,
			MaxFee:         tx.Fee,
			PutAddressInDa: true,
		}, nil
	}
	// if the to account not exist, just use the next free account id.
	toId := s.GetFreeAccountId()
	return operation.TransferOp{
		TransferTx:     tx,
		ToId:           toId,
		ConditionType:  0,
		MaxFee:         tx.Fee,
		PutAddressInDa: true,
	}, nil
}

func (s *State) applyTransferOp(op operation.TransferOp) (
	fee *CollectedFee, upds account.Updates, err error,
) {
	accfrom, _ := s.GetAccount(op.AccountId)
	var accto account.Account
	var hasacc bool
	if accto, hasacc = s.GetAccount(op.ToId); !hasacc {
		// if to account not exist means transfer to new, we should create the account first.
		tacc, upd := account.CreateAccount(op.ToId, op.To)
		upds = append(upds, upd...)
		accto = tacc
	}

	fromOldBalance := accfrom.GetBalance(op.Token)
	fromOldNonce := accfrom.Nonce

	if err = invariant(op.Nonce == fromOldNonce, NonceMismatch); err != nil {
		return
	}

	// fromOldBalance >= op.Amount + op.Fee
	if err = invariant(
		fromOldBalance.Cmp(new(big.Int).Add(op.Amount, op.Fee)) != -1,
		InsufficientBalance,
	); err != nil {
		return
	}

	accfrom.SubBalance(op.Token, new(big.Int).Add(op.Amount, op.Fee))
	accfrom.Nonce += 1

	fromNewBalance := accfrom.GetBalance(op.Token)
	fromNewNonce := accfrom.Nonce

	toOldBalance := accto.GetBalance(op.Token)
	toNonce := accto.Nonce
	accto.AddBalance(op.Token, op.Amount)
	toNewBalance := accto.GetBalance(op.Token)

	s.InsertAccount(op.AccountId, accfrom)
	s.InsertAccount(op.ToId, accto)

	upds = append(
		upds, account.AccUpdateTuple{
			AccId: op.AccountId,
			AccUpdate: account.UpdateBalance{
				OldNonce:   fromOldNonce,
				NewNonce:   fromNewNonce,
				TokenId:    op.AccountId,
				OldBalance: fromOldBalance,
				NewBalance: fromNewBalance,
			},
		},
	)

	upds = append(
		upds, account.AccUpdateTuple{
			AccId: op.ToId,
			AccUpdate: account.UpdateBalance{
				OldNonce:   toNonce,
				NewNonce:   toNonce,
				TokenId:    op.Token,
				OldBalance: toOldBalance,
				NewBalance: toNewBalance,
			},
		},
	)

	fee = &CollectedFee{
		Token:  op.FeeToken,
		Amount: op.Fee,
	}
	return
}
