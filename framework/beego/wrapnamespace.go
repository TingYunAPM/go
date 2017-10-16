// Copyright 2017 冯立强 fenglq@tingyun.com.  All rights reserved.
// Use of this source code is governed by a MIT license: https://opensource.org/licenses/MIT

package tingyun_beego

import (
	"net/http"

	"github.com/TingYunAPM/go"
	"github.com/astaxie/beego"
)

//替换beego.NSHandler
func NSHandler(rootpath string, h http.Handler) beego.LinkNamespace {
	if !tingyun.Running() {
		return beego.NSHandler(rootpath, h)
	}
	return beego.NSHandler(rootpath, wrapHandler(rootpath, h, false))
}

//替换beego.Namesapce.Handler
func NamespaceHandler(n *beego.Namespace, rootpath string, h http.Handler) *beego.Namespace {
	if !tingyun.Running() {
		return n.Handler(rootpath, h)
	}
	return n.Handler(rootpath, wrapHandler(rootpath, h, false))
}
