package btree

import "bytes"

type node struct {
	isLeaf   bool
	numKeys  int
	parent   *node
	children nodes
	keys     [][]byte
	values   [][]byte
}

type nodes []*node

// find equal or greater key pos. return exist(euqal) and index
func (n *node) findPos(key []byte) (bool, int) {
	for i := 0; i < n.numKeys; i++ {
		c := bytes.Compare(n.keys[i], key)
		if c == 0 {
			return true, i
		} else if c > 0 {
			return false, i
		}
	}

	return false, n.numKeys
}

func (n *node) gotoLeaf(key []byte) *node {
	if n.isLeaf {
		return n
	}

	exist, idx := n.findPos(key)
	// equal: go to right
	if exist {
		return n.children[idx+1].gotoLeaf(key)
	}
	return n.children[idx].gotoLeaf(key)
}

func (n *node) getValue(key []byte) []byte {
	if n.isLeaf {
		exist, idx := n.findPos(key)
		if !exist {
			return nil
		}
		return n.values[idx]
	}

	return n.gotoLeaf(key).getValue(key)
}

func (n *node) insertKeyAt(idx int, key []byte) {
	lastCnt := len(n.keys)
	n.keys = append(n.keys, nil)
	for i := lastCnt - 1; i >= idx; i-- {
		n.keys[i+1] = n.keys[i]
	}
	n.keys[idx] = key
}

func (n *node) insertValueAt(idx int, value []byte) {
	lastCnt := len(n.values)
	n.values = append(n.values, nil)
	// move backward
	for i := lastCnt - 1; i >= idx; i-- {
		n.values[i+1] = n.values[i]
	}
	n.values[idx] = value
}

func (n *node) insertChildAt(idx int, child *node) {
	lastCnt := len(n.children)
	n.children = append(n.children, nil)
	// move backward
	for i := lastCnt - 1; i >= idx; i-- {
		n.children[i+1] = n.children[i]
	}
	n.children[idx] = child
}

// if key is already exist return false, otherwise return true
func (n *node) insert(key, value []byte) (bool, *node) {
	if n.isLeaf {
		exist, idx := n.findPos(key)
		// replace
		if exist {
			n.values[idx] = value
			return false, n
		}
		// add
		n.insertKeyAt(idx, key)
		n.insertValueAt(idx, value)
		n.numKeys++
		return true, n
	}

	return n.gotoLeaf(key).insert(key, value)
}

func (n *node) split(order int) (newRoot *node) {
	if n.numKeys < order {
		return nil
	}

	mid := order / 2
	midKey := n.keys[mid]

	// current is root
	if n.parent == nil {
		newRoot = &node{
			isLeaf: false,
			parent: nil,
		}
		n.parent = newRoot
	}

	// deal right child
	rc := &node{
		isLeaf: n.isLeaf,
		parent: n.parent,
	}
	//TODO: first allcate enough memory then use copy
	if !n.isLeaf {
		for i := mid + 1; i < n.numKeys; i++ {
			rc.keys = append(rc.keys, n.keys[i])
			rc.numKeys++
		}
		for i := mid + 1; i <= n.numKeys; i++ {
			rc.children = append(rc.children, n.children[i])
		}
	} else {
		for i := mid; i < n.numKeys; i++ {
			rc.keys = append(rc.keys, n.keys[i])
			rc.values = append(rc.values, n.values[i])
			rc.numKeys++
		}
	}

	for _, c := range rc.children {
		c.parent = rc
	}

	// deal left child
	n.keys = n.keys[0:mid]
	n.numKeys = mid
	if !n.isLeaf {
		n.children = n.children[0 : mid+1]
	} else {
		n.values = n.values[0:mid]
	}

	// deal parent insert
	parent := n.parent
	if parent.numKeys == 0 {
		parent.numKeys++
		parent.keys = [][]byte{midKey}
		parent.children = []*node{n, rc}
	} else {
		exist, idx := parent.findPos(midKey)
		if exist {
			panic("should does not exist")
		}

		parent.children = append(parent.children, nil)
		for i := parent.numKeys; i > idx; i++ {
			parent.children[i+1] = parent.children[i]
		}
		parent.children[idx+1] = rc

		parent.insertKeyAt(idx, midKey)
		parent.numKeys++
	}

	if newRoot != nil {
		return
	}

	newRoot = parent.split(order)
	return
}

// only for leaf
func (n *node) removeAt(idx int) {
	for i := idx + 1; i < n.numKeys; i++ {
		n.keys[i-1] = n.keys[i]
		n.values[i-1] = n.values[i]
	}
	n.keys = n.keys[:n.numKeys-1]
	n.values = n.keys[:n.numKeys-1]
	n.numKeys--
}

func (n *node) remove(key []byte) (bool, *node) {
	if n.isLeaf {
		exist, idx := n.findPos(key)
		if !exist {
			return false, nil
		}
		n.removeAt(idx)
		return true, n
	}

	return n.gotoLeaf(key).remove(key)
}

// return is tree empty and new root is exist
func (n *node) mergeOrRedistribute(order int) (bool, *node) {
	minKeys := (order+1)/2 - 1
	if n.numKeys >= minKeys {
		return false, nil
	}

	// current is root
	if n.parent == nil {
		if n.numKeys == 0 {
			if n.isLeaf {
				return true, nil
			}

			// only one children
			newRoot := n.children[0]
			newRoot.parent = nil
			return false, newRoot
		}

		return false, nil
	}

	p := n.parent
	// get pos in parent children
	pos := -1
	for i, nod := range p.children {
		if nod == n {
			pos = i
		}
	}
	if pos == -1 {
		panic("child must in parent's children")
	}

	var leftSibling *node
	var rightSibling *node
	if pos != 0 {
		leftSibling = p.children[pos-1]
	}
	if pos != p.numKeys {
		rightSibling = p.children[pos+1]
	}

	if leftSibling != nil && leftSibling.numKeys > minKeys { // borrowFromLeftSibling
		// borrow key from paranet
		borrowKey := p.keys[pos-1]
		if n.isLeaf {
			borrowKey = leftSibling.keys[leftSibling.numKeys-1]
		}

		n.insertKeyAt(0, borrowKey)
		n.numKeys++
		if n.isLeaf {
			n.insertValueAt(0, leftSibling.values[leftSibling.numKeys-1])
		} else {
			n.insertChildAt(0, leftSibling.children[len(leftSibling.children)-1])
		}

		// set parent key to left sibling's last key
		p.keys[pos-1] = leftSibling.keys[leftSibling.numKeys-1]

		// remove left sibling last key, value or children
		leftSibling.keys = leftSibling.keys[:leftSibling.numKeys-1]
		leftSibling.numKeys--
		if leftSibling.isLeaf {
			leftSibling.values = leftSibling.values[:len(leftSibling.values)-1]
		} else {
			leftSibling.children = leftSibling.children[:len(leftSibling.children)-1]
		}

		return p.mergeOrRedistribute(order)
	} else if rightSibling != nil && rightSibling.numKeys > minKeys { // borrowFromRightSibling
		// borrow key from paranet
		borrowKey := p.keys[pos]

		n.insertKeyAt(n.numKeys, borrowKey)
		n.numKeys++
		if n.isLeaf {
			n.values = append(n.values, rightSibling.values[0])
		} else {
			n.children = append(n.children, rightSibling.children[0])
		}

		if n.isLeaf {
			p.keys[pos] = rightSibling.keys[1]
		} else {
			p.keys[pos] = rightSibling.keys[0]
		}

		rightSibling.keys = rightSibling.keys[1:]
		rightSibling.numKeys--
		if rightSibling.isLeaf {
			rightSibling.values = rightSibling.values[1:]
		} else {
			rightSibling.children = rightSibling.children[1:]
		}

		return p.mergeOrRedistribute(order)
	} else if leftSibling != nil { // merge left sibling

	} else { // merge right sibling
	}

	return false, nil
}
