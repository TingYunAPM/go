// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"github.com/TingYunAPM/go/utils/list"
)

type structActionTraceSet struct {
	traces   list.List
	limitMax uint16
}

func (p *structActionTraceSet) Init() {
	p.limitMax = 20
	p.traces.Init()
}

//!used
func (p *structActionTraceSet) Reset() {
	for p.traces.Size() > 0 {
		trace, _ := p.traces.PopFront()
		trace.(*structActionTrace).destroy()
	}
}
func (p *structActionTraceSet) Merge(q *structActionTraceSet) {
	for q.traces.Size() > 0 {
		trace, _ := q.traces.PopFront()
		p.Add(trace.(*structActionTrace))
	}
}
func (p *structActionTraceSet) Read(put func(interface{})) {
	for it := p.traces.Front(); !it.IsEnd(); it.MoveBack() {
		v, _ := it.Value()
		put(v.(*structActionTrace).Read())
	}
}
func (p *structActionTraceSet) Len() int { return int(p.traces.Size()) }
func (p *structActionTraceSet) Add(t *structActionTrace) {
	if p.traces.Size() == 0 {
		p.traces.PushBack(t)
		return
	}
	iter := p.traces.Back()
	defer iter.Destroy()
	v, _ := iter.Value()
	if p.traces.Size() == int(p.limitMax) && t.time.duration <= v.(*structActionTrace).time.duration {
		t.destroy()
		return
	}
	for !iter.IsEnd() {
		v, _ = iter.Value()
		if t.time.duration <= v.(*structActionTrace).time.duration {
			break
		}
		iter.MoveFront()
	}
	if iter.IsEnd() {
		p.traces.PushFront(t)
	} else {
		iter.InsertBack(t)
	}
	if p.traces.Size() > int(p.limitMax) {
		trace, _ := p.traces.PopBack()
		trace.(*structActionTrace).destroy()
	}
}
