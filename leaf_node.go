package btree

type leafNode struct {
	baseNode
	values [][]byte
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
