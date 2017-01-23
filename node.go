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

type nodes []node

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
		n.keys = append(n.keys, nil)
		n.values = append(n.values, nil)
		// move backward
		for i := n.numKeys - 1; i >= idx; i-- {
			n.keys[i+1] = n.keys[i]
			n.values[i+1] = n.values[i]
		}
		n.keys[idx] = key
		n.values[idx] = value
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

func (n *node) split(order int) {
	if n.numKeys < order {
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

	for _, c := range []node(n.children) {
		c.dump()
	}
}
