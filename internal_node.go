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
		if len(n.keys)+1 <= n.maxKeys() {
			return n.insertChild(pos, childIResult, revision)
		} else {
			return n.addAndSplit(pos, childIResult, revision)
		}
	}
}
