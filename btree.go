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
	if order <= 1 {
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

	exist, n := t.root.insert(key, value)
	if !exist {
		t.size++
	}

	newRoot := n.split(t.order)
	if newRoot != nil {
		t.height++
		t.root = newRoot
	}
	return nil
}

func (t *Btree) NewIterator(begin, end []byte) Iterator {
	return &iterator{}
}

func (t *Btree) dump() {
	if t.root != nil {
		t.root.dump()
	}
}
