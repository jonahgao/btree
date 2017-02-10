package btree

type leafNode struct {
	baseNode
	values [][]byte
}

func newLeafNode(t *MVCCBtree, r uint64) *leafNode {
	return &leafNode{
		baseNode: baseNode{
			tree:     t,
			revision: r,
			keys:     make([][]byte, 0, t.order-1),
		},
		values: make([][]byte, 0, t.order-1),
	}
}

func (n *leafNode) isLeaf() bool {
	return true
}

func (n *leafNode) leftMostKey() []byte {
	return n.keys[0]
}

func (n *leafNode) get(key []byte) []byte {
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
	return newLeaf
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
		return n.replaceValue(pos, value, revision)
	} else if len(n.keys)+1 <= n.maxKeys() {
		return n.addValue(pos, key, value, revision)
	} else {
		return n.addAndSplit(pos, key, value, revision)
	}
}

func (n *leafNode) removeValue(pos int, revision uint64) *deleteResult {
	newLeaf := newLeafNode(n.tree, revision)
	for i := 0; i < pos; i++ {
		newLeaf.keys = append(newLeaf.keys, n.keys[i])
		newLeaf.values = append(newLeaf.values, n.values[i])
	}
	for i := pos + 1; i < len(n.keys); i++ {
		newLeaf.keys = append(newLeaf.keys, n.keys[i])
		newLeaf.values = append(newLeaf.values, n.values[i])
	}

	return &deleteResult{
		rtype:    dRTypeRemoved,
		modified: newLeaf,
	}
}

func (n *leafNode) borrowFromLeft(pos int, revision uint64, sibling *leafNode) *deleteResult {
	newLeaf := newLeafNode(n.tree, revision)
	newSibling := newLeafNode(n.tree, revision)

	// copy to new sibling
	for i := 0; i < len(sibling.keys)-1; i++ {
		newSibling.keys = append(newSibling.keys, sibling.keys[i])
		newSibling.values = append(newSibling.values, sibling.values[i])
	}

	// insert borrowed kv
	newLeaf.keys = append(newLeaf.keys, sibling.keys[sibling.numOfKeys()-1])
	newLeaf.values = append(newLeaf.values, sibling.values[len(sibling.keys)-1])
	// copy kv before pos from n
	for i := 0; i < pos; i++ {
		newLeaf.keys = append(newLeaf.keys, n.keys[i])
		newLeaf.values = append(newLeaf.values, n.values[i])
	}
	// copy kv after pos from n
	for i := pos + 1; i < len(n.keys); i++ {
		newLeaf.keys = append(newLeaf.keys, n.keys[i])
		newLeaf.values = append(newLeaf.values, n.values[i])
	}

	return &deleteResult{
		rtype:           dRTypeBorrowFromLeft,
		modified:        newLeaf,
		modifiedSibling: newSibling,
	}
}

func (n *leafNode) borrowFromRight(pos int, revision uint64, sibling *leafNode) *deleteResult {
	newSibling := newLeafNode(n.tree, revision)
	// copy kv to new sibling
	for i := 1; i < len(sibling.keys); i++ {
		newSibling.keys = append(newSibling.keys, sibling.keys[i])
		newSibling.values = append(newSibling.values, sibling.values[i])
	}

	newLeaf := newLeafNode(n.tree, revision)
	// copy kv before pos from n
	for i := 0; i < pos; i++ {
		newLeaf.keys = append(newLeaf.keys, n.keys[i])
		newLeaf.values = append(newLeaf.values, n.values[i])
	}
	// copy kv after pos from n
	for i := pos + 1; i < len(n.keys); i++ {
		newLeaf.keys = append(newLeaf.keys, n.keys[i])
		newLeaf.values = append(newLeaf.values, n.values[i])
	}
	// insert borrowed kv
	newLeaf.keys = append(newLeaf.keys, sibling.keys[0])
	newLeaf.values = append(newLeaf.values, sibling.values[0])

	return &deleteResult{
		rtype:           dRTypeBorrowFromRight,
		modified:        newLeaf,
		modifiedSibling: newSibling,
	}
}

func (n *leafNode) mergeWithLeft(pos int, revision uint64, sibling *leafNode) *deleteResult {
	newLeaf := newLeafNode(n.tree, revision)
	for i := 0; i < len(sibling.keys); i++ {
		newLeaf.keys = append(newLeaf.keys, sibling.keys[i])
		newLeaf.values = append(newLeaf.values, sibling.values[i])
	}
	for i := 0; i < pos; i++ {
		newLeaf.keys = append(newLeaf.keys, n.keys[i])
		newLeaf.values = append(newLeaf.values, n.values[i])
	}
	for i := pos + 1; i < len(n.keys); i++ {
		newLeaf.keys = append(newLeaf.keys, n.keys[i])
		newLeaf.values = append(newLeaf.values, n.values[i])
	}

	return &deleteResult{
		rtype:    dRTypeMergeWithLeft,
		modified: newLeaf,
	}
}

func (n *leafNode) mergeWithRight(pos int, revision uint64, sibling *leafNode) *deleteResult {
	newLeaf := newLeafNode(n.tree, revision)
	for i := 0; i < pos; i++ {
		newLeaf.keys = append(newLeaf.keys, n.keys[i])
		newLeaf.values = append(newLeaf.values, n.values[i])
	}
	for i := pos + 1; i < len(n.keys); i++ {
		newLeaf.keys = append(newLeaf.keys, n.keys[i])
		newLeaf.values = append(newLeaf.values, n.values[i])
	}
	for i := 0; i < len(sibling.keys); i++ {
		newLeaf.keys = append(newLeaf.keys, sibling.keys[i])
		newLeaf.values = append(newLeaf.values, sibling.values[i])
	}

	return &deleteResult{
		rtype:    dRTypeMergeWithRight,
		modified: newLeaf,
	}
}

func (n *leafNode) delete(key []byte, revision uint64, parent node, parentPos int) *deleteResult {
	exist, pos := n.findPos(key)
	if !exist {
		return &deleteResult{
			rtype: dRTypeNotPresent,
		}
	}

	// current node is root or node's keys is adeuate
	if parent == nil || len(n.keys) > n.minKeys() {
		return n.removeValue(pos, revision)
	}

	siblingPos := n.selectSibling(parent, parentPos)
	sibling := ((parent.(*internalNode)).childAt(siblingPos)).(*leafNode)
	if sibling.numOfKeys() > n.minKeys() { // can borrow
		if siblingPos < parentPos {
			return n.borrowFromLeft(pos, revision, sibling)
		} else {
			return n.borrowFromRight(pos, revision, sibling)
		}
	} else {
		if siblingPos < parentPos {
			return n.mergeWithLeft(pos, revision, sibling)
		} else {
			return n.mergeWithRight(pos, revision, sibling)
		}
	}
}

func (n *leafNode) iterateNext(*iterator) bool {
	return false
}
