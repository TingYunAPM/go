// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"strings"
)

var g_opnames = [5]string{"INSERT", "UPDATE", "SELECT", "DELETE", "CALL"}
var g_opnameMap = map[string]int{"INSERT": 0, "UPDATE": 1, "SELECT": 2, "DELETE": 3, "CALL": 4}

func dbGetIdByOp(op string) int {
	if id, ok := g_opnameMap[strings.ToUpper(op)]; ok {
		return id
	}
	return -1
}

type structSqlPerf struct {
	OpTablePerf [5]mapPerformance
}

func newStructSqlPerf() *structSqlPerf {
	r := &structSqlPerf{}
	for i := 0; i < 5; i++ {
		r.OpTablePerf[i] = newMapPerformance()
	}
	return r
}

func (s *structSqlPerf) Merge(t *structSqlPerf) {
	for i := 0; i < 5; i++ {
		s.OpTablePerf[i].Merge(t.OpTablePerf[i])
	}
}
func (s *structSqlPerf) Add(table string, opId int, perf float64, excl float64) {
	s.OpTablePerf[opId].ExclAdd(table, perf, excl)
}
