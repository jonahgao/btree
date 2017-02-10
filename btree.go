package btree

import (
	"sync"
)

type Iterator interface {
	Next() bool
	Key() []byte
	Value() []byte
	Error() error
}

type Snapshot interface {
	Get(key []byte) []byte
	NewIterator([]byte, []byte) Iterator
	GetRevision() uint64
}

type MVCCBtree struct {
	order           int
	currentTree     *btreeHeader
	currentRevision uint64
	metaLock        sync.RWMutex //TODO: atomic.Value?
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

func (mt *MVCCBtree) Get(key []byte) []byte {
	return mt.getTree().Get(key)
}

func (mt *MVCCBtree) Put(key []byte, value []byte) {
	mt.writeLock.Lock()

	mt.currentRevision++

	oldTree := mt.getTree()
	if oldTree == nil {
		oldTree = &btreeHeader{
			mbtree: mt,
		}
	}
	mt.putTree(oldTree.put(key, value, mt.currentRevision))

	mt.writeLock.Unlock()
}

// TODO: return the deleted value
func (mt *MVCCBtree) Delete(key []byte) {
	mt.writeLock.Lock()
	mt.currentRevision++

	oldTree := mt.getTree()
	mt.putTree(oldTree.delete(key, mt.currentRevision))

	mt.writeLock.Unlock()
}

func (mt *MVCCBtree) GetSnapshot() Snapshot {
	return mt.getTree()
}

func (mt *MVCCBtree) NewIterator(beginKey, endKey []byte) Iterator {
	return mt.getTree().NewIterator(beginKey, endKey)
}

func (mt *MVCCBtree) getTree() *btreeHeader {
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
