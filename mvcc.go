package btree

import (
	"sync"
)

type MVCCBtree struct {
	order           int
	currentTree     *btree
	currentRevision uint64
	metaLock        sync.RWMutex
	writeLock       sync.Mutex
}

func NewMVCCBtree(order int) *MVCCBtree {
	if order <= 2 {
		panic("bad order")
	}

	return &MVCCBtree{
		order: order,
	}
}

func (mt *MVCCBtree) GetOrder() int {
	return mt.order
}

func (mt *MVCCBtree) GetTree() *btree {
	mt.metaLock.RLock()
	t := mt.currentTree
	mt.metaLock.RUnlock()
	return t
}

func (mt *MVCCBtree) putTree(bt *btree) {
	mt.metaLock.Lock()
	if mt.currentTree == nil || bt.getRevision() > mt.currentTree.getRevision() {
		mt.currentTree = bt
	}
	mt.metaLock.Unlock()
}

func (mt *MVCCBtree) Get(key []byte) []byte {
	return mt.GetTree().Get(key)
}

func (mt *MVCCBtree) Put(key []byte, value []byte) {
	mt.writeLock.Lock()
	defer mt.writeLock.Unlock()

	mt.currentRevision++

	oldTree := mt.GetTree()
	if oldTree == nil {
		newTree := initBtree(mt, mt.currentRevision, key, value)
		mt.putTree(newTree)
	} else {

	}
}

func (mt *MVCCBtree) Delete(key []byte) {
	mt.writeLock.Lock()
	defer mt.writeLock.Unlock()
}
