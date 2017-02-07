package btree

type MVCCBtree struct {
	order           int
	trees           map[uint64]node
	currentRevision uint64
}

func NewMVCCBtree(order int) *MVCCBtree {
	if order <= 2 {
		panic("bad order")
	}

	return &MVCCBtree{
		order:           order,
		trees:           make(map[uint64]node, 64),
		currentRevision: 0,
	}
}

func (t *MVCCBtree) Get(key []byte) []byte {

}

func (t *MVCCBtree) Put(key []byte, value []byte) {

}

func (t *MVCCBtree) Delete(key []byte) {

}
