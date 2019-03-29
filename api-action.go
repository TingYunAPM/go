// Copyright 2016-2019 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/TingYunAPM/go/utils/pool"
)

const (
	ComponentDefault    = 0
	ComponentDefaultDB  = 32
	ComponentMysql      = 33
	ComponentPostgreSql = 34
	ComponentMongo      = 48
	ComponentMemCache   = 49
	ComponentRedis      = 50
	ComponentExternal   = 64
	componentUnused     = 255
)

var dbNameMap = [32]string{0: "Database", 1: "Mysql", 2: "PostgreSql", 16: "MongoDB", 17: "MemCache", 18: "Redis"}

const (
	actionUsing    = 1
	actionFinished = 2
	actionUnused   = 0
)

type Action struct {
	name          string
	url           string
	method        string
	trackId       string
	track_enable  bool
	cache         pool.Pool
	errors        pool.Pool
	statusCode    uint16
	requestParams map[string]string
	customParams  map[string]interface{}
	time          timeRange
	stateUsed     uint8
}

//创建Web Service性能分解组件
//参数:
//    url    : 调用Web Service的url,格式: http(s)://host/uri, 例如 http://www.baidu.com/
//    method : 发起这个Web Service调用的类名.方法名, 例如 http.Get
func (a *Action) CreateExternalComponent(url string, method string) *Component {
	if a == nil || a.stateUsed != actionUsing {
		return nil
	}
	return a.createComponent(ComponentExternal, url, method)
}

//创建数据库或NOSQL性能分解组件
//参数:
//    dbType : 组件类型 (ComponentMysql, ComponentPostgreSql, ComponentMongo, ComponentMemCache, ComponentRedis)
//    host   : 主机地址，可空
//    dbname : 数据库名称，可空
//    table  : 数据库表名
//    op     : 操作类型, 关系型数据库("SELECT", "INSERT", "UPDATE", "DELETE" ...), NOSQL("GET", "SET" ...)
//    method : 发起这个数据库调用的类名.方法名, 例如 db.query redis.get
func (a *Action) CreateDBComponent(dbType uint8, host string, dbname string, table string, op string, method string) *Component {
	if a == nil || a.stateUsed != actionUsing {
		//fmt.Println("CreateDBComponent object is nil!")
		return nil
	}
	nameId := dbType - ComponentDefaultDB
	if nameId < 0 || nameId >= 32 {
		return nil
	}
	protocol := dbNameMap[nameId]
	if protocol == "" {
		protocol = "UnDefDatabase"
	}
	if dbname == "" {
		dbname = "(NULL)"
	}
	if host == "" {
		host = "NULL"
	}
	if table == "" {
		table = "NULL"
	}
	return a.createComponent(dbType, fmt.Sprintf("%s://%s/%s/%s/%s", protocol, host, dbname, table, op), method)
}
func makeComponentMetricName(method string) string {
	className, methodName := parseMethod(method)
	if className == "" {
		className = "NULL"
	}
	return "Go/" + url.QueryEscape(className) + "/" + url.QueryEscape(methodName)
}

//创建性能分解组件，作用为将一个HTTP请求拆分为多个可以度量的组件
//参数
//    method : 类名.方法名, 例如 main.user.login
func (a *Action) CreateComponent(method string) *Component {
	if a == nil || a.stateUsed != actionUsing {
		return nil
	}
	metric := makeComponentMetricName(method)
	return a.createComponent(ComponentDefault, metric, url.QueryEscape(method))
}
func (a *Action) AddRequestParam(k string, v string) {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	a.requestParams[k] = v
}
func (a *Action) AddCustomParam(k string, v string) {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	a.customParams[k] = v
}

//跨应用追踪接口,用于被调用端,获取当前事务的执行性能信息,通过http头或者自定义协议传回调用端
//
//返回值: 事务的性能数据
func (a *Action) GetTxData() string {
	if a == nil || a.stateUsed != actionUsing || len(a.trackId) == 0 {
		return ""
	}
	if enabled := readServerConfigBool(configServerConfigBoolTransactionTracerEnabled, false); !enabled {
		return ""
	}
	tr := 0
	if a.track_enable || a.Slow() {
		tr = 1
	}
	instance_id := ""
	if secId := app.configs.server.CStrings.Read(configServerStringTingyunIdSecret, ""); len(secId) == 0 {
		return ""
	} else {
		if parts := strings.Split(secId, "|"); len(parts) < 2 {
			return ""
		} else {
			instance_id = strings.Split(parts[1], ";")[0]
		}
	}
	curr_time := time.Now()
	duration_ := a.time.GetDuration(curr_time)
	var db_time time.Duration = 0
	var mc_time time.Duration = 0
	var mongo_time time.Duration = 0
	var redis_time time.Duration = 0
	var external_time time.Duration = 0
	var code_time time.Duration = 0
	a.cache.ForEach(func(v interface{}) {
		component := v.(*Component)
		duration := component.time.GetDuration(curr_time)
		if component._type == ComponentDefault {
			code_time += duration
		} else if component._type == ComponentMysql || component._type == ComponentDefaultDB || component._type == ComponentPostgreSql {
			db_time += duration
		} else if component._type == ComponentExternal {
			external_time += duration
		} else if component._type == ComponentMemCache {
			mc_time += duration
		} else if component._type == ComponentMongo {
			mongo_time += duration
		} else if component._type == ComponentRedis {
			redis_time += duration
		}
	})
	time_obj := map[string]interface{}{
		"duration": duration_ / time.Millisecond,
		"qu":       0,
	}
	if db_time > 0 {
		time_obj["db"] = db_time / time.Millisecond
	}
	if mc_time > 0 {
		time_obj["mc"] = mc_time / time.Millisecond
	}
	if mongo_time > 0 {
		time_obj["mon"] = mongo_time / time.Millisecond
	}
	if redis_time > 0 {
		time_obj["rds"] = redis_time / time.Millisecond
	}
	if external_time > 0 {
		time_obj["ex"] = external_time / time.Millisecond
	}
	if code_time > 0 {
		time_obj["code"] = code_time / time.Millisecond
	}
	res := map[string]interface{}{
		"id":     instance_id,
		"action": a.GetName(),
		"time":   time_obj,
		"tr":     tr,
	}
	if tr > 0 {
		res["trId"] = unicId(a.time.begin, a)
	}
	if json_byte, err := json.Marshal(res); err == nil {
		return string(json_byte)
	}
	return ""
}

//跨应用追踪接口,用于被调用端,保存调用端传递过来的跨应用追踪id
//
//参数: 跨应用追踪id
func (a *Action) SetTrackId(id string) {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	if enabled := readServerConfigBool(configServerConfigBoolTransactionTracerEnabled, false); !enabled {
		return
	}
	//解析id,匹配用户

	if secId := app.configs.server.CStrings.Read(configServerStringTingyunIdSecret, ""); len(secId) == 0 {
		return
	} else {
		clientUser := strings.Split(id, "|")[0]
		if len(clientUser) == 0 {
			return
		}
		localUser := strings.Split(secId, "|")[0]
		if localUser != clientUser {
			return
		}
		a.trackId = id
	}
}
func formatActionName(instance string, method string) string {
	if len(instance) == 0 {
		instance = "Go"
	}
	mlen := len(method)
	if mlen > 1 && method[0:1] == "/" {
		method = method[1:mlen]
	}
	return "WebAction/" + url.QueryEscape(instance) + "/" + url.QueryEscape(method)
}

//设置HTTP请求的友好名称
//参数:
//    instance   : 分类, 例如 loginController
//    method : 方法, 例如 POST
func (a *Action) SetName(instance string, method string) {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	a.name = formatActionName(instance, method)
}
func (a *Action) GetName() string {
	if a == nil {
		return ""
	}
	return a.name
}
func (a *Action) GetUrl() string {
	if a == nil {
		return ""
	}
	return a.url
}
func (a *Action) SetUrl(name string) {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	a.url = name
}

//返回当前HTTP请求是否为慢请求
//返回值: 当HTTP请求性能超出阈值时为true, 否则为false
func (a *Action) Slow() bool {
	if a == nil {
		return false
	}
	enabled := readServerConfigBool(configServerConfigBoolActionTracerEnabled, true)
	if !enabled {
		return false
	}
	if a.stateUsed == actionUnused {
		return false
	}
	threshold := readServerConfigInt(configServerConfigIntegerActionTracerActionThreshold, 500)
	if a.stateUsed == actionUsing {
		return time.Now().Sub(a.time.begin) >= time.Duration(threshold)*time.Millisecond
	} else if a.stateUsed == actionFinished {
		return a.time.duration >= time.Duration(threshold)*time.Millisecond
	}
	return false
}

//不统计此http请求的性能数据
func (a *Action) Ignore() {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	a.destroy()
}

func (a *Action) HasError() bool {
	if a == nil || a.stateUsed != actionUsing {
		return false
	}
	return a.errors.Size() > 0
}

//记录HTTP请求的运行时错误信息
func (a *Action) SetError(e interface{}) {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	a.setError(e, "RUNTIME_ERROR")
}

//停止HTTP请求的性能计时
//HTTP请求时长 = Finish时刻 - CreateAction时刻
func (a *Action) Finish() {
	if a == nil || a.stateUsed != actionUsing {
		return
	}
	a.stateUsed = actionFinished
	if a.statusCode == 0 {
		a.statusCode = 200
	}
	append_action(a)
}

//正常返回0
//无效的Action 返回1
//如果状态码是被忽略的错误,返回2
//错误的状态码,返回3
func (a *Action) SetStatusCode(code uint16) int {
	if a == nil || a.stateUsed != actionUsing {
		return 1
	}
	if a.statusCode == 0 {
		a.statusCode = code
	}
	if code >= 400 && code != 401 { //401认证失败，非错误代码
		ignoredStatusCodes := app.configs.server_arrays.Read(configServerConfigIArrayIgnoredStatusCodes, nil)
		if ignoredStatusCodes != nil {
			for _, v := range ignoredStatusCodes {
				if uint16(v) == code {
					return 2
				}
			}
		}
		a.setError(errors.New(fmt.Sprint("status code ", code)), "HTTP_ERROR")
		return 3
	}
	return 0
}
func agentEnabled() bool {
	if app == nil {
		return false
	} else if config_disabled {
		return false
	} else {
		return app.configs.server.CBools.Read(configServerBoolEnabled, true)
	}
}
func CreateAction(instance string, method string) (*Action, error) {
	//fmt.Printf("CreateAction(%s, %s)\n", instance, method)
	if app == nil {
		if config_disabled {
			return nil, errors.New("Agent disabled by local config file.")
		}
		return nil, errors.New("Agent not Inited, please call AppInit() first.")
	} else if app.actionPool.Size() > int32(app.configs.local.CIntegers.Read(configLocalIntegerNbsActionCacheMax, 10000)) {
		return nil, errors.New("Server busy, Skip one action.")
	}
	return app.createAction(instance, method)
}
func (a *Action) destroy() {
	if a == nil || a.stateUsed == actionUnused {
		return
	}
	a.name = ""
	a.url = ""
	a.trackId = ""
	for component, _ := a.cache.Get(); component != nil; component, _ = a.cache.Get() {
		component.(*Component).destroy()
	}
	for err, _ := a.errors.Get(); err != nil; err, _ = a.errors.Get() {
	}
	if a.requestParams != nil {
		for k := range a.requestParams {
			delete(a.requestParams, k)
		}
	}
	if a.customParams != nil {
		for k := range a.customParams {
			delete(a.customParams, k)
		}
	}
	a.stateUsed = actionUnused
	a.statusCode = 0
}
