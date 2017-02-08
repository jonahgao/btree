package btree

type leafNode struct {
	baseNode
	values [][]byte
}

func newLeafNode(t *MVCCBtree, r uint64) *leafNode {
	return &leafNode{
		tree:     t,
		revision: r,
		keys:     make([][]byte, 0, t.order-1),
		values:   make([][]byte, 0, t.order-1),
	}
}

func (n *leafNode) isLeaf() bool {
	return true
}

func (n *leafNode) get([]byte) []byte {
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
	newLeaf := newLeafNode(n.tree, revision)
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
	newLeaf := newLeafNode(n.tree, revision)
	for i := 0; i < pos; i++ {
		newLeaf.keys = append(newLeaf.keys, n.keys[i])
		newLeaf.values = append(newLeaf.values, n.values[i])
	}
	newLeaf.keys = append(newLeaf.keys, key)
	newLeaf.values = append(newLeaf.values, value)
	for i := pos; i < len(n.keys); i++ {
		newLeaf.keys = append(newLeaf.keys, n.keys[i])
		newLeaf.values = append(newLeaf.values, n.values[i])
	}

	return &insertResult{
		rtype:    iRTypeModified,
		modified: newLeaf,
	}
}

func (n *leafNode) addAndSplit(pos int, key, value []byte, revision uint64) *insertResult {
	leftLeaf := newLeafNode(n.tree, revision)
	rightLeaf := newLeafNode(n.tree, revision)
	middle := n.splitPivot()
	if pos <= middle {
		for i := 0; i < pos; i++ {
			leftLeaf.keys = append(leftLeaf.keys, n.keys[i])
			leftLeaf.values = append(leftLeaf.values, n.values[i])
		}
		leftLeaf.keys = append(leftLeaf.keys, key)
		leftLeaf.values = append(leftLeaf.values, value)
		for i := pos; i < middle; i++ {
			leftLeaf.keys = append(leftLeaf.keys, n.keys[i])
			leftLeaf.values = append(leftLeaf.values, n.values[i])
		}

		for i := middle; i < len(n.keys); i++ {
			rightLeaf.keys = append(rightLeaf.keys, n.keys[i])
			rightLeaf.values = append(rightLeaf.values, n.values[i])
		}
	} else {
		for i := 0; i < middle; i++ {
			leftLeaf.keys = append(leftLeaf.keys, n.keys[i])
			leftLeaf.values = append(leftLeaf.values, n.values[i])
		}

		for i := middle; i < pos; i++ {
			rightLeaf.keys = append(rightLeaf.keys, n.keys[i])
			rightLeaf.values = append(rightLeaf.values, n.values[i])
		}
		rightLeaf.keys = append(rightLeaf.keys, key)
		rightLeaf.values = append(rightLeaf.values, value)
		for i := pos; i < len(n.keys); i++ {
			rightLeaf.keys = append(rightLeaf.keys, n.keys[i])
			rightLeaf.values = append(rightLeaf.values, n.values[i])
		}
	}

	return &insertResult{
		rtype: iRTypeSplit,
		left:  leftLeaf,
		right: rightLeaf,
		pivot: rightLeaf.keys[0],
	}
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
