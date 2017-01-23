package btree

import (
	"bytes"
	"testing"
)

func TestBtreePutGet(t *testing.T) {
	btree := NewBtree(4)
	key := []byte("mykey")
	value := []byte("value")
	btree.Put(key, value)

	actualValue := btree.Get(key)
	if bytes.Compare(value, actualValue) != 0 {
		t.Errorf("expected=%v, actual=%v", string(value), string(actualValue))
	}

	btree.Put([]byte("key1"), value)
	btree.Put([]byte("key2"), value)

	// btree.dump()
	writeDot(btree.root, "dotfile")

}
