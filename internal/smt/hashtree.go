package smt

import (
	"math/big"

	"github.com/steinselite/zigo/crypto/poseidon"
)

// This file impl the hash calculate tree model(with no item)

const NARY = 4

// The hasher use in the merkle tree
var merkleHashParams = poseidon.NewParams(5, 6, 52)

func Hasher(inputS []string) string {
	inputBi := make([]*big.Int, len(inputS))
	for i := 0; i < len(inputS); i++ {
		inputBi[i], _ = new(big.Int).SetString(inputS[i], 10)
	}
	resBi := poseidon.Hash(inputBi, merkleHashParams)
	return resBi.String()
}

// This is the implemention of the sparse merkle tree with 4_n_ary
// This merkle tree only store the hash of item in the leaf(with no origin item data)
// all the hash is in string representation
type MerkleHashTree struct {
	depth int
	root  string
	// cache storing the already calculated hashes for nodes
	// it exist in [parent node hash] ==> [4]{4 children node hash}
	cache map[string][]string
}

// return an empty hash tree with default leaf hash
func NewHashTree(treeDepth int, defaultHash string) MerkleHashTree {
	chldHash := defaultHash
	cache := make(map[string][]string)

	var parentHash string
	for i := 0; i < treeDepth; i++ {
		chldGroup := []string{chldHash, chldHash, chldHash, chldHash}
		parentHash = Hasher(chldGroup)
		cache[parentHash] = chldGroup
		chldHash = parentHash
	}
	return MerkleHashTree{
		depth: treeDepth,
		root:  parentHash,
		cache: cache,
	}
}

// insert or update the tree,(all the operation could be thought as "update" due to sparseness)
func (ht *MerkleHashTree) Update(index int, itemHash string) {
	currentNodeRef := ht.root
	lookupPath := index
	upsertPath := index
	sideNode := [][]string{}

	// look up the path to find the item with index
	for i := 0; i < ht.depth; i++ {
		chld := ht.cache[currentNodeRef]
		sideNode = append(sideNode, chld)
		chldIndex := (lookupPath >> (2 * (ht.depth - 1))) % NARY
		currentNodeRef = chld[chldIndex]
		lookupPath = lookupPath << 2 // remove the first 2 bit
	}

	// update the merkle tree bottom up
	currentNodeRef = itemHash

	for i := 0; i < ht.depth; i++ {
		chldIndex := upsertPath % NARY
		leaves := []string{}

		for j := 0; j < NARY; j++ {
			if j != chldIndex {
				leaves = append(leaves, sideNode[ht.depth-1-i][j])
			} else {
				leaves = append(leaves, currentNodeRef)
			}
		}
		newRef := Hasher(leaves)
		ht.cache[newRef] = leaves
		upsertPath = upsertPath >> 2
		currentNodeRef = newRef
	}
	ht.root = currentNodeRef
}

// MerklePath create a proof of existence for a certain element of the tree.
// return value is list with length 4*depth(every level has 4 elemnt)
func (ht *MerkleHashTree) MerklePath(index int) []string {
	curItem := ht.root
	lookupPath := index
	// specify the capcity avoid the extend
	sideNodes := make([][]string, ht.depth)

	for i := 0; i < ht.depth; i++ {
		nodeRef := (lookupPath >> (2 * (ht.depth - 1))) % NARY
		level := ht.depth - 1 - i
		levelNode := []string{}

		for j := 0; j < NARY; j++ {
			if j != nodeRef {
				levelNode = append(levelNode, ht.cache[curItem][j])
			}
		}
		sideNodes[level] = levelNode
		curItem = ht.cache[curItem][nodeRef]
		lookupPath = lookupPath << 2
	}
	merkleProof := []string{}
	for i := 0; i < len(sideNodes); i++ {
		for j := 0; j < len(sideNodes[0]); j++ {
			merkleProof = append(merkleProof, sideNodes[i][j])
		}
	}
	return merkleProof

}

// VerifyProoof verify the given proof for the given element and index
func (ht *MerkleHashTree) VerifyProof(merkleProof []string, index int, itemHash string) bool {
	path := index
	curNode := itemHash
	proofIdx := 0

	for i := 0; i < ht.depth; i++ {
		levelNodes := []string{}
		for j := 0; j < NARY; j++ {
			if j == path%NARY {
				levelNodes = append(levelNodes, curNode)
			} else {
				levelNodes = append(levelNodes, merkleProof[proofIdx])
				proofIdx = proofIdx + 1
			}
		}
		curNode = Hasher(levelNodes)
		path = path >> 2
	}
	return ht.root == curNode
}
