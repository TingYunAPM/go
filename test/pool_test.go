package main

import (
	"testing"

	"github.com/TingYunAPM/go/utils/pool"
)

func TestPool(t *testing.T) {
	msgCache := pool.New()
	worker := func() {
		for i := 0; i < 10000; i++ {
			msgCache.Put(i)
		}

	}
	go worker()
	go worker()
	count := 0
	for count < 20000 {
		_, ok := msgCache.Get()
		if ok {
			count++
		}
	}

	serialPool := pool.SerialNew()
	worker1 := func() {
		for i := 0; i < 10000; i++ {
			serialPool.Put(i)
		}
	}
	go worker1()
	go worker1()
	count = 0
	for count < 20000 {
		t := serialPool.Get()
		if t != nil {
			count++
		}
	}

}
