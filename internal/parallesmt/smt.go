package parallesmt

import (
	"sync"

	"github.com/vivijj/zigo/crypto/ff"
	"github.com/vivijj/zigo/crypto/poseidon"
)

const NARY = 4

// NodeDirection is child node direction relatively to its parent.
type NodeDirection = int

//             ______________(root)______________
//            |          |            |          |
//           LeftUp   LeftDown     RightUp    RightDown

const (
	LeftUp    NodeDirection = 0
	LeftDown  NodeDirection = 1
	RightUp   NodeDirection = 2
	RightDown NodeDirection = 3
)

// depth 0, 1 item(root), level 0
// depth 1, 2 item, level: [0,1]
// depth N: 2 ^ N item, level: [0,N]

// NodeIndex index(root) = 1
type NodeIndex = int

// ItemIndex in [0, N)
type ItemIndex = int

// NodeRef is index of node in the slice.
type NodeRef = int

type RwCacheMap struct {
	sync.RWMutex

	Value map[NodeIndex]ff.Element
}

func

// QuadNode for easily serialize, we just use slice to store the tree.
type QuadNode struct {
	depth int
	index NodeIndex
	// root has the node ref 0,so this child field can never be 0.
	child [4]int
}

type IndexHashPair struct {
	Index NodeIndex
	Hash  ff.Element
}

type SparseQuadMerkleTree[T any] struct {
	// list of stored items.
	Items map[ItemIndex]T
	// hasher to calculate hash in the merkle tree
	Hasher poseidon.Hasher
	// fixed depth of the tree
	TreeDepth int
	Root      int
	// list of intermediate nodes
	Nodes []QuadNode
	// cache of hashes for `default` nodes(e.g. ones that are absent in the tree)
	preHashed []ff.Element
	// cache storing the already calculated hashes for node.
	// allowing us to avoid calculating the hash of element more than once.
	cache RwCacheMap
}

func NewQuadSMT[T any](treeDepth int) SparseQuadMerkleTree[T] {
	hasher := poseidon.NewHasher(5)
	items := make(map[ItemIndex]T)
	nodes := []QuadNode{
		{
			depth: 0,
			index: 1,
			child: [4]int{-1, -1, -1, -1},
		},
	}
	preHash := make([]ff.Element, treeDepth+1)
	cur := hasher.HashElements(T.getElemnt())
	preHash[treeDepth] = cur
	preHash = append(preHash, cur)
	for i := 0; i < treeDepth; i++ {
		cur := hasher.HashElements([]ff.Element{cur, cur, cur, cur})
		preHash[treeDepth-i-1] = cur
	}
	cache :=

}

// Insert an element to the tree.
func (t *SparseQuadMerkleTree[T]) Insert(itemIndex int, item T) {
	t.Items[itemIndex] = item
	// leafIndex = (4**d + 2)/3 + itemIndex
	leafIndex := ((1<<(2*t.TreeDepth))+2)/3 + itemIndex

	// invalidate the root cache.
	t.cache.Lock()
	delete(t.cache.Value, 1)
	t.cache.Unlock()

	currentNodeRef := t.Root

	for {
		currentNode := t.Nodes[currentNodeRef]
		currentLevel := t.calculateLevel(currentNode.depth)

		dir := (itemIndex << (2 * currentLevel)) % NARY
		if nextRef := currentNode.child[dir]; nextRef > 0 {
			nextNode := t.Nodes[nextRef]
			// Normalized leaf index is basically an index of the node parent
			// to our leaf on the current level.
			leafIndexNormalized := t.normalizeIndex(leafIndex, nextNode.depth)
			if leafIndexNormalized == nextNode.index {
				// the `next` node is the node we should update.
				t.wipeCache(nextNode.index, currentNode.index)

				// We should go at least one full level deeper.
				if nextNode.index == leafIndex {
					// We reached the leaf, no further updating required.
					// All the outdated caches are invalidated, and the leaf value
					// was inserted below.
					break
				} else {
					// We didn't reach the leaf layer, thus we should keep going down the tree.
					currentNodeRef = nextRef
					continue
				}
			} else {
				// Next node is not the node we must update.
				// We have to insert one additional node which will have the
				// `next` node and our node as children.
				firstNodeIdx := leafIndexNormalized
				secondNodeIdx := nextNode.index
				var dirFirst, dirSecond NodeDirection
				for firstNodeIdx != secondNodeIdx {
					dirFirst = (firstNodeIdx + 2) % NARY
					dirSecond = (secondNodeIdx + 2) % NARY
					firstNodeIdx = t.parentIndex(firstNodeIdx)
					secondNodeIdx = t.parentIndex(secondNodeIdx)
				}
				commonParentIdx := firstNodeIdx

				// Invalidate the cache for the intersection point.
				t.wipeCache(commonParentIdx, currentNode.index)

				leafRef := t.insertNode(leafIndex, t.TreeDepth, [NARY]int{-1, -1, -1, -1})

				commonParentChild := [NARY]int{-1, -1, -1, -1}
				commonParentChild[dirFirst] = leafRef
				commonParentChild[dirSecond] = nextRef

				splitNodeDepth := t.depth(commonParentIdx)
				splitNodeRef := t.insertNode(commonParentIdx, splitNodeDepth, commonParentChild)

				t.addChildNode(currentNodeRef, dir, splitNodeRef)
				break
			}
		} else {
			// There is no child within the direction of the node to insert.
			// We must simply insert the leaf and make it a child of the latest
			// existing parent node.
			// No further processing is required.
			leafRef := t.insertNode(leafIndex, t.TreeDepth, [NARY]int{-1, -1, -1, -1})
			t.addChildNode(currentNodeRef, dir, leafRef)
		}
	}
}

// calculate the direct parent index of current node with index
// for quad tree: [(n-2+4)/4] = [(n+2)/4]
func (t *SparseQuadMerkleTree[T]) parentIndex(index NodeIndex) NodeIndex {
	return (index + 2) >> 2
}

func (t *SparseQuadMerkleTree[T]) normalizeIndex(leafIndex NodeIndex, depth int) NodeIndex {
	currentIndex := leafIndex
	for i := 0; i < t.TreeDepth-depth; i++ {
		currentIndex = t.parentIndex(currentIndex)
	}
	return currentIndex
}

// calculate the depth ('layer") of the element with provided index.
func (t *SparseQuadMerkleTree[T]) depth(index NodeIndex) int {
	level := 0
	for i := index; i > 1; {
		level += 1
		i = t.parentIndex(i)
	}
	return level
}

// Removes the entry with provided index from hashes cache, as well as its parent entries, limited
// by `parent` index
func (t *SparseQuadMerkleTree[T]) wipeCache(child NodeIndex, parent NodeIndex) {
	t.cache.Lock()
	if _, ok := t.cache.Value[child]; ok {
		// Item existed in cache, now we should delete it and go up the tree
		// and remove parent hashes, until we reach the provided
		// `parent` index.
		delete(t.cache.Value, child)
		parentI := t.parentIndex(child)
		for parentI > parent {
			delete(t.cache.Value, parentI)
			parentI = t.parentIndex(parentI)
		}
	}
}

func (t *SparseQuadMerkleTree[T]) insertNode(index NodeIndex, depth int, child [NARY]int) int {
	t.Nodes = append(
		t.Nodes, QuadNode{
			depth: depth,
			index: index,
			child: child,
		},
	)
	return len(t.Nodes) - 1
}

func (t *SparseQuadMerkleTree[T]) addChildNode(nodeRef NodeRef, dir NodeDirection, child NodeRef) {
	t.Nodes[nodeRef].child[dir] = child
}

func (t *SparseQuadMerkleTree[T]) calculateLevel(curDepth int) int {
	return t.TreeDepth - curDepth - 1
}
