// Copyright 2017 冯立强 fenglq@tingyun.com.  All rights reserved.
// Use of this source code is governed by a MIT license: https://opensource.org/licenses/MIT

//beego's Wrapper
package tingyun_beego

import (
	"net/http"

	"github.com/TingYunAPM/go"
	"github.com/astaxie/beego/context"
)

//API: 在应用过程中
//对于没有使用beego.Handler, beego.NSHandler 函数和 beego.Namespace.Handler方法 的方式
//使用FindAction获取对应的tingyun.Action
func FindAction(ctx *context.Context) *tingyun.Action {
	return GetAction(ctx.ResponseWriter.ResponseWriter)
}

//API: 对于使用beego.Handler, beego.NSHandler函数 或 beego.Namespace.Handler方法的方式
//使用GetAction获取对应的tingyun.Action
//在没有hook的情况下,返回nil
func GetAction(rw http.ResponseWriter) *tingyun.Action {
	if !tingyun.Running() {
		return nil
	}
	if p, ok := rw.(*responseWriter); ok {
		return p.action
	}
	return nil
}
