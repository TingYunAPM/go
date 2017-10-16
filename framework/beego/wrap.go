// Copyright 2017 冯立强 fenglq@tingyun.com.  All rights reserved.
// Use of this source code is governed by a MIT license: https://opensource.org/licenses/MIT

package tingyun_beego

import (
	"fmt"
	"net/http"

	"github.com/TingYunAPM/go"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func finishAction(ctx *context.Context) {
	if p, ok := ctx.ResponseWriter.ResponseWriter.(*responseWriter); ok {
		if p.isStatic || !p.writed { //不采集静态路由性能
			p.action.Ignore()
		} else {
			p.action.Finish()
		}
		unWrapRequest(ctx.Request)
		unWrapContext(ctx)
	}
}

var wrapperInited = false

//hook context.Context.ResponseWriter.ResponseWriter
//通过这个hook的responseWriter捕获事件采集数据
func wrapContext(ctx *context.Context) {
	if _, ok := ctx.ResponseWriter.ResponseWriter.(*responseWriter); !ok {
		rs_writer := createResponseWriter(ctx.ResponseWriter.ResponseWriter, nil)
		ctx.ResponseWriter.ResponseWriter = rs_writer
		fmt.Printf("Context Create action %p\n", rs_writer.action)

		wrapRequest(ctx.Request, rs_writer.action)
	}
	FindAction(ctx).SetName(ctx.Request.Method, ctx.Request.RequestURI)
}

//应用过程结束时解除hook,释放内存,还原原来的结构
func unWrapContext(ctx *context.Context) {
	if p, ok := ctx.ResponseWriter.ResponseWriter.(*responseWriter); ok {
		ctx.ResponseWriter.ResponseWriter = p.ResponseWriter
		p.action = nil
	}
}

//初始化,插入filter,
//抓取panic
func beegoInit() {
	if wrapperInited {
		return
	}
	wrapperInited = true
	beego.InsertFilter("/*", beego.BeforeStatic, wrapContext)
	beego.InsertFilter("/*", beego.BeforeRouter, func(ctx *context.Context) {
		if action := FindAction(ctx); action != nil {
			p, _ := ctx.ResponseWriter.ResponseWriter.(*responseWriter)
			p.isStatic = false
		}
	})
	beego.InsertFilter("/*", beego.AfterExec, finishAction, false)
	beego.InsertFilter("/*", beego.FinishRouter, finishAction, false)
	//抓panic
	beego.BConfig.RecoverFunc = wrapPanic(beego.BConfig.RecoverFunc)
}
func wrapPanic(pre_recover func(*context.Context)) func(*context.Context) {
	return func(ctx *context.Context) {
		if pre_recover != nil {
			defer pre_recover(ctx)
		}
		if err := recover(); err != nil {
			action := FindAction(ctx)
			action.SetError(err)
			action.Finish()
			unWrapContext(ctx)
			panic(err)
		}
	}
}

//替换beego.Handler
func Handler(rootpath string, h http.Handler, options ...interface{}) *beego.App {
	if !tingyun.Running() {
		return beego.Handler(rootpath, h, options...)
	}
	return beego.Handler(rootpath, wrapHandler(rootpath, h, true), options...)
}

//替换beego.Run
func Run(params ...string) {
	if tingyun.Running() {
		beegoInit()
	}
	beego.Run(params...)
}
