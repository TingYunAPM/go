package main

import (
	"testing"

	"github.com/TingYunAPM/go/utils/rbtree"
)

func lessRBT(a interface{}, b interface{}) bool {
	res := a.(int) < b.(int)
	return res
}

const testRBTCount = 100

func TestRBT(t *testing.T) {

	tree := rbtree.NewTree(lessRBT)
	for i := 0; i < testRBTCount/2; i++ {
		if !tree.Insert(rbtree.NewNode(i)) {
			t.Error("insert %lld error\n", i)
		}
		if !tree.Insert(rbtree.NewNode(testRBTCount - i - 1)) {
			t.Error("insert %lld error\n", i)
		}
	}
	for i := 0; i < testRBTCount; i++ {
		if tree.Find(i) == nil {
			t.Error("%lld not found\n", i)
		}
	}
	id := 0
	for n := tree.Left(); n != nil; n, id = n.Right(), id+1 {
		if n.Value != id {
			t.Error("%lld : %lld failed\n", id, n.Value)
		}
	}
	if id != testRBTCount {
		t.Error("count %lld", id)
	}
}
