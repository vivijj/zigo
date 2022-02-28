package witness

import (
	"github.com/vivijj/zigo/internal/smt"
	"github.com/vivijj/zigo/pkg/account"
	"github.com/vivijj/zigo/types/fr"
)

// Fa is function account to deal with account
type Fa func(circuitAccount *account.CircuitAccount)

// Fb is function balance to deal with balance update
type Fb func(balance *account.CircuitBalance)

// GetAudits returns the merkle path as the audit of validation.
// NOTE: only account from will update 2 kind of token
func GetAudits(
	tree *smt.SparseMerkleTree[account.CircuitAccount],
	accId int,
	tokenIds ...int,
) ([]fr.Fr, []fr.Fr, []fr.Fr) {
	auditAccount := tree.HashTree.MerklePath(accId)
	acc, ok := tree.Get(accId)
	if !ok {
		acc = account.NewCircuitAccount()
	}
	auditBalanceA := acc.SubTree.HashTree.MerklePath(tokenIds[0])
	var auditBalanceB []fr.Fr
	if len(tokenIds) == 2 {
		auditBalanceB = acc.SubTree.HashTree.MerklePath(tokenIds[1])
	}
	return auditAccount, auditBalanceA, auditBalanceB
}

func ApplyLeafOperation(
	tree *smt.SparseMerkleTree[account.CircuitAccount],
	accId int,
	fa Fa,
	fb Fb,
	tokenIds ...int,
) {

}
