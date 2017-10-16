package tingyun

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strings"
)

// 1. For example, to instrument this code:
//
//    http.Handle("/login", http.HandlerFunc(loginHandler))
//
// Perform this replacement:
//
//    tingyun.Handle("/login", http.HandlerFunc(loginHandler))
//
//
// 2. For example, to instrument this code:
//
//    http.HandleFunc("/login", loginHandler)
//
// Perform this replacement:
//
//    tingyun.HandleFunc("/login", loginHandler)
//
//
// 3. For example, to instrument this code:
//
//	  //handlers Implement http.Handler interface.
//    Server.Handler = handlers
//
// Perform this replacement:
//
//    Server.Handler = tingyun.WrapHandler(handlers)
//
// The Action is passed to the handler in place of the original
// http.ResponseWriter, so it can be accessed using type assertion.
//
//  func loginHandler(w http.ResponseWriter, r *http.Request) {
//      action := tingyun.GetAction(w)
//	    action.SetName("bee", "login")
//		...
//  }

type wrapIf interface {
	http.ResponseWriter
	Context() *Action
	Finish()
}

type wrapObj struct {
	wrapIf
	action  *Action
	writer  http.ResponseWriter
	request *http.Request
}

const (
	hasC = 1 << iota // CloseNotifier
	hasF             // Flusher
	hasH             // Hijacker
	hasR             // ReaderFrom
)

type wrap struct{ *wrapObj }
type wrapR struct{ *wrapObj }
type wrapH struct{ *wrapObj }
type wrapHR struct{ *wrapObj }
type wrapF struct{ *wrapObj }
type wrapFR struct{ *wrapObj }
type wrapFH struct{ *wrapObj }
type wrapFHR struct{ *wrapObj }
type wrapC struct{ *wrapObj }
type wrapCR struct{ *wrapObj }
type wrapCH struct{ *wrapObj }
type wrapCHR struct{ *wrapObj }
type wrapCF struct{ *wrapObj }
type wrapCFR struct{ *wrapObj }
type wrapCFH struct{ *wrapObj }
type wrapCFHR struct{ *wrapObj }

func (x wrapC) CloseNotify() <-chan bool    { return x.writer.(http.CloseNotifier).CloseNotify() }
func (x wrapCR) CloseNotify() <-chan bool   { return x.writer.(http.CloseNotifier).CloseNotify() }
func (x wrapCH) CloseNotify() <-chan bool   { return x.writer.(http.CloseNotifier).CloseNotify() }
func (x wrapCHR) CloseNotify() <-chan bool  { return x.writer.(http.CloseNotifier).CloseNotify() }
func (x wrapCF) CloseNotify() <-chan bool   { return x.writer.(http.CloseNotifier).CloseNotify() }
func (x wrapCFR) CloseNotify() <-chan bool  { return x.writer.(http.CloseNotifier).CloseNotify() }
func (x wrapCFH) CloseNotify() <-chan bool  { return x.writer.(http.CloseNotifier).CloseNotify() }
func (x wrapCFHR) CloseNotify() <-chan bool { return x.writer.(http.CloseNotifier).CloseNotify() }

func (x wrapF) Flush()    { x.writer.(http.Flusher).Flush() }
func (x wrapFR) Flush()   { x.writer.(http.Flusher).Flush() }
func (x wrapFH) Flush()   { x.writer.(http.Flusher).Flush() }
func (x wrapFHR) Flush()  { x.writer.(http.Flusher).Flush() }
func (x wrapCF) Flush()   { x.writer.(http.Flusher).Flush() }
func (x wrapCFR) Flush()  { x.writer.(http.Flusher).Flush() }
func (x wrapCFH) Flush()  { x.writer.(http.Flusher).Flush() }
func (x wrapCFHR) Flush() { x.writer.(http.Flusher).Flush() }

func (x wrapH) Hijack() (net.Conn, *bufio.ReadWriter, error) { return x.writer.(http.Hijacker).Hijack() }
func (x wrapHR) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return x.writer.(http.Hijacker).Hijack()
}
func (x wrapFH) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return x.writer.(http.Hijacker).Hijack()
}
func (x wrapFHR) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return x.writer.(http.Hijacker).Hijack()
}
func (x wrapCH) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return x.writer.(http.Hijacker).Hijack()
}
func (x wrapCHR) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return x.writer.(http.Hijacker).Hijack()
}
func (x wrapCFH) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return x.writer.(http.Hijacker).Hijack()
}
func (x wrapCFHR) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return x.writer.(http.Hijacker).Hijack()
}

func (x wrapR) ReadFrom(r io.Reader) (int64, error)    { return x.writer.(io.ReaderFrom).ReadFrom(r) }
func (x wrapHR) ReadFrom(r io.Reader) (int64, error)   { return x.writer.(io.ReaderFrom).ReadFrom(r) }
func (x wrapFR) ReadFrom(r io.Reader) (int64, error)   { return x.writer.(io.ReaderFrom).ReadFrom(r) }
func (x wrapFHR) ReadFrom(r io.Reader) (int64, error)  { return x.writer.(io.ReaderFrom).ReadFrom(r) }
func (x wrapCR) ReadFrom(r io.Reader) (int64, error)   { return x.writer.(io.ReaderFrom).ReadFrom(r) }
func (x wrapCHR) ReadFrom(r io.Reader) (int64, error)  { return x.writer.(io.ReaderFrom).ReadFrom(r) }
func (x wrapCFR) ReadFrom(r io.Reader) (int64, error)  { return x.writer.(io.ReaderFrom).ReadFrom(r) }
func (x wrapCFHR) ReadFrom(r io.Reader) (int64, error) { return x.writer.(io.ReaderFrom).ReadFrom(r) }

func structToInterface(obj *wrapObj) wrapIf {
	if obj == nil {
		return nil
	}
	x := 0
	if _, ok := obj.writer.(http.CloseNotifier); ok {
		x |= hasC
	}
	if _, ok := obj.writer.(http.Flusher); ok {
		x |= hasF
	}
	if _, ok := obj.writer.(http.Hijacker); ok {
		x |= hasH
	}
	if _, ok := obj.writer.(io.ReaderFrom); ok {
		x |= hasR
	}

	switch x {
	default:
		return wrap{obj}
	case hasR:
		return wrapR{obj}
	case hasH:
		return wrapH{obj}
	case hasH | hasR:
		return wrapHR{obj}
	case hasF:
		return wrapF{obj}
	case hasF | hasR:
		return wrapFR{obj}
	case hasF | hasH:
		return wrapFH{obj}
	case hasF | hasH | hasR:
		return wrapFHR{obj}
	case hasC:
		return wrapC{obj}
	case hasC | hasR:
		return wrapCR{obj}
	case hasC | hasH:
		return wrapCH{obj}
	case hasC | hasH | hasR:
		return wrapCHR{obj}
	case hasC | hasF:
		return wrapCF{obj}
	case hasC | hasF | hasR:
		return wrapCFR{obj}
	case hasC | hasF | hasH:
		return wrapCFH{obj}
	case hasC | hasF | hasH | hasR:
		return wrapCFHR{obj}
	}
}

func (a *wrapObj) Context() *Action { return a.action }

func (a *wrapObj) Header() http.Header { return a.writer.Header() }

func (a *wrapObj) Write(b []byte) (int, error) { return a.writer.Write(b) }

func (a *wrapObj) WriteHeader(code int) {
	a.writer.WriteHeader(code)
	a.action.SetStatusCode(uint16(code))
}

func createWrapper(instance string, pattern string, w http.ResponseWriter, req *http.Request) *wrapObj {
	action, _ := CreateAction(instance, pattern)
	if action == nil {
		return nil
	}
	action.url = req.URL.RequestURI()
	r := new(wrapObj)
	r.action = action
	r.request = req
	r.writer = w
	return r
}

func (a *wrapObj) Finish() {
	if a.action != nil {
		if a.action.Slow() || a.action.HasError() {
			a.action.AddCustomParam("referer", a.request.Referer())
			a.action.AddCustomParam("user-agent", a.request.UserAgent())
			a.action.AddCustomParam("IP", a.request.RemoteAddr)
			for k, v := range a.request.Form {
				s := strings.Join(v, ",")
				if len(s) > 128 {
					s = s[:128] + "..."
				}
				a.action.AddRequestParam(k, s)
			}
		}
		a.action.Finish()
	}
	a.action = nil
	a.writer = nil
	a.request = nil
}

func WrapHandle(pattern string, handler http.Handler) (string, http.Handler) {
	return pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i := structToInterface(createWrapper("Handler", pattern, w, r))
		if i != nil {
			i.Context().SetName("Handler", getnameByAddr(handler))
			defer i.Finish()
			handler.ServeHTTP(i, r)
		} else {
			//agent is disabled
			handler.ServeHTTP(w, r)
		}
	})
}

//封装http.Handle
//例如：
//初始代码为 http.Handle("/login", http.HandlerFunc(loginHandler))
//应该替换为 tingyun.Handle("/login", http.HandlerFunc(loginHandler))
func Handle(pattern string, handler http.Handler) {
	http.Handle(WrapHandle(pattern, handler))
}

func WrapHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
	p, h := WrapHandle(pattern, http.HandlerFunc(handler))
	return p, func(w http.ResponseWriter, r *http.Request) { h.ServeHTTP(w, r) }
}

//封装http.HandleFunc函数
//例如：
//初始代码为 http.HandleFunc("/login", loginHandler)
//应该替换为 tingyun.HandleFunc("/login", loginHandler)
func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(WrapHandleFunc(pattern, handler))
}

//
//  modify --> beego/app.go
//  func (app *App) Run(){
//      ...
//      app.Server.Handler = tingyun.WrapHandler(app.Handlers)
//  }

//  modify -->  beego/controller.go
//  func (c *Controller) Init(ctx *context.Context, controllerName, actionName string, app interface{}) {
//		...
//		tingyun.GetAction(ctx.ResponseWriter.ResponseWriter).SetName(controllerName, actionName)
//  }

//  modify -->  beego/config.go
//  func recoverPanic(ctx *context.Context) {
//  	if err := recover(); err != nil {
//  		tingyun.GetAction(ctx.ResponseWriter.ResponseWriter).SetError(err)
//			...
//      }
//  }

type handlerObj struct {
	handler http.Handler
}

func (h *handlerObj) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i := structToInterface(createWrapper("URI", r.URL.Path, w, r))
	if i != nil {
		defer i.Finish()
		h.handler.ServeHTTP(i, r)
	} else {
		//agent is disabled
		h.handler.ServeHTTP(w, r)
	}
}

//封装http.Handler对象
//例如：
//初始代码为 app.Server.Handler = app.Handlers
//应该替换为 app.Server.Handler = tingyun.WrapHandler(app.Handlers)
func WrapHandler(h http.Handler) http.Handler {
	return &handlerObj{handler: h}
}

// 从封装的http.ResponseWriter中获取tingyun上下文
// 参数:
//    http.ResponseWriter:  被tingyun.wrapper封装的http.ResponseWriter
//
// 例如
// func loginHandler(w http.ResponseWriter, r *http.Request) {
//    action := tingyun.GetAction(w)
//    action.SetName("bee", "login")
// }
func GetAction(w http.ResponseWriter) *Action {
	if wrap, ok := w.(wrapIf); ok {
		return wrap.Context()
	}
	return nil
}
