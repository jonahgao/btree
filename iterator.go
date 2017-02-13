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
	return it.stackTop().node.iterateNext(it)
}

func (it *iterator) stackTop() iteratorPos {
	return it.stack[len(it.stack)-1]
}

func (it *iterator) stackPop() iteratorPos {
	top := it.stack[len(it.stack)-1]
	it.stack = it.stack[:len(it.stack)-1]
	return top
}

func (it *iterator) stackPush(pos iteratorPos) {
	it.stack = append(it.stack, pos)
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
