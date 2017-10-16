// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"encoding/json"
	"fmt"
	"time"
)

type structErrorTrace struct {
	errorTime     time.Time
	action        string
	className     string
	count         uint32
	httpStatus    int
	uri           string
	requestParams map[string]string
	customParams  map[string]interface{}
}
type structErrorTraceSet struct {
	errors map[string]map[string]*structErrorTrace
}

func (t *structErrorTrace) Destroy() {
	if t.requestParams != nil {
		for s, _ := range t.requestParams {
			delete(t.requestParams, s)
		}
		t.requestParams = nil
	}
	if t.customParams != nil {
		for s, _ := range t.customParams {
			delete(t.customParams, s)
		}
		t.customParams = nil
	}
}

func (e *structErrorTraceSet) Init() {
	e.errors = make(map[string]map[string]*structErrorTrace)
}
func (e *structErrorTraceSet) initStack(stack string) map[string]*structErrorTrace {
	if v, ok := e.errors[stack]; ok {
		return v
	} else {
		v = make(map[string]*structErrorTrace)
		e.errors[stack] = v
		return v
	}
}
func mergeMessage(local map[string]*structErrorTrace, peer map[string]*structErrorTrace) {
	for k, v := range peer {
		if addValue(local, k, v) {
			delete(peer, k)
		}
	}
}
func addValue(m map[string]*structErrorTrace, key string, val *structErrorTrace) bool {
	if v, ok := m[key]; ok {
		m[key].count = v.count + val.count
		return false
	}
	m[key] = val
	return true
}

//!used
func (e *structErrorTraceSet) Len() int {
	return len(e.errors)
}
func (e *structErrorTraceSet) Add(err_info *errInfo, action *Action) {
	stack, message := err_info.stack, fmt.Sprint(err_info.e)
	trace := &structErrorTrace{}
	trace.action = action.name
	trace.errorTime = err_info.happenTime
	trace.customParams = action.customParams
	trace.requestParams = action.requestParams
	trace.uri = action.url
	trace.className = err_info.eType
	trace.count = 1
	trace.httpStatus = int(action.statusCode)

	addValue(e.initStack(stack), message, trace)
}
func (e *structErrorTraceSet) Merge(t *structErrorTraceSet) {
	for k, v := range e.errors {
		mergeMessage(e.initStack(k), v)
	}
}

//			ERROR_TIME_IN_SECONDS,
//			“ACTION_METRIC_NAME”,
//			HTTP_STATUS,
//			“EXCEPTION_CLASS_NAME”,
//			“ERROR_MESSAGE”,
//			COUNT_OF_ERRORS,
//			“REQUEST_URI”,
//			“ERROR_PARAMS”
func makeErrorTraceJson(stack string, message string, trace *structErrorTrace) interface{} {
	ret := make([]interface{}, 8)
	ret[0] = trace.errorTime.Unix()
	ret[1] = trace.action
	ret[2] = trace.httpStatus
	ret[3] = trace.className
	ret[4] = message
	ret[5] = trace.count
	ret[6] = trace.uri
	params := make(map[string]interface{})
	params["params"] = trace.customParams
	params["requestParams"] = trace.requestParams
	params["stacktrace"] = jsonDecodeArray(stack)
	paramsJson, _ := json.Marshal(params)
	ret[7] = string(paramsJson)
	return ret
}
func (e *structErrorTraceSet) Read(put func(interface{})) {
	for k, v := range e.errors {
		for msg, trace := range v {
			put(makeErrorTraceJson(k, msg, trace))
		}
	}
}

//!used
func (e *structErrorTraceSet) Destroy() {
	if e.errors == nil {
		return
	}
	for k, v := range e.errors {
		for s, t := range v {
			t.Destroy()
			delete(v, s)
		}
		delete(e.errors, k)
	}
	e.errors = nil
}
