package main

import (
	"testing"

	"github.com/TingYunAPM/go/utils/map"
)

func lessMap(a interface{}, b interface{}) bool {
	res := a.(int) < b.(int)
	return res
}

const testMapCount = 100

func TestMap(t *testing.T) {
	mymap := ordermap.New(lessMap)
	for i := 0; i < testMapCount/2; i++ {
		mymap.Update(i, testMapCount-i-1)
		mymap.Update(testMapCount-i-1, i)
	}
	end := mymap.End()
	for i := 0; i < testMapCount; i++ {
		if mymap.Find(i).Eq(end) {
			t.Error("%lld not found\n", i)
		}
	}
	id := 0
	for it := mymap.Begin(); !it.Eq(end); it.MoveNext() {
		if it.First().(int) != id {
			t.Error("%d : %d failed\n", id, it.First().(int))
		}
		id++
	}
	if id != testMapCount {
		t.Log("count %d", id)
	}
}
