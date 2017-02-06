package btree

import (
	"bytes"
	"fmt"
	"testing"
)

var testDotExePath = "dot"

func TestBtreePutGet(t *testing.T) {
	m := 4
	n := 20

	btree := NewBtree(m)
	for i := 1; i <= n; i++ {
		key := []byte(fmt.Sprintf("%04d", i))
		value := []byte(fmt.Sprintf("%04d", i))
		btree.Put(key, value)
	}

	for i := 1; i <= n; i++ {
		key := []byte(fmt.Sprintf("%04d", i))
		expectedValue := []byte(fmt.Sprintf("%04d", i))
		actualValue := btree.Get(key)
		if bytes.Compare(expectedValue, actualValue) != 0 {
			t.Errorf("expected=%v, actual=%v", string(expectedValue), string(actualValue))
		}
	}

	err := writeDotSvg(testDotExePath, "output.svg", btree)
	if err != nil {
		t.Error(err)
	}

	btree.Delete([]byte("0013"))
	err = writeDotSvg(testDotExePath, "output2.svg", btree)
	if err != nil {
		t.Error(err)
	}

	btree.Delete([]byte("0014"))
	err = writeDotSvg(testDotExePath, "output3.svg", btree)
	if err != nil {
		t.Error(err)
	}
}
