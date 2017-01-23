package btree

type Iterator interface {
	Next() bool
	Key() []byte
	Value() []byte
	Error() error
}

type iterator struct {
}

func (it *iterator) Next() bool {
	//TODO:
	return false
}

func (it *iterator) Key() []byte {
	//TODO:
	return nil
}

func (it *iterator) Value() []byte {
	//TODO:
	return nil
}

func (it *iterator) Error() error {
	//TODO:
	return nil
}
