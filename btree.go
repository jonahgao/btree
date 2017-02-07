package btree

type insertResult struct {
}

type deleteResult struct {
}

var gInsertResult insertResult
var gDeleteResult deleteResult

type btree struct {
	mbtree   *MVCCBtree
	order    int
	root     node
	revision uint64
}

func initBtree(mbtree *MVCCBtree, revision uint64, key, value []byte) *btree {
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

func (t *btree) getRevision() uint64 {
	return t.revision
}

func (t *btree) Get(key []byte) []byte {
	if t == nil || t.root == nil {
		return nil
	}
	return t.root.getValue(key)
}
