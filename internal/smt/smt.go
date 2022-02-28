// Package smt implements Sparse Merkle Tree.
// Sparse Merkle Tree is a variation off the standard Merkle Tree, where contained
// data is indexed and each data point is placed in the leaf that corresponds to
// its index.
// The sparseness of the tree is implementing through a "default leaf" - an item which
// hash will be used for the missing indices instead of the actual element hash.
package smt

import (
	"github.com/vivijj/zigo/crypto/ff"
)

type MerkleItem interface {
	// HashFrContent return the poseidon hash of the item itself
	HashFrContent() ff.Element
}

type SparseMerkleTree[T MerkleItem] struct {
	HashTree MerkleHashTree `json:"hash_tree"`
	Items    map[int]T      `json:"-"`
}

func NewSparseMerkleTree[T MerkleItem](depth int) SparseMerkleTree[T] {
	var t T
	hashTree := NewHashTree(depth, t.HashFrContent())
	items := make(map[int]T)
	return SparseMerkleTree[T]{
		HashTree: hashTree,
		Items:    items,
	}
}

// Insert will insert an element to the tree, once complete,the Root hash will change
func (smt *SparseMerkleTree[T]) Insert(index int, item T) {
	newHash := item.HashFrContent()
	smt.Items[index] = item
	smt.HashTree.Update(index, newHash)
}

func (smt *SparseMerkleTree[T]) InsertJust(index int, item T) {
	smt.Items[index] = item
}

// Get return a flag to indicate that if the item exist
func (smt *SparseMerkleTree[T]) Get(index int) (T, bool) {
	res, ok := smt.Items[index]
	return res, ok
}

func (smt *SparseMerkleTree[T]) RootHash() ff.Element {
	return smt.HashTree.Root
}
