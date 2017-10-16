// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/TingYunAPM/go/utils/list"
)

type structAppData struct {
	actions        map[string]*structAction
	sqlTraces      mapStructSqlTrace
	generalMetrics mapPerformance
	startTime      time.Time
	endTime        time.Time
	sys            *sysInfo
	runtimeData    runtimeBlock
}

func (r *structAppData) init(startTime time.Time) *structAppData {
	r.actions = make(map[string]*structAction)
	r.sqlTraces = newMapStructSqlTrace()
	r.startTime = startTime
	r.generalMetrics = newMapPerformance()
	return r
}

func (r *structAppData) end(perf *runtimePerf) {
	perf.Snap()
	r.runtimeData.Read(perf)
	perf.Reset()
}
func (r *structAppData) getStructAction(name string) *structAction {
	ret, ok := r.actions[name]
	if !ok {
		ret = createStructAction(app.ReadApdex(name))
		r.actions[name] = ret
	}
	return ret
}

//追加action数据
func (r *structAppData) Append(action *Action) {
	if r.actions == nil || action.stateUsed != actionFinished {
		return
	}
	//慢sql
	enableActionTrace := readServerConfigBool(configServerConfigBoolActionTracerEnabled, true)
	enableSlowSql := readServerConfigBool(configServerConfigBoolActionTracerSlowSQL, true)
	enabled := enableActionTrace && enableSlowSql
	slowValue := time.Duration(readServerConfigInt(configServerConfigIntegerActionTracerSlowSQLThreshold, 500)) * time.Millisecond
	onComponent := func(component *Component) {
		timeVal := float64(component.time.duration / time.Millisecond)
		r.generalMetrics.Add(component.metricName(), timeVal)
		if enabled && component.time.duration >= slowValue && component.isDatabaseComponent() {
			r.sqlTraces.Add(action, component)
		}
	}
	r.getStructAction(action.name).Add(action, onComponent)
	if action.trackId != "" {
		if metric := getTopMetric(action.trackId); metric != "" {
			r.generalMetrics.Add(metric, float64(action.time.duration/time.Millisecond))
		}
	}
}

func (r *structAppData) ReadActionMetrics() []interface{} {
	count := len(r.actions)
	ret := make([]interface{}, count)
	i := 0
	for k, v := range r.actions {
		name := map[string]string{"name": k}
		ret[i] = []interface{}{name, v.perf.action_time.IntSlice()}
		i++
	}
	return ret
}
func actionToApdex(metric string) string {
	return strings.Replace(metric, "WebAction", "Apdex", 1)
}
func (r *structAppData) ReadApdex() []interface{} {
	count := len(r.actions)
	ret := make([]interface{}, count)
	i := 0
	for k, v := range r.actions {
		key := map[string]string{"name": actionToApdex(k)}
		ret[i] = []interface{}{key, v.perf.apdex.IntSlice()}
		i++
	}
	return ret

}
func (r *structAppData) ReadComponents() []interface{} {
	var rlist list.List
	rlist.Init()
	for k, v := range r.actions {
		onElement := func(name string, perf *structPerformance) {
			key := make(map[string]string)
			key["name"] = name
			key["parent"] = k
			element := [...]interface{}{key, perf.IntSlice()}
			rlist.PushBack(element)
		}
		v.perf.componentPerf.Read(onElement)
		v.perf.sql_perfs.Read(onElement)
		v.perf.externalPerf.Read(onElement)
		v.perf.nosqlPerf.Read(onElement)
	}
	ret := make([]interface{}, rlist.Size())
	for i := 0; rlist.Size() > 0; i++ {
		element, _ := rlist.PopFront()
		ret[i] = element
	}
	return ret
}

const (
	//采样值,采样时刻进程内的go程数
	metricNumGoroutine = "GoRuntime/NULL/Goroutine"
	//单位时间内事件次数, =>一个累加值 在两次采样之间的差值。
	metricNumCgoCall = "GoRuntime/NULL/CgoCall"
	//单位时间内GC耗时的累加和,单位毫秒
	metricPauseTotalMs = "GoRuntime/NULL/PauseTotalMs"
	//单位时间内,每次GC耗时的5值统计性能数据
	metricGCTime = "GC/NULL/Time"
	//单位时间内 Free的次数
	metricFrees = "GoRuntime/NULL/Frees"
	//单位时间内 Malloc的次数
	metricMallocs = "GoRuntime/NULL/Mallocs"
	//单位时间内 Lookup的次数
	metricLookups = "GoRuntime/NULL/Lookups"
	//采样值,系统总的申请内存数 MB
	metricMemTotalSys = "Memory/NULL/MemSys"
	//采样值,系统栈内存数 MB
	metricMemStackSys = "Memory/Stack/StackSys"
	//采样值,系统堆内存数 MB
	metricMemHeapSys = "Memory/Heap/HeapSys"
	//采样值,系统内存区间结构数
	metricMSpanSys = "Memory/MSpan/MSpanSys"
	//采样值,系统内存Cache结构数
	metricMCacheSys = "Memory/MCache/MCacheSys"
	//采样值,系统内存BuckHash数
	metricBuckHashSys = "Memory/NULL/BuckHashSys"
	//采样值,使用中的堆内存数 MB
	metricHeapInuse = "Memory/Heap/HeapInuse"
	//采样值,使用中的栈内存数 MB
	metricStackInuse = "Memory/Stack/StackInuse"
	//采样值,使用中的内存区间结构数
	metricMSpanInuse = "Memory/MSpan/MSpanInuse"
	//采样值,使用中的内存Cache结构数
	metricMCacheInuse     = "Memory/MCache/MCacheInuse"
	metricUserTime        = "CPU/NULL/UserTime"
	metricUserUtilization = "CPU/NULL/UserUtilization"
	metricmem             = "Memory/NULL/PhysicalUsed"
	//采样值,进程打开文件句柄数(linux)
	metricFDSize = "Process/NULL/FD"
	//采样值,进程内的系统线程数(linux)
	metricThreads = "Process/NULL/Threads"
)

func (r *structAppData) ReadGeneral() []interface{} {
	if !agentEnabled() {
		return []interface{}{}
	}
	sumMetrics := make(map[string]*structPerformance)
	add := func(metric string, perf *structPerformance) {
		v, ok := sumMetrics[metric]
		if !ok {
			v = newStructPerformance()
			sumMetrics[metric] = v
		}
		v.Append(perf)
	}
	sumAdd := func(metric string, perf *structPerformance) bool {
		array := strings.Split(metric, "/")
		if array[0] == "EntryTransaction" {
			add(array[0]+"/NULL/"+array[2], perf)
			return array[1] != "NULL"
		}
		if strings.Index(array[0], "Database") != -1 {
			add(array[0]+"/NULL/"+array[2], perf)
		}
		add(array[0]+"/NULL/All", perf)
		add(array[0]+"/NULL/AllWeb", perf)
		return array[1] != "NULL"
	}
	tempMetrics := make(map[int]interface{})
	i := 0
	for k, v := range r.generalMetrics {
		if sumAdd(k, v) { //过滤掉 xxxx/NULL/yyy 数据
			tempMetrics[i] = []interface{}{map[string]string{"name": k}, v.IntSlice()}
			i++
		}
	}
	for k, v := range sumMetrics {
		tempMetrics[i] = []interface{}{map[string]string{"name": k}, v.IntSlice()}
		i++
	}
	count := len(tempMetrics)
	//	if r.runtimeData.UserTime.access_count > 0 {
	count += 21
	//	}
	ret := make([]interface{}, count)
	for k, v := range tempMetrics {
		ret[k] = v
		delete(tempMetrics, k)
	}
	if i < count {
		ret[i+0] = []interface{}{map[string]string{"name": metricUserTime}, r.runtimeData.UserTime.FloatSlice()}
		ret[i+1] = []interface{}{map[string]string{"name": metricUserUtilization}, r.runtimeData.UserUtilization.FloatSlice()}
		ret[i+2] = []interface{}{map[string]string{"name": metricmem}, r.runtimeData.mem.FloatSlice()}
		ret[i+3] = []interface{}{map[string]string{"name": metricNumGoroutine}, r.runtimeData.NumGoroutine.FloatSlice()}
		ret[i+4] = []interface{}{map[string]string{"name": metricNumCgoCall}, r.runtimeData.NumCgoCall.FloatSlice()}

		//		ret[i+5] = []interface{}{map[string]string{"name": metricPauseTotalMs}, r.runtimeData.PauseTotalNs.FloatSlice()}
		//		ret[i+6] = []interface{}{map[string]string{"name": metricNumGC}, r.runtimeData.NumGC.FloatSlice()}
		ret[i+5] = []interface{}{map[string]string{"name": metricGCTime}, r.runtimeData.GCTime.FloatSlice()}
		ret[i+6] = []interface{}{map[string]string{"name": metricFrees}, r.runtimeData.Frees.FloatSlice()}
		ret[i+7] = []interface{}{map[string]string{"name": metricMallocs}, r.runtimeData.Mallocs.FloatSlice()}
		ret[i+8] = []interface{}{map[string]string{"name": metricLookups}, r.runtimeData.Lookups.FloatSlice()}
		ret[i+9] = []interface{}{map[string]string{"name": metricMemTotalSys}, r.runtimeData.MemTotalSys.FloatSlice()}
		ret[i+10] = []interface{}{map[string]string{"name": metricMemStackSys}, r.runtimeData.MemStackSys.FloatSlice()}
		ret[i+11] = []interface{}{map[string]string{"name": metricMemHeapSys}, r.runtimeData.MemHeapSys.FloatSlice()}
		ret[i+12] = []interface{}{map[string]string{"name": metricMSpanSys}, r.runtimeData.MSpanSys.FloatSlice()}
		ret[i+13] = []interface{}{map[string]string{"name": metricMCacheSys}, r.runtimeData.MCacheSys.FloatSlice()}
		ret[i+14] = []interface{}{map[string]string{"name": metricBuckHashSys}, r.runtimeData.BuckHashSys.FloatSlice()}

		ret[i+15] = []interface{}{map[string]string{"name": metricHeapInuse}, r.runtimeData.HeapInuse.FloatSlice()}
		ret[i+16] = []interface{}{map[string]string{"name": metricStackInuse}, r.runtimeData.StackInuse.FloatSlice()}
		ret[i+17] = []interface{}{map[string]string{"name": metricMSpanInuse}, r.runtimeData.MSpanInuse.FloatSlice()}
		ret[i+18] = []interface{}{map[string]string{"name": metricMCacheInuse}, r.runtimeData.MCacheInuse.FloatSlice()}
		ret[i+19] = []interface{}{map[string]string{"name": metricFDSize}, r.runtimeData.FDSize.FloatSlice()}
		ret[i+20] = []interface{}{map[string]string{"name": metricThreads}, r.runtimeData.Threads.FloatSlice()}
	}
	return ret
}

func (r *structAppData) ReadErrors() []interface{} {

	count := 3
	error_count_all := uint32(0)
	for _, v := range r.actions {
		error_count_all += v.perf.errorCount
		if v.perf.errorCount > 0 {
			count++
		}
	}
	if count == 3 {
		return []interface{}{}
	}
	ret := make([]interface{}, count)
	i := 0
	for k, v := range r.actions {
		if v.perf.errorCount > 0 {
			name := map[string]string{"name": "Errors/Count/" + url.QueryEscape(k)}
			ret[i] = []interface{}{name, []interface{}{v.perf.errorCount}}
			i++
		}
	}
	error_names := [3]string{"Errors/Count/All", "Errors/Count/AllWeb", "Errors/Count/AllBackground"}
	error_values := [3]uint32{error_count_all, error_count_all, 0}
	for k := 0; k < 3; k++ {
		ret[i+k] = []interface{}{map[string]string{"name": error_names[k]}, []interface{}{error_values[k]}}
	}

	return ret
}

//数据序列化
func (r *structAppData) Serialize() ([]byte, error) {
	perfMetrics := make(map[string]interface{})
	perfMetrics["type"] = "perfMetrics"
	perfMetrics["timeFrom"] = r.startTime.Unix()
	perfMetrics["timeTo"] = r.endTime.Unix()
	perfMetrics["interval"] = app.configs.server.CIntegers.Read(configServerIntegerDataSentInterval, 60)
	perfMetrics["actions"] = r.ReadActionMetrics()
	perfMetrics["apdex"] = r.ReadApdex()
	perfMetrics["components"] = r.ReadComponents()
	perfMetrics["general"] = r.ReadGeneral()
	perfMetrics["errors"] = r.ReadErrors()
	data := [...]interface{}{perfMetrics, r.ReadActionTraces(), r.ReadErrorTraces(), r.ReadSqlTraces()}
	defer func() {
		delete(perfMetrics, "type")
		delete(perfMetrics, "timeFrom")
		delete(perfMetrics, "timeTo")
		delete(perfMetrics, "interval")
		delete(perfMetrics, "actions")
		delete(perfMetrics, "apdex")
		delete(perfMetrics, "components")
		delete(perfMetrics, "general")
		delete(perfMetrics, "errors")
		for i := 0; i < 4; i++ {
			data[i] = nil
		}
	}()
	return json.Marshal(data)
}

func (r *structAppData) ReadActionTraces() interface{} {
	trace_count := 0
	for _, v := range r.actions {
		trace_count += v.commentTrace.Len()
	}
	action_traces := make([]interface{}, trace_count)
	index := 0
	append_trace := func(trace interface{}) {
		action_traces[index] = trace
		index++
	}
	for _, v := range r.actions {
		v.commentTrace.Read(append_trace)
	}
	traces := make(map[string]interface{})
	traces["type"] = "actionTraceData"
	traces["actionTraces"] = action_traces
	return traces
}
func (r *structAppData) ReadSqlTraces() interface{} {
	traces := make(map[string]interface{})
	traces["type"] = "sqlTraceData"
	traces["sqlTraces"] = r.sqlTraces.Read()
	return traces
}
func (r *structAppData) ReadErrorTraces() interface{} {
	var trace_cache list.List
	trace_cache.Init()
	for _, v := range r.actions {
		v.errorTrace.Read(func(trace interface{}) {
			//fmt.Println("Got trace:", trace)
			trace_cache.PushBack(trace)
		})
	}

	traceArray := make([]interface{}, trace_cache.Size())
	for index := 0; trace_cache.Size() > 0; index++ {
		trace, _ := trace_cache.PopFront()
		traceArray[index] = trace
	}
	traces := make(map[string]interface{})
	traces["type"] = "errorTraceData"
	traces["errors"] = traceArray

	return traces
}

//合并两块数据
func (r *structAppData) Merge(t *structAppData) {
	for k, v := range t.actions {
		if val, found := r.actions[k]; found {
			val.Merge(v)
		} else {
			r.actions[k] = v
			delete(t.actions, k)
		}
	}
	for k, v := range t.generalMetrics {
		if val, found := r.generalMetrics[k]; found {
			val.Append(v)
		} else {
			r.generalMetrics[k] = v
			delete(t.generalMetrics, k)
		}
	}
	r.sqlTraces.Merge(t.sqlTraces)
	if t.startTime.Before(r.startTime) {
		r.startTime = t.startTime
	}
	if t.endTime.After(r.endTime) {
		r.endTime = t.endTime
	}
	r.runtimeData.Merge(&t.runtimeData)
}

//释放内存
func (r *structAppData) destroy() {
	r.sqlTraces = nil
	if r.actions != nil {
		for k, v := range r.actions {
			v.Destroy()
			delete(r.actions, k)
		}
		r.actions = nil
	}
	r.sys = nil
	if r.generalMetrics != nil {
		r.generalMetrics.Reset()
		r.generalMetrics = nil
	}
}
func (a *application) CreateReportBlock(startTime time.Time) {
	if a.dataBlock == nil {
		a.dataBlock = new(structAppData).init(startTime)
	}
}
