// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
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
	exId      bool
	callStack interface{}
	time      timeRange
	aloneTime time.Duration
	subs      list.List
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
func (c *Component) CreateComponent(method string) *Component {
	if c == nil || c.action == nil || c._type != ComponentDefault {
		return nil
	}
	metric := makeComponentMetricName(method)
	r := c.action.createComponent(c._type, metric, url.QueryEscape(method))
	r.pre = c
	return r
}
