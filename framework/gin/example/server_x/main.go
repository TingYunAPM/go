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
	//	r := gin.New()
	r := tingyun_gin.New()
	r.Use(Logger())

	r.GET("/test", func(c *gin.Context) {
		example := c.MustGet("example").(string)

		// it would print: "12345"
		fmt.Println(example)
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8042")
}
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		fmt.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		fmt.Println(status)
	}
}
