package btree

type internalNode struct {
	baseNode
	children []node
}

func newInternalNode(t *MVCCBtree, r uint64) *internalNode {
	return &internalNode{
		baseNode: baseNode{
			tree:     t,
			revision: r,
			keys:     make([][]byte, 0, t.order-1),
		},
		children: make([]node, 0, t.order),
	}
}

func (n *internalNode) isLeaf() bool {
	return false
}

func (n *internalNode) leftMostKey() []byte {
	return n.children[0].leftMostKey()
}

func (n *internalNode) get(key []byte) []byte {
	exist, idx := n.findPos(key)
	// equal: go to right
	if exist {
		return n.children[idx+1].get(key)
	}
	return n.children[idx].get(key)
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

func (n *internalNode) childAt(idx int) node {
	return n.children[idx]
}

func (n *internalNode) clone(revision uint64) *internalNode {
	newINode := newInternalNode(n.tree, revision)
	for _, k := range n.keys {
		newINode.keys = append(newINode.keys, k)
	}
	for _, c := range n.children {
		newINode.children = append(newINode.children, c)
	}
	return newINode
}

func (n *internalNode) replaceChild(pos int, childIResult *insertResult, revision uint64) *insertResult {
	newNode := n.clone(revision)
	newNode.children[pos] = childIResult.modified
	return &insertResult{
		rtype:    iRTypeModified,
		modified: newNode,
	}
}

func (n *internalNode) insertChild(pos int, childIResult *insertResult, revision uint64) *insertResult {
	newNode := newInternalNode(n.tree, revision)

	// inset key
	for i := 0; i < pos; i++ {
		newNode.keys = append(newNode.keys, n.keys[i])
	}
	newNode.keys = append(newNode.keys, childIResult.pivot)
	for i := pos; i < len(n.keys); i++ {
		newNode.keys = append(newNode.keys, n.keys[i])
	}

	// insert child
	for i := 0; i < pos; i++ {
		newNode.children = append(newNode.children, n.children[i])
	}
	newNode.children = append(newNode.children, childIResult.left)
	newNode.children = append(newNode.children, childIResult.right)
	for i := pos + 1; i < len(n.children); i++ {
		newNode.children = append(newNode.children, n.children[i])
	}

	return &insertResult{
		rtype:    iRTypeModified,
		modified: newNode,
	}
}

func (n *internalNode) addAndSplit(pos int, childIResult *insertResult, revision uint64) *insertResult {
	leftNode := newInternalNode(n.tree, revision)
	rightNode := newInternalNode(n.tree, revision)
	middle := n.splitPivot()
	var pivot []byte
	if pos < middle {
		// left
		for i := 0; i < pos; i++ {
			leftNode.keys = append(leftNode.keys, n.keys[i])
			leftNode.children = append(leftNode.children, n.children[i])
		}
		leftNode.keys = append(leftNode.keys, childIResult.pivot)
		leftNode.children = append(leftNode.children, childIResult.left)
		leftNode.children = append(leftNode.children, childIResult.right)
		for i := pos; i < middle-1; i++ {
			leftNode.keys = append(leftNode.keys, n.keys[i])
		}
		for i := pos + 1; i < middle; i++ {
			leftNode.children = append(leftNode.children, n.children[i])
		}

		// right
		for i := middle; i < len(n.keys); i++ {
			rightNode.keys = append(rightNode.keys, n.keys[i])
		}
		for i := middle; i < len(n.children); i++ {
			rightNode.children = append(rightNode.children, n.children[i])
		}

		pivot = n.keys[middle-1]
	} else if pos == middle {
		for i := 0; i < middle; i++ {
			leftNode.keys = append(leftNode.keys, n.keys[i])
			leftNode.children = append(leftNode.children, n.children[i])
		}
		leftNode.children = append(leftNode.children, childIResult.left)

		// right
		for i := middle; i < len(n.keys); i++ {
			rightNode.keys = append(rightNode.keys, n.keys[i])
		}
		rightNode.children = append(rightNode.children, childIResult.right)
		for i := middle + 1; i < len(n.children); i++ {
			rightNode.children = append(rightNode.children, n.children[i])
		}

		pivot = childIResult.pivot
	} else {
		// left
		for i := 0; i < middle; i++ {
			leftNode.keys = append(leftNode.keys, n.keys[i])
		}
		for i := 0; i < middle+1; i++ {
			leftNode.children = append(leftNode.children, n.children[i])
		}

		// right
		for i := middle + 1; i < pos; i++ {
			rightNode.keys = append(rightNode.keys, n.keys[i])
			rightNode.children = append(rightNode.children, n.children[i])
		}
		rightNode.keys = append(rightNode.keys, childIResult.pivot)
		rightNode.children = append(rightNode.children, childIResult.left)
		rightNode.children = append(rightNode.children, childIResult.right)

		for i := pos; i < len(n.keys); i++ {
			rightNode.keys = append(rightNode.keys, n.keys[i])
		}
		for i := pos + 1; i < len(n.children); i++ {
			rightNode.children = append(rightNode.children, n.children[i])
		}

		pivot = n.keys[middle]
	}

	return &insertResult{
		rtype: iRTypeSplit,
		left:  leftNode,
		right: rightNode,
		pivot: pivot,
	}
}

func (n *internalNode) insert(key, value []byte, revision uint64) *insertResult {
	exist, pos := n.findPos(key)
	if exist {
		pos++
	}
	childIResult := n.children[pos].insert(key, value, revision)
	if childIResult.rtype == iRTypeModified { // modified result
		return n.replaceChild(pos, childIResult, revision)
	} else { // split result
		if len(n.keys) < n.maxKeys() {
			return n.insertChild(pos, childIResult, revision)
		} else {
			return n.addAndSplit(pos, childIResult, revision)
		}
	}
}

func (n *internalNode) handleRemovedResult(childResult *deleteResult, pos int, exist bool, revision uint64) *deleteResult {
	newNode := n.clone(revision)
	newNode.children[pos] = childResult.modified
	// the deleted key is in our node's keys, update to right child's left most
	if exist {
		newNode.keys[pos-1] = childResult.modified.leftMostKey()
	}
	return &deleteResult{
		rtype:    dRTypeRemoved,
		modified: newNode,
	}
}

func (n *internalNode) handleBorrowedResult(childResult *deleteResult, pos int, exist bool, revision uint64) *deleteResult {
	newInternalNode := n.clone(revision)
	if exist {
		if childResult.rtype == dRTypeBorrowFromLeft {
			newInternalNode.keys[pos-1] = childResult.modified.leftMostKey()
			newInternalNode.children[pos] = childResult.modified
			newInternalNode.children[pos-1] = childResult.modifiedSibling
		} else {
			newInternalNode.keys[pos-1] = childResult.modified.leftMostKey()
			newInternalNode.keys[pos] = childResult.modifiedSibling.leftMostKey()
			newInternalNode.children[pos] = childResult.modified
			newInternalNode.children[pos+1] = childResult.modifiedSibling
		}
	} else {
		if childResult.rtype == dRTypeBorrowFromLeft {
			newInternalNode.keys[pos-1] = childResult.modified.leftMostKey()
			newInternalNode.children[pos-1] = childResult.modifiedSibling
			newInternalNode.children[pos] = childResult.modified
		} else {
			newInternalNode.keys[pos-1] = childResult.modifiedSibling.leftMostKey()
			newInternalNode.children[pos] = childResult.modifiedSibling
			newInternalNode.children[pos+1] = childResult.modified
		}
	}
	return &deleteResult{
		rtype:    dRTypeRemoved,
		modified: newInternalNode,
	}
}

// append current node's key and children to new internal node
// And apply the child merge result to the new internal node
func (n *internalNode) childMergeAppendTo(newNode *internalNode, childResult *deleteResult, pos int, exist bool) {
	// copy and append children, skip two children; insert the modified one
	skipChildPosStart := -1
	if childResult.rtype == dRTypeMergeWithLeft {
		skipChildPosStart = pos - 1
	} else if childResult.rtype == dRTypeMergeWithRight {
		skipChildPosStart = pos
	} else {
		panic("unexpect merge type")
	}
	for i := 0; i < skipChildPosStart; i++ {
		newNode.children = append(newNode.children, n.children[i])
	}
	newNode.children = append(newNode.children, childResult.modified)
	for i := skipChildPosStart + 2; i < len(n.children); i++ {
		newNode.children = append(newNode.children, n.children[i])
	}

	// copy and append keys, skip the merged key; update key before skip pos if necessary
	skipKeyPos := -1
	if childResult.rtype == dRTypeMergeWithLeft {
		skipKeyPos = pos
	} else if childResult.rtype == dRTypeMergeWithRight {
		skipKeyPos = pos + 1
	} else {
		panic("unexpected merge type")
	}
	if exist {
		skipKeyPos--
	}
	for i := 0; i < skipKeyPos; i++ {
		newNode.keys = append(newNode.keys, n.keys[i])
	}
	if exist && childResult.rtype == dRTypeMergeWithRight {
		newNode.keys[len(newNode.keys)-1] = childResult.modified.leftMostKey()
	}
	for i := skipKeyPos + 1; i < len(n.keys[i]); i++ {
		newNode.keys = append(newNode.keys, n.keys[i])
	}
}

func (n *internalNode) removeKey(childResult *deleteResult, pos int, exist bool, revision uint64) *deleteResult {
	newNode := newInternalNode(n.tree, revision)
	n.childMergeAppendTo(newNode, childResult, pos, exist)
	return &deleteResult{
		rtype:    dRTypeRemoved,
		modified: newNode,
	}
}

func (n *internalNode) borrowFromLeft(childResult *deleteResult, pos int, exist bool, revision uint64, sibling *internalNode) *deleteResult {
	newNode := newInternalNode(n.tree, revision)
	newSibling := newInternalNode(n.tree, revision)

	// handle new sibling
	for i := 0; i < len(sibling.keys)-1; i++ {
		newSibling.keys = append(newSibling.keys, sibling.keys[i])
	}
	for i := 0; i < len(sibling.children)-1; i++ {
		newSibling.children = append(newSibling.children, sibling.children[i])
	}

	// handle new node children
	newNode.keys = append(newNode.keys, nil) // placeholder
	newNode.children = append(newNode.children, sibling.children[len(sibling.children)-1])
	n.childMergeAppendTo(newNode, childResult, pos, exist)
	// the borrowed key must be right child's left most
	newNode.keys[0] = newNode.children[1].leftMostKey()

	return &deleteResult{
		rtype:           dRTypeBorrowFromLeft,
		modified:        newNode,
		modifiedSibling: newSibling,
	}
}

func (n *internalNode) borrowFromRight(childResult *deleteResult, pos int, exist bool, revision uint64, sibling *internalNode) *deleteResult {
	newNode := newInternalNode(n.tree, revision)
	newSibling := newInternalNode(n.tree, revision)

	// handle new sibling
	for i := 1; i < len(sibling.keys); i++ {
		newSibling.keys = append(newSibling.keys, sibling.keys[i])
	}
	for i := 1; i < len(sibling.children); i++ {
		newSibling.children = append(newSibling.children, sibling.children[i])
	}

	n.childMergeAppendTo(newNode, childResult, pos, exist)
	newNode.children = append(newNode.children, sibling.children[0])
	newNode.keys = append(newNode.keys, sibling.children[0].leftMostKey()) // the borrowed key must be right child's left most

	return &deleteResult{
		rtype:           dRTypeBorrowFromRight,
		modified:        newNode,
		modifiedSibling: newSibling,
	}
}

func (n *internalNode) mergeWithLeft(childResult *deleteResult, pos int, exist bool, revision uint64, sibling *internalNode) *deleteResult {
	newNode := newInternalNode(n.tree, revision)

	// add left sibling's children
	for i := 0; i < len(sibling.children); i++ {
		newNode.children = append(newNode.children, sibling.children[i])
	}
	// add left sibling's keys
	for i := 0; i < len(sibling.keys); i++ {
		newNode.keys = append(newNode.keys, sibling.keys[i])
	}
	// middle Key, placeholder
	newNode.keys = append(newNode.keys, nil)

	n.childMergeAppendTo(newNode, childResult, pos, exist)

	// set midkey
	newNode.keys[len(sibling.keys)] = newNode.children[len(sibling.keys)+1].leftMostKey()

	return &deleteResult{
		rtype:    dRTypeMergeWithLeft,
		modified: newNode,
	}
}

func (n *internalNode) mergeWithRight(childResult *deleteResult, pos int, exist bool, revision uint64, sibling *internalNode) *deleteResult {
	newNode := newInternalNode(n.tree, revision)

	n.childMergeAppendTo(newNode, childResult, pos, exist)

	// add right sibling's children
	for i := 0; i < len(sibling.children); i++ {
		newNode.children = append(newNode.children, sibling.children[i])
	}

	// middle key, placeholder
	placeholderPos := len(newNode.keys)
	newNode.keys = append(newNode.keys, nil)

	// add right sibling's keys
	for i := 0; i < len(sibling.keys); i++ {
		newNode.keys = append(newNode.keys, sibling.keys[i])
	}

	// update middle key
	newNode.keys[placeholderPos] = newNode.children[placeholderPos+1].leftMostKey()

	return &deleteResult{
		rtype:    dRTypeMergeWithRight,
		modified: newNode,
	}
}

func (n *internalNode) delete(key []byte, revision uint64, parent node, parentPos int) *deleteResult {
	exist, pos := n.findPos(key)
	if exist {
		pos++
	}

	childDResult := n.children[pos].delete(key, revision, n, pos)
	// not present
	if childDResult.rtype == dRTypeNotPresent {
		return childDResult
	} else if childDResult.rtype == dRTypeRemoved {
		return n.handleRemovedResult(childDResult, pos, exist, revision)
	} else if childDResult.rtype == dRTypeBorrowFromLeft || childDResult.rtype == dRTypeBorrowFromRight {
		return n.handleBorrowedResult(childDResult, pos, exist, revision)
	} else { // merge
		if parent == nil || len(n.keys) > n.minKeys() {
			return n.removeKey(childDResult, pos, exist, revision)
		} else {
			siblingPos := n.selectSibling(parent, parentPos)
			sibling := (parent.(*internalNode)).children[siblingPos].(*internalNode)
			if sibling.numOfKeys() > n.minKeys() {
				if siblingPos < parentPos {
					return n.borrowFromLeft(childDResult, pos, exist, revision, sibling)
				} else {
					return n.borrowFromRight(childDResult, pos, exist, revision, sibling)
				}
			} else {
				if siblingPos < parentPos {
					return n.mergeWithLeft(childDResult, pos, exist, revision, sibling)
				} else {
					return n.mergeWithRight(childDResult, pos, exist, revision, sibling)
				}
			}
		}
	}
	return nil
}
