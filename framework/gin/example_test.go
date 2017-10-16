// Copyright 2017 冯立强 fenglq@tingyun.com.  All rights reserved.
// Use of this source code is governed by a MIT license: https://opensource.org/licenses/MIT
package tingyun_gin

import (
	"fmt"

	"github.com/TingYunAPM/go"
	"github.com/TingYunAPM/go/framework/gin"
	"github.com/gin-gonic/gin"
)

func ExampleDefault() {
	tingyun.AppInit("tingyun.json") //初始化tingyun探针
	defer tingyun.AppStop()         //退出时关闭探针
	r := tingyun_gin.Default()      //这里替换掉原来的gin.Default
	r.GET("/ping", func(c *gin.Context) {
		component := tingyun_gin.FindAction(c).CreateComponent("gin.Context::JSON")
		c.JSON(200, gin.H{
			"message": "pong",
		})
		component.Finish()
	})
	r.Run(":8000")
}

func ExampleNew() {
	tingyun.AppInit("tingyun.json") //初始化tingyun探针
	defer tingyun.AppStop()         //退出时关闭探针
	r := tingyun_gin.New()          //这里替换掉原来的gin.New
	r.GET("/ping", func(c *gin.Context) {
		component := tingyun_gin.FindAction(c).CreateComponent("gin.Context::JSON")
		c.JSON(200, gin.H{
			"message": "pong",
		})
		component.Finish()
	})
	r.Run(":8000")
}
