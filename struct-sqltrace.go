// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"strings"
)

type structSqlTrace struct {
	perf       structSqlTracePerf
	time       timeRange
	stackTrace string
	uri        string
}

func newStructSqlTrace() *structSqlTrace {
	ret := &structSqlTrace{}
	ret.perf.Reset()
	ret.time.Init()
	ret.stackTrace = "{}"
	ret.uri = ""
	return ret
}

//			SQL_TRACE_TIME_IN_SECONDS,
//			“ACTION_ ROOT_METRIC_NAME”,
//			“METRIC_NAME”,
//			“REQUEST_URI”,
//			“SQL_OBFUSCATED”,
//			CALL_COUNT,
//			SUM_PERF,
//			MAX_PERF,
//			MIN_PERF,
//			“SQL_PARAMS”
func parseSqlTraceKey(key string) (string, string) {
	array := strings.Split(key, "|")
	cMetric := array[0]
	aMetric := array[1]
	for i := 2; i < len(array); i++ {
		aMetric = aMetric + "|" + array[i]
	}
	return aMetric, cMetric
}
func (s *structSqlTrace) Read(key string) interface{} {
	traceItem := make([]interface{}, 10)
	traceItem[0] = s.time.begin.Unix()
	traceItem[1], traceItem[2] = parseSqlTraceKey(key)
	traceItem[3] = s.uri
	traceItem[4] = traceItem[2]
	traceItem[5] = s.perf.Count
	traceItem[6] = s.perf.Sum
	traceItem[7] = s.perf.Max
	traceItem[8] = s.perf.Min
	traceItem[9] = s.stackTrace
	return traceItem

}
