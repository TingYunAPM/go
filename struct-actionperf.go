// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"time"
)

//structActionPerf Action性能数据
type structActionPerf struct {
	apdex         structApdex
	action_time   structPerformance
	componentPerf mapPerformance
	externalPerf  mapPerformance
	nosqlPerf     mapPerformance
	sql_perfs     mapSqlPerf
	errorCount    uint32
	accessCount   uint32
}

func (p *structActionPerf) Init(apdex_t int32) {
	p.Reset()
	p.apdex.apdex_t_in_millis = apdex_t
}

func (p *structActionPerf) Reset() {
	p.apdex.Reset()
	p.action_time.Reset()
	if p.componentPerf == nil {
		p.componentPerf = newMapPerformance()
	}
	p.componentPerf.Reset()
	if p.externalPerf == nil {
		p.externalPerf = newMapPerformance()
	}
	p.externalPerf.Reset()
	if p.nosqlPerf == nil {
		p.nosqlPerf = newMapPerformance()
	}
	p.nosqlPerf.Reset()
	if p.sql_perfs == nil {
		p.sql_perfs = newMapSqlPerf()
	}
	p.sql_perfs.Reset()
	p.errorCount = 0
	p.accessCount = 0
}
func (p *structActionPerf) Merge(q *structActionPerf) {
	p.apdex.Merge(&q.apdex)
	p.action_time.Append(&q.action_time)
	p.componentPerf.Merge(q.componentPerf)
	p.sql_perfs.Merge(q.sql_perfs)
	p.externalPerf.Merge(q.externalPerf)
	p.nosqlPerf.Merge(q.nosqlPerf)
	p.errorCount += q.errorCount
	p.accessCount += q.accessCount
}
func (p *structActionPerf) Add(action *Action, cb func(*Component)) {
	p.accessCount++
	p.errorCount += uint32(action.errors.Size())
	p.action_time.AddValue(float64(action.time.duration/time.Millisecond), 0)
	p.apdex.Add(int32(action.time.duration/time.Millisecond), action.errors.Size() != 0)
	if action.errors.Size() != 0 {
		p.errorCount += 1
	}
	for action.cache.Size() > 0 {
		c, found := action.cache.Get()
		if !found {
			continue
		}
		component := c.(*Component)
		//fmt.Printf("addComponent(%d,%s,%s)\n", component._type, component.name, component.method)
		component.fixEnd(&action.time)
		timeVal := float64(component.time.duration / time.Millisecond)
		switch component._type {
		case ComponentExternal:
			p.externalPerf.ExclAdd(component.metricName(), timeVal, timeVal)
		case ComponentDefaultDB:
			p.sql_perfs.Add(component.name, timeVal, timeVal)
		case ComponentMemCache:
			p.nosqlPerf.ExclAdd(component.metricName(), timeVal, timeVal)
		case ComponentMongo:
			p.nosqlPerf.ExclAdd(component.metricName(), timeVal, timeVal)
		case ComponentMysql:
			p.sql_perfs.Add(component.name, timeVal, timeVal)
		case ComponentPostgreSql:
			p.sql_perfs.Add(component.name, timeVal, timeVal)
		case ComponentRedis:
			p.nosqlPerf.ExclAdd(component.metricName(), timeVal, timeVal)
		case ComponentDefault:
		default:
		}
		cb(component)
	}
}
