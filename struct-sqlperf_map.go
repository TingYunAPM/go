// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"fmt"
	"net/url"
	"strings"
)

type mapSqlPerf map[string]*structSqlPerf

func newMapSqlPerf() mapSqlPerf {
	return mapSqlPerf{}
}

func (p mapSqlPerf) Read(onElement func(string, *structPerformance)) {
	for dbHost, v := range p {
		for i := 0; i < 5; i++ {
			opName := g_opnames[i]
			v.OpTablePerf[i].Read(func(name string, perf *structPerformance) {
				onElement(fmt.Sprintf("Database %s/%s/%s", url.QueryEscape(dbHost), name, opName), perf)
			})
		}
	}
}
func (p mapSqlPerf) Reset() {
	for k := range p {
		delete(p, k)
	}
}

func (p mapSqlPerf) Add(name string, perf float64, excl float64) {
	//name : Mysql://host/dbname/tablename/op
	array := strings.Split(name, "://")
	serverDb, dburl := array[0], array[1]
	array = strings.Split(dburl, "/")
	_, _, table, op := array[0], array[1], array[2], array[3]
	opId := dbGetIdByOp(op)
	if opId == -1 { //无效OP类型
		return
	}
	v, ok := p[serverDb]
	if !ok {
		v = newStructSqlPerf()
		p[serverDb] = v
	}
	v.Add(table, opId, perf, excl)
}

func (p mapSqlPerf) Merge(q mapSqlPerf) {
	for k, v := range q {
		if s, ok := p[k]; ok {
			s.Merge(v)
		} else {
			p[k] = v
		}
	}
}
