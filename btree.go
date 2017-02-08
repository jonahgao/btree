package btree

import (
	"sync"
)

type Snapshot interface {
	Get(key []byte) []byte
}

type Iterator interface {
	Next() bool
	Key() []byte
	Value() []byte
	Error() error
}

type MVCCBtree struct {
	order           int
	currentTree     *btreeHeader
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

func (mt *MVCCBtree) GetTree() *btreeHeader {
	mt.metaLock.RLock()
	t := mt.currentTree
	mt.metaLock.RUnlock()
	return t
}

func (mt *MVCCBtree) putTree(bt *btreeHeader) {
	mt.metaLock.Lock()
	if mt.currentTree == nil || bt.GetRevision() > mt.currentTree.GetRevision() {
		mt.currentTree = bt
	}
	mt.metaLock.Unlock()
}

func (mt *MVCCBtree) GetSnapshot() Snapshot {
	return mt.GetTree()
}

func (mt *MVCCBtree) Get(key []byte) []byte {
	return mt.GetTree().Get(key)
}

func (mt *MVCCBtree) Put(key []byte, value []byte) {
	mt.writeLock.Lock()

	mt.currentRevision++

	oldTree := mt.GetTree()
	if oldTree == nil {
		oldTree = &btreeHeader{
			mbtree: mt,
		}
	}
	mt.putTree(oldTree.put(key, value, mt.currentRevision))

	mt.writeLock.Unlock()
}

func (mt *MVCCBtree) Delete(key []byte) {
	mt.writeLock.Lock()
	defer mt.writeLock.Unlock()
}
