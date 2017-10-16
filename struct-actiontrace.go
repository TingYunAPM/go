// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"encoding/json"
	"time"

	"github.com/TingYunAPM/go/utils/list"
)

type structActionTrace struct {
	name          string
	uri           string
	txId          string
	statusCode    uint16
	requestParams map[string]string
	customParams  map[string]interface{}
	time          timeRange
	components    list.List
}

func clearComponent(nlist *list.List) {
	for nlist.Size() > 0 {
		v, _ := nlist.PopFront()
		v.(*Component).destroy()
	}
}
func appendComponent(nlist *list.List, new_component *Component) {
	//  front                                               back
	//{time.begin:1}<-{time.begin:2}<-{time.begin:3}<-{time.begin:4}
	//
	for iter := nlist.Back(); !iter.IsEnd(); iter.MoveFront() {
		curr, _ := iter.Value()
		//新增项不早于当前项
		if !new_component.time.begin.Before(curr.(*Component).time.begin) {
			iter.InsertBack(new_component)
			iter.Destroy()
			return
		}
	}
	nlist.PushFront(new_component)
}

func (c *Component) appendSub(sub *Component) {
	appendComponent(&c.subs, sub)
}
func (a *structActionTrace) AddSegment(component *Component) {
	if component.pre == nil {
		appendComponent(&a.components, component)
	} else {
		component.pre.appendSub(component)
	}
}
func (a *structActionTrace) Empty() bool { return a.components.Size() == 0 }
func newStructActionTrace(time *timeRange, uri string, statusCode uint16) *structActionTrace {
	r := &structActionTrace{name: "", uri: uri, txId: "", statusCode: statusCode, requestParams: nil, customParams: nil, time: *time}
	r.components.Init()
	return r
}
func serialJson(components *list.List, begin time.Time) []interface{} {
	if components.Size() == 0 {
		return []interface{}{}
	}
	ret := make([]interface{}, components.Size())
	for i, it := 0, components.Front(); !it.IsEnd(); it.MoveBack() {
		v, _ := it.Value()
		ret[i] = v.(*Component).serialJson(begin)
		i++
	}
	return ret
}
func (c *Component) serialJson(begin time.Time) interface{} {
	segment := make([]interface{}, 9)
	start := c.time.begin.Sub(begin) / time.Millisecond
	segment[0] = start
	segment[1] = start + c.time.duration/time.Millisecond
	segment[2] = c.metricName()
	segment[3] = c.getURL()
	segment[4] = 1
	segment[5], segment[6] = parseMethod(c.method)
	params := make(map[string]interface{})
	if c.callStack != nil {
		sql := c.getSQL()
		if sql != "" {
			params["sql"] = sql
		}
		params["stacktrace"] = c.callStack
	}
	if c.exId {
		params["externalId"] = c.unicId()
	}
	segment[7] = params
	segment[8] = serialJson(&c.subs, begin)
	return segment
}
func (p *structActionTrace) Read() interface{} {
	var traceData [4]interface{}
	traceData[0] = p.time.begin.Unix()
	traceData[1] = p.requestParams
	traceData[2] = p.customParams
	var rootSegment [9]interface{}
	rootSegment[0] = 0
	rootSegment[1] = p.time.duration / time.Millisecond
	rootSegment[2] = ""
	rootSegment[3] = p.uri
	rootSegment[4] = 1
	rootSegment[5] = "Go"
	rootSegment[6] = "execute"
	rootSegment[7] = make(map[string]interface{})
	rootSegment[8] = serialJson(&p.components, p.time.begin)
	traceData[3] = &rootSegment
	var traceItem [7]interface{}
	traceItem[0] = p.time.begin.Unix()
	traceItem[1] = p.time.duration / time.Millisecond
	traceItem[2] = p.name
	traceItem[3] = p.uri
	traceByte, _ := json.Marshal(&traceData)
	//traceData释放内存
	for i := 0; i < 4; i++ {
		traceData[i] = nil
	}
	for i := 0; i < 9; i++ {
		rootSegment[i] = nil
	}
	//traceData End
	traceString := string(traceByte)
	traceItem[4] = traceString
	traceItem[5] = p.txId
	traceItem[6] = md5sum(p.name + md5sum(traceString) + p.time.begin.String() + p.time.duration.String())

	return &traceItem
}
func formatList(components *list.List) time.Duration {
	ret := time.Duration(0)
	for it := components.Front(); !it.IsEnd(); it.MoveBack() {
		v, _ := it.Value()
		component := v.(*Component)
		ret += component.time.duration
		if component._type == ComponentDefault {
			component.aloneTime = component.time.duration - formatList(&component.subs)
			if component.aloneTime < 0 {
				component.aloneTime = 0
			}
		}
	}
	return ret
}
func segmentForEach(components *list.List, cb func(*Component)) {
	if components != nil {
		for it := components.Front(); !it.IsEnd(); it.MoveBack() {
			v, _ := it.Value()
			component := v.(*Component)
			cb(component)
			segmentForEach(&component.subs, cb)
		}
	}
}
func (p *structActionTrace) forEachComponent(cb func(*Component)) {
	segmentForEach(&p.components, cb)
}
func (p *structActionTrace) formatComponent() {
	formatList(&p.components)
}
func (p *structActionTrace) destroy() {
	clearComponent(&p.components)
	p.txId = ""
	p.name = ""
	p.uri = ""
	p.requestParams = nil
	p.customParams = nil
}
