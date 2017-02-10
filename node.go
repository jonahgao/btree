package btree

import "bytes"

type node interface {
	isLeaf() bool
	keyAt(int) []byte
	numOfKeys() int
	leftMostKey() []byte

	get([]byte) []byte
	insert([]byte, []byte, uint64) *insertResult
	delete([]byte, uint64, node, int) *deleteResult
}

type baseNode struct {
	tree     *MVCCBtree
	revision uint64
	keys     [][]byte
}

func (n *baseNode) minKeys() int {
	return (n.tree.GetOrder()+1)/2 - 1
}

func (n *baseNode) maxKeys() int {
	return n.tree.GetOrder() - 1
}

func (n *baseNode) keyAt(pos int) []byte {
	return n.keys[pos]
}

func (n *baseNode) numOfKeys() int {
	return len(n.keys)
}

func (n *baseNode) splitPivot() int {
	return n.tree.GetOrder() / 2
}

// find equal or greater key pos. return exist(euqal) and index
func (n *baseNode) findPos(key []byte) (bool, int) {
	for i := 0; i < len(n.keys); i++ {
		c := bytes.Compare(n.keys[i], key)
		if c == 0 {
			return true, i
		} else if c > 0 {
			return false, i
		}
	}

	return false, len(n.keys)
}

func (n *baseNode) insertKeyAt(idx int, key []byte) {
	lastCnt := len(n.keys)
	n.keys = append(n.keys, nil)
	for i := lastCnt - 1; i >= idx; i-- {
		n.keys[i+1] = n.keys[i]
	}
	n.keys[idx] = key
}

// for delete: select a sibling for borrow or merge
func (n *baseNode) selectSibling(parent node, parentPos int) int {
	if parentPos == 0 {
		return 1
	}

	if parentPos == parent.numOfKeys() {
		return parentPos - 1
	}

	leftSibling := (parent.(*internalNode)).childAt(parentPos - 1)
	rightSibling := (parent.(*internalNode)).childAt(parentPos + 1)
	if leftSibling.numOfKeys() >= rightSibling.numOfKeys() {
		return parentPos - 1
	} else {
		return parentPos + 1
	}
}
