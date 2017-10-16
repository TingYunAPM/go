// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

//无锁消息池，多读多写, 用于goroutine 间收发消息
package pool

//接口: Put, Get, Size

import "sync/atomic"

type node struct {
	next  *node
	value interface{}
}
type lockList struct {
	lock int32
	head *node
	end  *node
}

func (l *lockList) init() {
	l.lock = 0
	l.head = nil
	l.end = nil
}

func (l *lockList) PushBack(n *node) bool {
	used := atomic.AddInt32(&l.lock, 1)
	defer atomic.AddInt32(&l.lock, -1)
	if used != 1 {
		return false
	}
	n.next = nil
	if l.end == nil {
		l.head = n
		l.end = n
	} else {
		l.end.next = n
		l.end = n
	}
	return true
}
func (l *lockList) PopFront() *node {
	if l.head == nil {
		return nil
	}
	used := atomic.AddInt32(&l.lock, 1)
	defer atomic.AddInt32(&l.lock, -1)
	if used != 1 {
		return nil
	}
	if l.head == nil {
		return nil
	}
	ret := l.head
	l.head = ret.next
	ret.next = nil
	if l.head == nil {
		l.end = nil
	}
	return ret
}

const bucketCount = 8

type nodePool struct {
	count      int32
	indexRead  int32
	indexWrite int32
	array      [bucketCount]lockList
}

func (p *nodePool) init() *nodePool {
	p.count = 0
	p.indexRead = 0
	p.indexWrite = 0
	for i := 0; i < bucketCount; i++ {
		p.array[i].init()
	}
	return p
}

func (p *nodePool) Put(n *node) {
	pwrite := p.indexWrite
	for {
		for i := pwrite; i-pwrite < bucketCount; i++ {
			listId := i % bucketCount
			if p.array[listId].PushBack(n) {
				atomic.AddInt32(&p.count, 1)
				p.indexWrite = listId + 1
				return
			}
		}
	}
}

func (p *nodePool) Size() int32 {
	return p.count
}

func (p *nodePool) Get() *node {
	if p.count == 0 {
		return nil
	}
	pread := p.indexRead
	for i := pread; i-pread < bucketCount; i++ {
		readListId := i % bucketCount
		r := p.array[readListId].PopFront()
		if r != nil {
			atomic.AddInt32(&p.count, -1)
			p.indexRead = readListId + 1
			return r
		}
	}
	return nil
}

type Pool struct {
	pool nodePool
}

func (p *Pool) Init() *Pool {
	p.pool.init()
	return p
}
func New() *Pool {
	return new(Pool).Init()
}

func (p *Pool) Put(v interface{}) {
	p.pool.Put(&node{next: nil, value: v})
}
func (p *Pool) Size() int32 {
	return p.pool.Size()
}
func (p *Pool) Get() (interface{}, bool) {
	n := p.pool.Get()
	if n == nil {
		return nil, false
	}
	ret := n.value
	n.value = nil
	return ret, true
}
