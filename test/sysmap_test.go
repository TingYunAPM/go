package main

import (
	"testing"
)

const testSysMapCount = 1000000

func TestSysMap(t *testing.T) {
	mymap := make(map[int]int)
	for i := 0; i < testSysMapCount/2; i++ {
		mymap[i] = testSysMapCount - i - 1
		mymap[testSysMapCount-i-1] = i
	}

	for i := 0; i < testSysMapCount; i++ {
		if _, ok := mymap[i]; !ok {
			t.Error("%d not found\n", i)
		}
	}

	id := 0
	for key, value := range mymap {
		if key+value != testSysMapCount-1 {
			t.Error("%d : %d %d failed\n", id, key, value)
		}
		id++
	}
	if id != testSysMapCount {
		t.Error("count %d", id)
	}

}
