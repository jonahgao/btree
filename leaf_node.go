package btree

type leafNode struct {
	baseNode
	values [][]byte
}

func newLeafNode(t *btree, p node, r uint64) *leafNode {
	return &leafNode{
		tree:     t,
		parent:   p,
		revision: r,
		keys:     make([][]byte, 0, t.order-1),
		values:   make([][]byte, 0, t.order-1),
	}
}

func (n *leafNode) isLeaf() bool {
	return true
}

func (n *leafNode) getValue([]byte) []byte {
	exist, idx := n.findPos(key)
	if !exist {
		return nil
	}
	return n.values[idx]
}

func (n *leafNode) insertValueAt(idx int, value []byte) {
	lastCnt := len(n.values)
	n.values = append(n.values, nil)
	// move backward
	for i := lastCnt - 1; i >= idx; i-- {
		n.values[i+1] = n.values[i]
	}
	n.values[idx] = value
}

func (n *leafNode) clone(revision uint64) *leafNode {
	newLeaf := newLeafNode(n.tree, n.parent, revision)
	for _, k := range n.keys {
		newLeaf.keys = append(newLeaf.keys, k)
	}
	for _, v := range n.values {
		newLeaf.values = append(newLeaf.values, v)
	}
}

func (n *leafNode) replaceValue(pos int, value []byte, revision uint64) *insertResult {
	newLeaf := n
	if n.revision != revision {
		newLeaf = n.clone(revision)
	}
	newLeaf.values[pos] = value

	return &insertResult{
		rtype:    iRTypeModified,
		modified: newLeaf,
	}
}

func (n *leafNode) addValue(pos int, key, value []byte, revision uint64) *insertResult {

}

func (n *leafNode) addAndSplit(pos int, key, value []byte, revision uint64) *insertResult {
}

func (n *leafNode) insert(key, value []byte, revision uint64) *insertResult {
	exist, pos := n.findPos(key)
	if exist {
		return n.replaceValue()
	} else if len(n.keys)+1 <= n.maxKeys() {
		return n.addValue(pos, key, value, revsion)
	} else {
		return n.addAndSplit(pos, key, value, revision)
	}
}
