// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"encoding/json"
	"strings"
	"time"
)

func (a *Action) setError(e interface{}, errType string) {
	if a == nil || a.stateUsed != actionUsing {
		return
	} //errorTrace 聚合,以 callstack + message
	errTime := time.Now()
	callstackBytes, _ := json.Marshal(callStack(1))
	a.errors.Put(&errInfo{errTime, e, string(callstackBytes), errType})
}
func (a *Action) createComponent(componentType uint8, name string, method string) *Component {
	//fmt.Printf("createComponent(%d,%s,%s)\n", componentType, name, method)
	if a == nil || a.stateUsed != actionUsing {
		return nil
	}
	r := newComponent(name, method, componentType)
	a.cache.Put(r)
	r.action = a
	return r
}
func getTxId(id string) string {
	array := strings.Split(id, ";")
	if len(array) < 4 {
		return ""
	}
	for i := 0; i < len(array); i++ {
		paire := strings.Split(array[i], "=")
		if paire[0] == "x" {
			if len(paire) < 2 {
				return ""
			} else {
				xid := paire[1]
				for i := 2; i < len(paire); i++ {
					xid += "=" + paire[i]
				}
				return xid
			}
		}
	}
	return ""
}
func getTopMetric(id string) string {
	array := strings.Split(id, ";")
	if len(array) < 4 {
		return ""
	}
	for i := 1; i < len(array); i++ {
		paire := strings.Split(array[i], "=")
		if paire[0] == "p" {
			if len(paire) < 2 {
				return ""
			} else {
				protocol := paire[1]
				for i := 2; i < len(paire); i++ {
					protocol += "=" + paire[i]
				}
				return "EntryTransaction/" + protocol + "/" + getAppSecId(array[0])
			}
		}
	}
	return ""
}

func (a *Action) unicId() string {
	txId := getTxId(a.trackId)
	if txId == "" {
		return unicId(a.time.begin, a)
	}
	return txId
}
func (a *Action) reinit(instance string, method string) *Action {
	a.stateUsed = actionUsing
	a.SetName(instance, method)
	a.method = method
	a.time.Init()
	a.statusCode = 0
	return a
}

func (a *Action) init(name string) *Action {
	a.cache.Init()
	a.errors.Init()
	a.name = name
	a.url = ""
	a.trackId = ""
	a.statusCode = 0
	a.requestParams = make(map[string]string)
	a.customParams = make(map[string]interface{})
	a.time.Init()
	a.stateUsed = actionUsing
	return a
}
