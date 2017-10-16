// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"strings"
	"time"
)

type structAction struct {
	perf         structActionPerf
	commentTrace structActionTraceSet
	errorTrace   structErrorTraceSet
}

func (a *structAction) Destroy() {
	a.perf.Reset()
	a.commentTrace.Reset()
	a.errorTrace.Destroy()
}

func (a *structAction) init(apdex_t int32) *structAction {
	a.perf.Init(apdex_t)
	a.commentTrace.Init()
	a.errorTrace.Init()
	return a
}
func createStructAction(apdex_t int32) *structAction {
	return (&structAction{}).init(apdex_t)
}
func (a *structAction) Merge(t *structAction) {
	a.perf.Merge(&t.perf)
	a.commentTrace.Merge(&t.commentTrace)
	a.errorTrace.Merge(&t.errorTrace)
}

func getAppSecId(id string) string {

	if sidArray := strings.Split(id, "|"); len(sidArray) < 2 {
		return ""
	} else {
		trackId := sidArray[1]
		for i := 2; i < len(sidArray); i++ {
			trackId += "|" + sidArray[i]
		}
		return trackId
	}
}

//TINGYUN_ID_SECRET;c=CALLER_TYPE;r=REQ_ID;x=TX_ID;e=EXTERNAL_ID

func parseTrackId(id string) map[string]interface{} {
	if array := strings.Split(id, ";"); len(array) < 4 { // sid, type, tid, eid
		return nil
	} else if trackId := getAppSecId(array[0]); trackId == "" {
		return nil
	} else {
		if tmpArray := strings.Split(trackId, "#"); len(tmpArray) != 2 {
			return nil
		} else if len(tmpArray[0]) == 0 || len(tmpArray[1]) == 0 {
			return nil
		}
		tid, eid := "", ""
		for i := 1; i < len(array); i++ {
			parser := strings.Split(array[i], "=")
			if len(parser) != 2 {
				return nil
			}
			switch parser[0] {
			case "x":
				tid = parser[1]
			case "e":
				eid = parser[1]
			default:
			}
		}
		if tid == "" || eid == "" {
			return nil
		}
		r := make(map[string]interface{})
		r["applicationId"] = trackId
		r["transactionId"] = tid
		r["externalId"] = eid
		return r
	}

}
func mkTime(action time.Duration, dbTime time.Duration, exTime time.Duration, redisTime time.Duration, mcTime time.Duration, mongoTime time.Duration, aloneTime time.Duration) map[string]interface{} {
	r := make(map[string]interface{})
	r["duration"] = action / time.Millisecond
	r["qu"] = 0
	r["db"] = dbTime / time.Millisecond
	r["ex"] = exTime / time.Millisecond
	r["rds"] = redisTime / time.Millisecond
	r["mc"] = mcTime / time.Millisecond
	r["mon"] = mongoTime / time.Millisecond
	r["code"] = aloneTime / time.Millisecond
	return r
}
func (a *structAction) Add(action *Action, onComponent func(*Component)) {
	trace := newStructActionTrace(&action.time, action.url, action.statusCode)
	x := &trace
	trace.name = action.name
	trace.requestParams, trace.customParams = action.requestParams, action.customParams
	action_slow := action.Slow()
	a.perf.Add(action, func(component *Component) {
		trace := *x
		trace.AddSegment(component)
		if component._type != ComponentDefault {
			onComponent(component)
		}
		trace = nil
	})
	x = nil
	trace.formatComponent()
	aloneTime := time.Duration(0)
	dbTime := aloneTime
	mcTime := aloneTime
	mongoTime := aloneTime
	redisTime := aloneTime
	exTime := aloneTime
	trace.forEachComponent(func(component *Component) {
		switch component._type {
		case ComponentDefault:
			aloneTime += component.aloneTime
			timeVal := float64(component.time.duration / time.Millisecond)
			excl := float64(component.aloneTime / time.Millisecond)
			a.perf.componentPerf.ExclAdd(component.name, timeVal, excl)
			return
		case ComponentDefaultDB:
			dbTime += component.time.duration
			return
		case ComponentMysql:
			dbTime += component.time.duration
			return
		case ComponentPostgreSql:
			dbTime += component.time.duration
			return
		case ComponentMongo:
			mongoTime += component.time.duration
			return
		case ComponentRedis:
			redisTime += component.time.duration
			return
		case ComponentMemCache:
			mcTime += component.time.duration
			return
		case ComponentExternal:
			exTime += component.time.duration
			return
		default:
			return
		}
	})
	params_moved := false
	if action_slow && !trace.Empty() {
		trace.txId = action.unicId()
		if len(action.trackId) > 0 {
			entryTrace := parseTrackId(action.trackId)
			if entryTrace != nil {
				entryTrace["time"] = mkTime(action.time.duration, dbTime, exTime, redisTime, mcTime, mongoTime, aloneTime)
				trace.customParams["entryTrace"] = entryTrace
			}
		}
		a.commentTrace.Add(trace)
		params_moved = true
	} else {
		trace.destroy()
	}
	//Error Trace
	for pointer, _ := action.errors.Get(); pointer != nil; pointer, _ = action.errors.Get() {
		errinfo := pointer.(*errInfo)
		a.errorTrace.Add(errinfo, action)
		params_moved = true
		errinfo.Destroy()
	}
	if params_moved {
		action.requestParams, action.customParams = nil, nil
	}
}
