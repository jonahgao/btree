package btree

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var testDotExePath = "dot"

func TestBtreeRandPutGet(t *testing.T) {
	m := 5
	n := 10

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	results := make(map[string]string)

	btree := NewMVCCBtree(m)
	for i := 1; i <= n; i++ {
		r := rnd.Int() % 10000
		key := []byte(fmt.Sprintf("%06d", r))
		value := []byte(fmt.Sprintf("Value%06d", r))
		results[string(key)] = string(value)
		btree.Put(key, value)
		fmt.Printf("Put %d\n", r)
		writeDotSvg(testDotExePath, fmt.Sprintf("%02d.svg", i), btree, fmt.Sprintf("Insert %v:", string(key)))
	}

	writeDotSvg(testDotExePath, "output.svg", btree, "")

	for k, v := range results {
		value := btree.Get([]byte(k))
		if string(value) != v {
			t.Errorf("expected=%v, actual=%v", v, string(value))
		}
	}
}

func TestBtreeDelete(t *testing.T) {
	m := 3
	n := 50

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	results := make(map[string]string)

	btree := NewMVCCBtree(m)
	for i := 1; i <= n; i++ {
		r := rnd.Int() % 100
		key := []byte(fmt.Sprintf("%02d", r))
		value := []byte(fmt.Sprintf("Value%02d", r))
		results[string(key)] = string(value)
		btree.Put(key, value)
		fmt.Printf("Put %d\n", r)
	}

	writeDotSvg(testDotExePath, "output.svg", btree, "")

	for k, v := range results {
		value := btree.Get([]byte(k))
		if string(value) != v {
			t.Errorf("expected=%v, actual=%v", v, string(value))
		}
	}

	idx := 0
	for key := range results {
		t.Logf("Delete %v", key)
		btree.Delete([]byte(key))
		idx++
		// writeDotSvg(testDotExePath, fmt.Sprintf("output%02d.svg", idx), btree, fmt.Sprintf("Delete %v:", key))
	}
}

func TestBtreePutGet(t *testing.T) {
	m := 4
	n := 20

	btree := NewMVCCBtree(m)

	for i := 1; i <= n; i++ {
		key := []byte(fmt.Sprintf("%04d", i))
		value := []byte(fmt.Sprintf("Value%04d", i))
		btree.Put(key, value)
	}

	for i := 1; i <= n; i++ {
		key := []byte(fmt.Sprintf("%04d", i))
		expectedValue := []byte(fmt.Sprintf("Value%04d", i))
		actualValue := btree.Get(key)
		if bytes.Compare(expectedValue, actualValue) != 0 {
			t.Errorf("expected=%v, actual=%v", string(expectedValue), string(actualValue))
		}
	}

	writeDotSvg(testDotExePath, "output.svg", btree, "")
}

func TestBtreeIterator(t *testing.T) {
	m := 5
	n := 30

	btree := NewMVCCBtree(m)
	for i := 1; i <= n; i++ {
		key := []byte(fmt.Sprintf("%04d", i))
		value := []byte(fmt.Sprintf("Value%04d", i))
		btree.Put(key, value)
	}

	writeDotSvg(testDotExePath, "output.svg", btree, "")

	begin := []byte("0011")
	end := []byte("0022")
	iter := btree.NewIterator(begin, end)
	for iter.Next() {
		t.Logf("key: %v, value: %v", string(iter.Key()), string(iter.Value()))
	}
}
