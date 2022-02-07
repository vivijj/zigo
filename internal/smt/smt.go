// Package smt implements Sparse Merkle Tree.
/// Sparse Merkle Tree is a variation off the standard Merkle Tree, where contained
/// data is indexed and each data point is placed in the leaf that corresponds to
/// its index.
/// The sparseness of the tree is implementing through a "default leaf" - an item which
/// hash will be used for the missing indices instead of the actual element hash.
package smt

type MerkleItem interface {
	// ItemHasher return the poseidon hash of the item itself
	ItemHasher() string
}

type SparseMerkleTree[T MerkleItem] struct {
	HashTree MerkleHashTree
	items    map[int]T
}

func NewSparseMerkleTree[T MerkleItem](depth int) SparseMerkleTree[T] {
	var t T
	hashTree := NewHashTree(depth, t.ItemHasher())
	items := make(map[int]T)
	return SparseMerkleTree[T]{
		HashTree: hashTree,
		items:    items,
	}
}

// insert an element to the tree, once complete,the root hash will change
func (smt *SparseMerkleTree[T]) Insert(index int, item T) {
	newHash := item.ItemHasher()
	smt.items[index] = item
	smt.HashTree.Update(index, newHash)
}

// return a flag to indicate that if the item exist
func (smt *SparseMerkleTree[T]) Get(index int) (T, bool) {
	res, ok := smt.items[index]
	return res, ok
}

func (smt *SparseMerkleTree[T]) RootHash() string {
	return smt.HashTree.root
}
