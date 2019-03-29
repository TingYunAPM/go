// Copyright 2016-2019 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/TingYunAPM/go/utils/list"
)

type Component struct {
	action    *Action
	pre       *Component
	name      string
	method    string
	txdata    string
	exId      bool
	callStack interface{}
	time      timeRange
	aloneTime time.Duration
	subs      list.List
	sql       string
	_type     uint8
}

func (c *Component) GetAction() *Action {
	if c == nil {
		return nil
	}
	return c.action
}

// 停止性能分解组件计时
// 性能分解组件时长 = Finish时刻 - CreateComponent时刻
// 当时长超出堆栈阈值时，记录当前组件的代码堆栈
func (c *Component) Finish() {
	if c != nil && c.time.duration == -1 {
		c.time.End()
		if c._type != ComponentDefault {
			c.aloneTime = c.time.duration
		}
		if readServerConfigBool(configServerConfigBoolActionTracerEnabled, true) {
			//超阈值取callstack
			if c.time.duration >= time.Duration(readServerConfigInt(configServerConfigIntegerActionTracerStacktraceThreshold, 500))*time.Millisecond {
				c.callStack = callStack(1)
			}
		}
	}
}

//跨应用追踪接口,用于调用端,生成一个跨应用追踪id,通过http头或者私有协议发送到被调用端
//
//返回值: 字符串,一个包含授权id,应用id,实例id,事务id等信息的追踪id
func (c *Component) CreateTrackId() string {
	if app == nil || c == nil || c.action == nil || c._type != ComponentExternal {
		return ""
	}
	if enabled := readServerConfigBool(configServerConfigBoolTransactionTracerEnabled, false); !enabled {
		return ""
	}
	//TINGYUN_ID_SECRET;c=CALLER_TYPE;r=REQ_ID;x=TX_ID;e=EXTERNAL_ID;p=PROTOCOL
	//时间+对象地址=>生成exId
	if secId := app.configs.server.CStrings.Read(configServerStringTingyunIdSecret, ""); len(secId) == 0 {
		return ""
	} else {
		c.exId = true
		protocol := "http"
		if arr := strings.Split(c.name, "://"); len(arr) > 1 {
			protocol = arr[0]
		}
		return secId + ";c=1;x=" + c.action.unicId() + ";e=" + c.unicId() + ";p=" + protocol
	}
}

//跨应用追踪接口,用于调用端,将被调用端返回的事务性能数据保存到外部调用组件
//
//参数: 被调用端返回的事务的性能数据
func (c *Component) SetTxData(txData string) {
	if app == nil || c == nil || c.action == nil || c._type != ComponentExternal {
		return
	}
	jsonData := map[string]interface{}{}
	if err := json.Unmarshal([]byte(txData), &jsonData); err != nil {
		return
	}
	if err, tr := jsonReadInt(jsonData, "tr"); err == nil {
		c.action.track_enable = (tr != 0)
	}
	c.txdata = txData
}

//用于数据库组件,通过此接口将sql查询语句保存到数据库组件,在报表慢事务追踪列表展示
//
//参数: sql语句
func (c *Component) AppendSQL(sql string) {
	if app == nil || c == nil || c.action == nil || (c._type != ComponentExternal && c._type != ComponentDefaultDB && c._type != ComponentMysql && c._type != ComponentPostgreSql) {
		return
	}
	c.sql = sql
}
func (c *Component) CreateComponent(method string) *Component {
	if c == nil || c.action == nil || c._type != ComponentDefault {
		return nil
	}
	metric := makeComponentMetricName(method)
	r := c.action.createComponent(c._type, metric, url.QueryEscape(method))
	r.pre = c
	return r
}
