// Copyright 2017 冯立强 fenglq@tingyun.com.  All rights reserved.
// Use of this source code is governed by a MIT license: https://opensource.org/licenses/MIT

package tingyun_gin

import (
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

// Wrap gin.RouterGroup
type WrapGroup struct {
	g      *gin.RouterGroup
	engine *WrapEngine
}

func (g *WrapGroup) TingYunGinRouterGroup() *gin.RouterGroup {
	return g.g
}

///////////////////////////////////////////////////////////////////////////////
func (g *WrapGroup) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	g.engine.wrap()
	g.TingYunGinRouterGroup().Use(middleware...)
	return g
}

func (g *WrapGroup) Group(relativePath string, handlers ...gin.HandlerFunc) *WrapGroup {
	g.engine.wrap()
	return &WrapGroup{g: g.TingYunGinRouterGroup().Group(relativePath, handlers...), engine: g.engine}
}

func (g *WrapGroup) BasePath() string {
	return g.TingYunGinRouterGroup().BasePath()
}

func (g *WrapGroup) wrap_method(httpMethod, relativePath string, handlers []gin.HandlerFunc) []gin.HandlerFunc {
	g.engine.wrap()
	count := len(handlers)
	arr := make([]gin.HandlerFunc, count+1)
	arr[0] = preHandler(httpMethod, path.Join(g.g.BasePath(), relativePath), count)
	for i, _ := range handlers {
		arr[i+1] = handlers[i]
	}
	return arr
}
func (g *WrapGroup) Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	arr := g.wrap_method(httpMethod, relativePath, handlers)
	g.TingYunGinRouterGroup().Handle(httpMethod, relativePath, arr...)
	return g
}
func (g *WrapGroup) POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	g.Handle("POST", relativePath, handlers...)
	return g
}
func (g *WrapGroup) GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	g.Handle("GET", relativePath, handlers...)
	return g
}
func (g *WrapGroup) DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	g.Handle("DELETE", relativePath, handlers...)
	return g
}
func (g *WrapGroup) PATCH(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	g.Handle("PATCH", relativePath, handlers...)
	return g
}
func (g *WrapGroup) PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	g.Handle("PUT", relativePath, handlers...)
	return g
}
func (g *WrapGroup) OPTIONS(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	g.Handle("OPTIONS", relativePath, handlers...)
	return g
}
func (g *WrapGroup) HEAD(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	g.Handle("HEAD", relativePath, handlers...)
	return g
}
func (g *WrapGroup) Any(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	arr := g.wrap_method("ANY", relativePath, handlers)
	g.TingYunGinRouterGroup().Any(relativePath, arr...)
	return g
}

func (g *WrapGroup) StaticFile(relativePath, filepath string) gin.IRoutes {
	g.engine.wrap()
	ginGroup := g.TingYunGinRouterGroup()
	handler := ginGroup.Handlers
	ginGroup.Use(ignore(relativePath))
	ginGroup.StaticFile(relativePath, filepath)
	ginGroup.Handlers = handler
	return g
}
func (g *WrapGroup) Static(relativePath, root string) gin.IRoutes {
	g.engine.wrap()
	ginGroup := g.TingYunGinRouterGroup()
	handler := ginGroup.Handlers
	ginGroup.Use(ignore(relativePath))
	ginGroup.Static(relativePath, root)
	ginGroup.Handlers = handler
	return g
}
func (g *WrapGroup) StaticFS(relativePath string, fs http.FileSystem) gin.IRoutes {
	g.engine.wrap()
	ginGroup := g.TingYunGinRouterGroup()
	handler := ginGroup.Handlers
	ginGroup.Use(ignore(relativePath))
	ginGroup.StaticFS(relativePath, fs)
	ginGroup.Handlers = handler
	return g
}
