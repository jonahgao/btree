package btree

import (
	"fmt"
	"testing"
)

var testDotExePath = "D:\\Develop\\graphviz-2.38\\bin\\dot.exe"

func TestBtreePutGet(t *testing.T) {
	btree := NewBtree(4)

	for i := 1; i <= 9; i++ {
		key := []byte(fmt.Sprintf("%02d", i))
		value := []byte(fmt.Sprintf("%02d", i))
		btree.Put(key, value)

	}

	// btree.dump()

	err := writeDotSvg(testDotExePath, "output.svg", btree)
	if err != nil {
		t.Error(err)
	}
}
