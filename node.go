package btree

import "bytes"

type node interface {
	isLeaf() bool
	getValue([]byte) []byte
	insert([]byte, []byte, uint64) *insertResult
}

type baseNode struct {
	tree     *btree
	parent   node
	revision uint64
	keys     [][]byte
}

func (n *baseNode) isRoot() bool {
	return n.parent == nil
}

func (n *baseNode) minKeys() int {
	if n.isRoot() {
		return 1
	}
	return (n.tree.order+1)/2 - 1
}

func (n *baseNode) maxKeys() int {
	return n.tree.order - 1
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

func (n *baseNode) setParent(p node) {
	n.parent = p
}
