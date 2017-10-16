package main

import (
	"fmt"
	"time"

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
	// Creates a router without any middleware by default
	//	r := gin.New()
	r := tingyun_gin.New()

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Per route middleware, you can add as many as you desire.
	r.GET("/benchmark", MyBenchLogger(), benchmark())

	authorized := r.Group("/")
	authorized.Use(AuthRequired())
	{
		authorized.POST("/login", login("/"))
		authorized.POST("/submit", submit("/"))
		authorized.POST("/read", read("/"))

		testing := authorized.Group("testing")
		testing.GET("/analytics", analytics("/testing"))
	}
	r.Run(":8029")
}
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func MyBenchLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		fmt.Println("Used :", time.Now().Sub(start))
	}
}
func benchmark() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.JSON(200, gin.H{
			"req": fmt.Sprintf("/benchmark"),
		})

	}
}
func analytics(parent string) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.JSON(200, gin.H{
			"req": fmt.Sprintf("%s/analytics", parent),
		})

	}
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

//curl -d "message=message2&nick=fengliqiang" http://127.0.021:8029/login -v
//curl -d "message=message2&nick=fengliqiang" http://127.0.021:8029/submit -v
//curl -d "message=message2&nick=fengliqiang" http://127.0.021:8029/read -v
