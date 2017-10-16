// Copyright 2017 冯立强 fenglq@tingyun.com.  All rights reserved.
// Use of this source code is governed by a MIT license: https://opensource.org/licenses/MIT

package tingyun_gin

import (
	"github.com/TingYunAPM/go"
	"github.com/gin-gonic/gin"
)

//gin's middleware,采集Action性能数据
func profiler() gin.HandlerFunc {
	return func(c *gin.Context) {
		//检测是否已经创建Action(用户多次插入中间件)
		action := FindAction(c)
		if action == nil {
			action, _ = tingyun.CreateAction("URI", c.Request.URL.Path)
			if action != nil {
				c.Set(context_id, action)
				defer func() {
					if action == nil { //No Panic?
						return
					}
					//Panic! Action::SetError
					if err := recover(); err != nil {
						action.SetError(err)
						action.Finish()
						c.Set(context_id, nil)
						panic(err) //重新抛出,恢复应用Panic处理过程
					}
				}()
			}
		}
		c.Next()
		c.Set(context_id, nil)
		status := c.Writer.Status()
		action.SetStatusCode(uint16(status))
		if status/100 >= 4 {
			querys := c.Request.URL.Query()
			for k, v := range querys {
				value := ""
				if len(v) > 0 {
					value = v[0]
				}
				action.AddRequestParam(k, value)
			}
			headers := c.Request.Header
			for k, v := range headers {
				value := ""
				if len(v) > 0 {
					value = v[0]
				}
				action.AddCustomParam(k, value)
			}

		}
		action.Finish()
		action = nil //Panic flag
	}
}

//静态路由,不捕获
func ignore(relativePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		//必须是已经创建Action
		action := FindAction(c)
		c.Next()
		if action != nil {
			status := c.Writer.Status()
			if status/100 < 4 {
				action.Ignore()
			}
		}
	}
}

//添加到 GET/HEAD/POST/PUT/OPTIONS/DELETE/PATCH handler 的前边
func preHandler(httpMethod, relativePath string, handler_count int) gin.HandlerFunc {
	return func(c *gin.Context) {
		action := FindAction(c)
		action.SetName(httpMethod, relativePath)
	}
}
