package btree

type btreeHeader struct {
	mbtree   *MVCCBtree
	root     node
	revision uint64
}

func initBtreeHeader(mbtree *MVCCBtree, revision uint64, key, value []byte) *btree {
	tree := &btree{
		mbtree:   mbtree,
		revision: revision,
		order:    mbtree.GetOrder(),
	}

	root := newLeafNode(mbtree, nil, revision)
	root.insertKeyAt(0, key)
	root.insertValueAt(0, value)
	tree.root = root

	return tree
}

func (h *btreeHeader) GetRevision() uint64 {
	return h.revision
}

func (h *btreeHeader) GetOrder() int {
	return h.mbtree.GetOrder()
}

func (h *btreeHeader) Get(key []byte) []byte {
	if h == nil || h.root == nil {
		return nil
	}
	return h.root.getValue(key)
}
