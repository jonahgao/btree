package btree

//Btree btree
type Btree struct {
	order  int
	root   *node
	size   int // nums of kv pair
	height int // tree height
}

// NewBtree create new btree
func NewBtree(order int) *Btree {
	if order <= 2 {
		panic("bad order")
	}

	return &Btree{
		order: order,
	}
}

func (t *Btree) Get(key []byte) []byte {
	if t.root == nil {
		return nil
	}
	return t.root.getValue(key)
}

func (t *Btree) Put(key []byte, value []byte) error {
	if t.root == nil {
		t.root = &node{
			isLeaf:  true,
			numKeys: 1,
			parent:  nil,
			keys:    [][]byte{key},
			values:  [][]byte{value},
		}
		t.size++
		t.height++
		return nil
	}

	b, n := t.root.insert(key, value)
	if b {
		t.size++
		// handle split
		newRoot := n.split(t.order)
		if newRoot != nil {
			t.height++
			t.root = newRoot
		}
	}

	return nil
}

func (t *Btree) Delete(key []byte) error {
	if t.root != nil {
		b, n := t.root.remove(key)
		if b {
			t.size--
			empty, newRoot := n.mergeOrRedistribute(t.order)
			if empty {
				t.root = nil
				t.height--
			} else if newRoot != nil {
				t.root = newRoot
				t.height--
			}
		}
	}
	return nil
}

func (t *Btree) NewIterator(begin, end []byte) Iterator {
	return &iterator{}
}
