package account

import (
	"fmt"
	"math/big"

	"github.com/vivijj/zigo/crypto/poseidon"
	"github.com/vivijj/zigo/internal/param"
	"github.com/vivijj/zigo/internal/smt"
	"github.com/vivijj/zigo/types/fr"
)

var (
	balancePoseidonParam = poseidon.NewParams(5, 6, 52)
	accountPoseidonParam = poseidon.NewParams(6, 6, 52)
)

type (
	CircuitAccountTree = smt.SparseMerkleTree[CircuitAccount]
	CircuitBalanceTree = smt.SparseMerkleTree[CircuitBalance]
)

type CircuitBalance struct {
	Value fr.Fr
}

func (b CircuitBalance) HashFrContent() fr.Fr {
	if b.Value == "" {
		b.Value = "0"
	}
	res := poseidon.Hash([]*big.Int{b.Value.ToBigInt()}, balancePoseidonParam)
	fmt.Println(res)
	return fr.FromBigInt(res)
}

// CircuitAccount Representation of account used in the `circuit`.
type CircuitAccount struct {
	SubTree smt.SparseMerkleTree[CircuitBalance]
	Nonce   fr.Fr
	PubkeyX fr.Fr
	PubkeyY fr.Fr
	Address fr.Fr
}

func NewCircuitAccount() CircuitAccount {
	return CircuitAccount{
		SubTree: smt.NewSparseMerkleTree[CircuitBalance](param.QuadAccountTreeDepth),
		Nonce:   "",
		PubkeyX: "",
		PubkeyY: "",
		Address: "",
	}
}

func (b CircuitAccount) HashFrContent() fr.Fr {
	content := make([]*big.Int, 0, 5)
	content = append(content, b.Address.ToBigInt())
	content = append(content, b.PubkeyX.ToBigInt())
	content = append(content, b.PubkeyY.ToBigInt())
	content = append(content, b.Nonce.ToBigInt())

	root := b.SubTree.RootHash()
	content = append(content, root.ToBigInt())
	res := poseidon.Hash(content, accountPoseidonParam)
	return fr.FromBigInt(res)
}
