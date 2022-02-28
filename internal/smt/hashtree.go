package smt

import (
	"fmt"

	"github.com/vivijj/zigo/crypto/ff"
	"github.com/vivijj/zigo/crypto/poseidon"
)

// This file impl the hash calculate tree model

const NARY = 4

// MerkleHashTree This is the implementation of the sparse merkle tree with 4_n_ary
// This merkle tree only store the hash of item in the leaf(with no origin item data)
// all the hash is in fr representation
type MerkleHashTree struct {
	Depth  int                             `json:"depth"`
	Root   ff.Element                      `json:"root"`
	Cache  map[ff.Element][NARY]ff.Element `json:"cache"`
	Hasher poseidon.Hasher                 `json:"-"`
}

// NewHashTree return an empty hash tree with default leaf hash
func NewHashTree(treeDepth int, defaultHash ff.Element) MerkleHashTree {
	hasher := poseidon.NewHasher(NARY + 1)
	fmt.Println("default is ", defaultHash)
	// calculate the prehash
	cache := make(map[ff.Element][NARY]ff.Element)
	curHash := defaultHash
	for i := 0; i < treeDepth; i++ {
		chldGroup := [NARY]ff.Element{curHash, curHash, curHash, curHash}
		curHash := hasher.HashElements(chldGroup[:])
		cache[curHash] = chldGroup
	}
	return MerkleHashTree{
		Depth:  treeDepth,
		Root:   curHash,
		Cache:  cache,
		Hasher: *hasher,
	}
}

// Update insert or update the tree,(all the operation could be thought as "update" due to sparseness)
func (ht *MerkleHashTree) Update(index int, itemHash ff.Element) {
	curNodeRef := ht.Root
	lookupPath := index
	upsertPath := index
	var sideNode [][NARY]ff.Element

	// look up the path to find the item with index
	for i := 0; i < ht.Depth; i++ {
		chld := ht.Cache[curNodeRef]
		sideNode = append(sideNode, chld)
		chldIndex := (lookupPath >> (2 * (ht.Depth - 1))) % NARY
		curNodeRef = chld[chldIndex]
		lookupPath = lookupPath << 2 // remove the first 2 bit
	}

	// update the merkle tree bottom up
	curNodeRef = itemHash

	for i := 0; i < ht.Depth; i++ {
		chldIndex := upsertPath % NARY
		var chldGroup [NARY]ff.Element

		for j := 0; j < NARY; j++ {
			if j != chldIndex {
				chldGroup[j] = sideNode[ht.Depth-1-i][j]
			} else {
				chldGroup[j] = curNodeRef
			}
		}
		newRef := ht.Hasher.HashElements(chldGroup[:])
		ht.Cache[newRef] = chldGroup
		upsertPath = upsertPath >> 2
		curNodeRef = newRef
	}
	ht.Root = curNodeRef
}

// MerklePath create a proof of existence for a certain element of the tree.
// return value is list with length 4*Depth(every level has 4 element)
func (ht *MerkleHashTree) MerklePath(index int) []ff.Element {
	curItem := ht.Root
	lookupPath := index
	// specify the capacity avoid the extent
	sideNodes := make([][]ff.Element, ht.Depth)

	for i := 0; i < ht.Depth; i++ {
		nodeRef := (lookupPath >> (2 * (ht.Depth - 1))) % NARY
		level := ht.Depth - 1 - i
		var levelNode []ff.Element

		for j := 0; j < NARY; j++ {
			if j != nodeRef {
				levelNode = append(levelNode, ht.Cache[curItem][j])
			}
		}
		sideNodes[level] = levelNode
		curItem = ht.Cache[curItem][nodeRef]
		lookupPath = lookupPath << 2
	}
	var merkleProof []ff.Element
	for i := 0; i < len(sideNodes); i++ {
		for j := 0; j < len(sideNodes[0]); j++ {
			merkleProof = append(merkleProof, sideNodes[i][j])
		}
	}
	return merkleProof

}

// VerifyProof verify the given proof for the given element and index
func (ht *MerkleHashTree) VerifyProof(
	merkleProof []ff.Element,
	index int,
	itemHash ff.Element,
) bool {
	path := index
	curNode := itemHash
	proofIdx := 0

	for i := 0; i < ht.Depth; i++ {
		var levelNodes []ff.Element
		for j := 0; j < NARY; j++ {
			if j == path%NARY {
				levelNodes = append(levelNodes, curNode)
			} else {
				levelNodes = append(levelNodes, merkleProof[proofIdx])
				proofIdx = proofIdx + 1
			}
		}
		curNode = ht.Hasher.HashElements(levelNodes)
		path = path >> 2
	}
	return ht.Root == curNode
}
