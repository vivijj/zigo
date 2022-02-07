package smt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashTree(t *testing.T) {
	ht := NewHashTree(16, "1234")
	ht.Update(1, "5678")
	proof := ht.MerklePath(2)
	res := ht.VerifyProof(proof, 2, "1234")
	assert.Equal(
		t,
		true,
		res,
	)
}

func BenchmarkHashTree(b *testing.B) {
	ht := NewHashTree(16, "1234")
	ht.Update(1, "5678")
	proof := ht.MerklePath(2)
	ht.VerifyProof(proof, 2, "1234")
}
