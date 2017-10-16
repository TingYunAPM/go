// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package pool

import (
	"sync/atomic"
)

//无锁消息池，多写,有序单读(按写入时间排序)。适用于多个消息发送者,单个消息接收者模式
type SerialReadPool struct {
	//由于atomic.AddUint64 在非8字节对齐的地址上会崩溃，所以这两个id必须放在开始位置
	writeId uint64
	readId  uint64
	p       Pool
	cache   map[uint64]interface{}
}
type msg struct {
	id uint64
	o  interface{}
}

func SerialNew() *SerialReadPool {
	return (&SerialReadPool{}).Init()
}
func (p *SerialReadPool) Init() *SerialReadPool {
	p.writeId = 0
	p.readId = 1
	p.cache = make(map[uint64]interface{})
	p.p.Init()
	return p
}

func (p *SerialReadPool) Size() int32 {
	return p.p.Size() + int32(len(p.cache))
}
func (p *SerialReadPool) cacheGet() interface{} {
	if r, exist := p.cache[p.readId]; exist {
		delete(p.cache, p.readId)
		p.readId++
		return r
	}
	return nil
}
func (p *SerialReadPool) Get() interface{} {
	if r := p.cacheGet(); r != nil { //cache里有可用的数据,从cache里取
		return r
	}
	for {
		u, poped := p.p.Get()
		if !poped { //底层的pool里没有message
			return nil
		}
		m := u.(*msg)
		r := m.o
		m.o = nil
		id := m.id
		if id == p.readId { //是最早入队的message
			p.readId++
			return r
		}
		//不是最早的那个message，扔到cache里
		p.cache[id] = r
	}
}
func (p *SerialReadPool) Put(o interface{}) {
	id := atomic.AddUint64(&(p.writeId), 1)
	p.p.Put(&msg{id, o})
}
