package main

import (
	"fmt"
	"net/http"

	"github.com/TingYunAPM/go"
	"github.com/TingYunAPM/go/framework/gin"
	"github.com/gin-gonic/gin"
)

func main() {
	err := tingyun.AppInit("tingyun.json")
	if err != nil {
		fmt.Println(err)
	}
	defer tingyun.AppStop()
	//	router := gin.Default()
	router := tingyun_gin.Default()
	router.LoadHTMLGlob("templates/**/*")
	router.GET("/posts/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
			"title": "Posts",
		})
	})
	router.GET("/users/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
			"title": "Users",
		})
	})
	router.Run(":8039")
}
