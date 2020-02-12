// Copyright 2016-2019 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"
)

func (c *Component) fixEnd(t *timeRange) {
	if c.time.duration == -1 {
		c.time.duration = t.duration - c.time.begin.Sub(t.begin)
		if c.time.duration < 0 {
			c.time.duration = 0
		}
	}
}
func (c *Component) clean_subs() {
	clearComponent(&c.subs)
}
func (c *Component) destroy() {
	c.clean_subs()
	if c._type == componentUnused {
		return
	}
	c.name = ""
	c.method = ""
	c.txdata = ""
	c.callStack = nil
	c.pre = nil
	c.action = nil
	c._type = componentUnused
	//	app.componentTemps.Put(c)
}

//汇总还要修改
func dbMetricName(name string) string {
	array := strings.Split(name, "://")
	serverDb, dburl := array[0], array[1]
	array = strings.Split(dburl, "/")
	host, db, table, op := array[0], array[1], array[2], array[3]
	return "Database " + serverDb + "/" + strings.Replace(url.QueryEscape(host), "%3A", ":", -1) + "%2F" + url.QueryEscape(db) + "%2F" + url.QueryEscape(table) + "/" + url.QueryEscape(op)
}
func nosqlMetricName(name string) string {
	array := strings.Split(name, "://")
	serverDb, dburl := array[0], array[1]
	array = strings.Split(dburl, "/")
	_, _, table, op := array[0], array[1], array[2], array[3]
	return serverDb + "/" + url.QueryEscape(table) + "/" + url.QueryEscape(op)
}
func (c *Component) getSQL() string {
	if c.sql != "" {
		return c.sql
	}
	if c.isDatabaseComponent() {
		return c.metricName()
	}
	return ""
}
func (c *Component) isDatabaseComponent() bool {
	return c._type == ComponentMysql || c._type == ComponentPostgreSql || c._type == ComponentDefaultDB
}
func (c *Component) getURL() string {
	if c._type == ComponentExternal {
		return c.name
	}
	return ""
}
func (c *Component) externalTransaction() (metricName string, duration float64, remote_duration float64, found bool) {
	if len(c.txdata) == 0 {
		return "", 0, 0, false
	}

	jsonData := map[string]interface{}{}
	if err := json.Unmarshal([]byte(c.txdata), &jsonData); err != nil {
		return "", 0, 0, false
	}
	if err, secret_id := jsonReadString(jsonData, "id"); err != nil {
		return "", 0, 0, false
	} else if err, action := jsonReadString(jsonData, "action"); err != nil {
		return "", 0, 0, false
	} else if err, time_object := jsonReadObjects(jsonData, "time"); err != nil {
		return "", 0, 0, false
	} else if err, remote_duration := jsonReadFloat(time_object, "duration"); err != nil {
		return "", 0, 0, false
	} else {
		c.extSecretId = secret_id
		c.remoteDuration = remote_duration
		return "ExternalTransaction/" + strings.Replace(c.name, "/", "%2F", -1) + "/" + secret_id + "%2F" + action, float64(c.time.duration / time.Millisecond), remote_duration, true
	}
}
func (c *Component) metricName() string {
	switch c._type {
	case ComponentExternal:
		return "External/" + strings.Replace(c.name, "/", "%2F", -1) + "/" + url.QueryEscape(c.method)
	case ComponentDefaultDB:
		return dbMetricName(c.name)
	case ComponentMemCache:
		return nosqlMetricName(c.name)
	case ComponentMongo:
		return nosqlMetricName(c.name)
	case ComponentMysql:
		return dbMetricName(c.name)
	case ComponentPostgreSql:
		return dbMetricName(c.name)
	case ComponentRedis:
		return nosqlMetricName(c.name)
	case ComponentDefault:
		return ""
	default:
		return ""
	}
}

func (c *Component) unicId() string {
	if c.exId {
		return unicId(c.time.begin, c)
	}
	return ""
}
func (c *Component) init(component string, method string, _type uint8) *Component {
	c.action = nil
	c.pre = nil
	c.name = component
	c.method = method
	c.txdata = ""
	c.extSecretId = ""
	c.exId = false
	c.callStack = nil
	c.time = timeRange{time.Now(), -1}
	c.sql = ""
	c.aloneTime = 0
	c.remoteDuration = 0
	c.subs.Init()
	c._type = _type
	return c
}
func newComponent(component string, method string, _type uint8) *Component {
	if app == nil {
		return nil
	}
	//	if comp, found := app.componentTemps.Get(); found {
	//		return comp.(*Component).init(component, method, _type)
	//	}
	return (&Component{}).init(component, method, _type)
}
