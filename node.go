package btree

import (
	"bytes"
	"fmt"
)

type node struct {
	isLeaf   bool
	numKeys  int
	parent   *node
	children nodes
	keys     [][]byte
	values   [][]byte
}

type nodes []*node

func (n *node) getValue(key []byte) []byte {
	idx := n.findPos(key)
	if n.isLeaf {
		if idx == n.numKeys || bytes.Compare(n.keys[idx], key) != 0 {
			return nil
		}
		return n.values[idx]
	}

	// equal: go to right child
	if idx < n.numKeys && bytes.Compare(n.keys[idx], key) == 0 {
		return n.children[idx+1].getValue(key)
	}

	return n.children[idx].getValue(key)
}

func (n *node) findPos(key []byte) int {
	for i := 0; i < n.numKeys; i++ {
		if bytes.Compare(n.keys[i], key) >= 0 {
			return i
		}
	}

	return n.numKeys
}

func (n *node) insertAt(idx int, key []byte, value []byte) {
	fmt.Printf("(%d, %d) %d", len(n.keys), len(n.values), idx)

	// handle keys
	n.keys = append(n.keys, nil)
	for i := n.numKeys - 1; i >= idx; i-- {
		n.keys[i+1] = n.keys[i]
	}
	n.keys[idx] = key

	if n.isLeaf {
		n.values = append(n.values, nil)
		// move backward
		for i := n.numKeys - 1; i >= idx; i-- {
			n.values[i+1] = n.values[i]
		}
		n.values[idx] = value
	}
}

func (n *node) insert(key, value []byte) (exist bool, modified *node) {
	idx := n.findPos(key)
	if n.isLeaf {
		modified = n

		// replace
		if idx < n.numKeys && bytes.Compare(n.keys[idx], key) == 0 {
			n.values[idx] = value
			exist = true
			return
		}

		// add
		exist = false
		n.insertAt(idx, key, value)
		n.numKeys++
		return
	}

	// equal: go to right child
	if idx < n.numKeys && bytes.Compare(n.keys[idx], key) == 0 {
		exist, modified = n.children[idx+1].insert(key, value)
	} else {
		exist, modified = n.children[idx].insert(key, value)
	}
	return
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
		idx := parent.findPos(midKey)

		parent.children = append(parent.children, nil)
		for i := parent.numKeys; i > idx; i++ {
			parent.children[i+1] = parent.children[i]
		}
		parent.children[idx+1] = rc

		parent.insertAt(idx, midKey, nil)
		parent.numKeys++
	}

	if newRoot != nil {
		return
	} else {
		newRoot = parent.split(order)
		return
	}
}

func (n *node) dump() {
	if n.isLeaf {
		for i := 0; i < n.numKeys; i++ {
			fmt.Printf("(%v, %v) ", string(n.keys[i]), string(n.values[i]))
		}
		fmt.Println()
	} else {
		for i := 0; i < n.numKeys; i++ {
			fmt.Printf("(%v) ", string(n.keys[i]))
		}
		fmt.Println()
	}

	for _, c := range n.children {
		c.dump()
	}
}
