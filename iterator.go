package btree

type iteratorPos struct {
	node node
	pos  int
}

type iterator struct {
	beginKey     []byte
	endKey       []byte
	stack        []iteratorPos
	currentKey   []byte
	currentValue []byte
}

func (it *iterator) Next() bool {
	if len(it.stack) == 0 {
		return false
	}
	return it.stack[len(it.stack)-1].node.iterateNext(it)
}

func (it *iterator) Key() []byte {
	return it.currentKey
}

func (it *iterator) Value() []byte {
	return it.currentValue
}

func (it *iterator) Error() error {
	return nil
}
