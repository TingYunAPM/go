// Copyright 2017 冯立强 fenglq@tingyun.com.  All rights reserved.
// Use of this source code is governed by a MIT license: https://opensource.org/licenses/MIT

package tingyun_beego

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/TingYunAPM/go"
)

//包装http.ResponseWriter,用来关联tingyun.Action
type responseWriter struct {
	http.ResponseWriter
	action   *tingyun.Action
	isStatic bool
	writed   bool
}

func (w *responseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

//抓取Write事件,default status code = 200
func (w *responseWriter) Write(s []byte) (int, error) {
	if !w.writed {
		w.action.SetStatusCode(200)
	}
	w.writed = true
	return w.ResponseWriter.Write(s)
}

//抓取状态码
func (w *responseWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.action.SetStatusCode(uint16(status))
	fmt.Printf("WriteHeader:%d\n", status)
	w.writed = true
	if status/100 > 3 && status != 401 {
		w.action.Finish()
	}
}

func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("webserver doesn't support hijacking")
	}
	return hj.Hijack()
}

func (w *responseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *responseWriter) CloseNotify() <-chan bool {
	if cn, ok := w.ResponseWriter.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	return nil
}

func (w *responseWriter) init(pre http.ResponseWriter, a *tingyun.Action) *responseWriter {
	w.ResponseWriter = pre
	w.isStatic = true
	if a == nil {
		w.action, _ = tingyun.CreateAction("URI", "/")
	} else {
		w.action = a
	}
	w.writed = false
	return w
}

//起因: hook方式有两种,一种是Filter,一种是对http.Handler的包装,
//  beego内部会copy http.ResponseWriter 指针到 context.Context内部
//  对context.Context的hook 不会延伸到 beego.Handler方式中的 http.ResponseWriter
//  所以两边都要Hook, 但是Action要保持一个，不能重复创建
//  Action的唯一性通过 hook http.Request.Body来保证
//  这里传递的tingyun.Action就是要保证对应的Web过程要有唯一的Action
func createResponseWriter(pre http.ResponseWriter, a *tingyun.Action) *responseWriter {
	r := (&responseWriter{}).init(pre, a)
	return r
}
