package btree

type internalNode struct {
	baseNode
	children []node
}

func (n *internalNode) isLeaf() bool {
	return false
}

func (n *internalNode) getValue([]byte) []byte {
	exist, idx := n.findPos(key)
	// equal: go to right
	if exist {
		return n.children[idx+1].getValue(key)
	}
	return n.children[idx].getValue(key)
}

func (n *internalNode) insertChildAt(idx int, child node) {
	lastCnt := len(n.children)
	n.children = append(n.children, nil)
	// move backward
	for i := lastCnt - 1; i >= idx; i-- {
		n.children[i+1] = n.children[i]
	}
	n.children[idx] = child
}
