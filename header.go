package btree

type btreeHeader struct {
	mbtree   *MVCCBtree
	root     node
	revision uint64
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
	return h.root.get(key)
}

func (h *btreeHeader) put(key, value []byte, revision uint64) *btreeHeader {
	if h == nil || h.root == nil {
		root := newLeafNode(h.mbtree, revision)
		root.insertKeyAt(0, key)
		root.insertValueAt(0, value)

		return &btreeHeader{
			mbtree:   h.mbtree,
			root:     root,
			revision: revision,
		}
	}

	newHeader := &btreeHeader{
		mbtree:   h.mbtree,
		revision: revision,
	}
	insertResult := h.root.insert(key, value, revision)
	if insertResult.rtype == iRTypeModified {
		newHeader.root = insertResult.modified
	} else {
		root := newInternalNode(h.mbtree, revision)
		root.keys = append(root.keys, insertResult.pivot)
		root.children = append(root.children, insertResult.left)
		root.children = append(root.children, insertResult.right)
		newHeader.root = root
	}
	return newHeader
}
