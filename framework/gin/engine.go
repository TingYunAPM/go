// Copyright 2017 冯立强 fenglq@tingyun.com.  All rights reserved.
// Use of this source code is governed by a MIT license: https://opensource.org/licenses/MIT

//gin's wrapper
package tingyun_gin

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Wrap gin.Engine
type WrapEngine struct {
	WrapGroup
	e      *gin.Engine
	wraped bool
}

func New() *WrapEngine {
	ret := &WrapEngine{e: gin.New(), wraped: false}
	ret.WrapGroup.g = &ret.e.RouterGroup
	ret.WrapGroup.engine = ret
	return ret
}
func Default() *WrapEngine {
	ret := &WrapEngine{e: gin.Default(), wraped: false}
	ret.WrapGroup.g = &ret.e.RouterGroup
	ret.WrapGroup.engine = ret
	return ret
}
func (e *WrapEngine) wrap() {
	if !e.wraped {
		e.e.Use(profiler())
		e.wraped = true
	}
}
func (e *WrapEngine) TingYunGinEngine() *gin.Engine {
	return e.e
}

/////////////////////////////////////////////////////////////////////////////
func (e *WrapEngine) Delims(left, right string) *WrapEngine {
	e.e.Delims(left, right)
	return e
}

func (e *WrapEngine) SecureJsonPrefix(prefix string) *WrapEngine {
	e.e.SecureJsonPrefix(prefix)
	return e
}

func (e *WrapEngine) LoadHTMLGlob(pattern string) {
	e.e.LoadHTMLGlob(pattern)
}

func (e *WrapEngine) LoadHTMLFiles(files ...string) {
	e.e.LoadHTMLFiles(files...)
}

func (e *WrapEngine) SetHTMLTemplate(templ *template.Template) {
	e.e.SetHTMLTemplate(templ)
}

func (e *WrapEngine) SetFuncMap(funcMap template.FuncMap) {
	e.e.SetFuncMap(funcMap)
}

func (e *WrapEngine) NoRoute(handlers ...gin.HandlerFunc) {
	e.e.NoRoute(handlers...)
}

func (e *WrapEngine) NoMethod(handlers ...gin.HandlerFunc) {
	e.e.NoMethod(handlers...)
}

func (e *WrapEngine) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	e.wrap()
	e.e.Use(middleware...)
	return e
}

func (e *WrapEngine) Routes() gin.RoutesInfo {
	return e.e.Routes()
}

func (e *WrapEngine) Run(addr ...string) error {
	e.wrap()
	return e.e.Run(addr...)
}

func (e *WrapEngine) RunTLS(addr string, certFile string, keyFile string) error {
	e.wrap()
	return e.e.RunTLS(addr, certFile, keyFile)
}
func (e *WrapEngine) RunUnix(file string) error {
	return e.e.RunUnix(file)
}

func (e *WrapEngine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.wrap()
	e.e.ServeHTTP(w, req)
}
func (e *WrapEngine) HandleContext(c *gin.Context) {
	e.e.HandleContext(c)
}
