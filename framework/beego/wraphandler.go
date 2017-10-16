// Copyright 2017 冯立强 fenglq@tingyun.com.  All rights reserved.
// Use of this source code is governed by a MIT license: https://opensource.org/licenses/MIT

package tingyun_beego

import (
	"net/http"
)

//对于beego.Handler/begoo.NSHandler/begoo.Namespace.Handler 这种方式
//使用context.Context方式的hook 抓不到http.ResponseWriter的 Write和WriteHeader事件
//故此处增加了对这种方式的hook. 需要在被监控代码里替换以上两个函数和一个方法的调用
type handlerWrapper struct {
	http.Handler
	rootpath string
	usepath  bool
}

func (h *handlerWrapper) getPath(r *http.Request) string {
	if !h.usepath {
		return r.RequestURI
	}
	return h.rootpath
}
func (h *handlerWrapper) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	action := getActionByRequest(r)
	wrapRW := createResponseWriter(rw, action)
	if action == nil {
		wrapRequest(r, wrapRW.action)
	}
	wrapRW.action.SetName("URI", h.getPath(r))
	success := false
	defer func() {
		if action == nil {
			unWrapRequest(r) //解wrap,释放内存
		}
		if !success {
			if err := recover(); err != nil {
				wrapRW.action.SetError(err)
				wrapRW.action.Finish()
				wrapRW.action = nil
				panic(err)
			}
		}
	}()
	h.Handler.ServeHTTP(wrapRW, r)
	wrapRW.action.Finish()
	wrapRW.action = nil
	success = true
}
func wrapHandler(path string, h http.Handler, pathUsed bool) http.Handler {
	return &handlerWrapper{Handler: h, rootpath: path, usepath: pathUsed}
}
