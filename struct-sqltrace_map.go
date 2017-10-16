// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"encoding/json"
	"time"
)

type mapStructSqlTrace map[string]*structSqlTrace

func newMapStructSqlTrace() mapStructSqlTrace {
	return mapStructSqlTrace{}
}
func (m mapStructSqlTrace) Add(action *Action, component *Component) {
	keyName := component.metricName() + "|" + action.name
	val, ok := m[keyName]
	if !ok {
		val = newStructSqlTrace()
		m[keyName] = val
	}
	perf := float64(component.time.duration / time.Millisecond)
	if val.perf.Count == 0 || perf > val.perf.Max {
		jsonbyte, _ := json.Marshal(map[string]interface{}{"stacktrace": component.callStack})
		val.stackTrace = string(jsonbyte)
		val.uri = action.url
		val.time = component.time
	}
	val.perf.Add(perf)
}
func (m mapStructSqlTrace) Merge(t mapStructSqlTrace) {
	for k, v := range t {
		val, ok := m[k]
		if !ok {
			val = newStructSqlTrace()
			m[k] = val
		}
		if val.perf.Count == 0 || v.perf.Max > val.perf.Max {
			val.stackTrace = v.stackTrace
			val.time = v.time
			val.uri = v.uri
		}
		val.perf.Append(&v.perf)
	}
}
func (m mapStructSqlTrace) Read() interface{} {
	traceArray := make([]interface{}, len(m))
	i := 0
	for k, v := range m {
		traceArray[i] = v.Read(k)
		i++
	}
	return traceArray
}
