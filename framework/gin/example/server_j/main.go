package main

import (
	"fmt"

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

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.POST("/login", login("/v1"))
		v1.POST("/submit", submit("/v1"))
		v1.POST("/read", read("/v1"))
	}

	// Simple group: v2
	v2 := router.Group("/v2")
	{
		v2.POST("/login", login("/v2"))
		v2.POST("/submit", submit("/v2"))
		v2.POST("/read", read("/v2"))
	}

	router.Run(":8028")
}
func login(parent string) gin.HandlerFunc {
	return func(c *gin.Context) {
		message := c.DefaultPostForm("message", "default_message")
		nick := c.DefaultPostForm("nick", "anonymous")

		c.JSON(200, gin.H{
			"req":     fmt.Sprintf("%s/login", parent),
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})

	}
}
func submit(parent string) gin.HandlerFunc {
	return func(c *gin.Context) {
		message := c.DefaultPostForm("message", "default_message")
		nick := c.DefaultPostForm("nick", "anonymous")

		c.JSON(200, gin.H{
			"req":     fmt.Sprintf("%s/submit", parent),
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	}
}
func read(parent string) gin.HandlerFunc {
	return func(c *gin.Context) {
		message := c.DefaultPostForm("message", "default_message")
		nick := c.DefaultPostForm("nick", "anonymous")

		c.JSON(200, gin.H{
			"req":     fmt.Sprintf("%s/read", parent),
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	}
}
