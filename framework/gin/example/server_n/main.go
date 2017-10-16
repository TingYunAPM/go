package main

import (
	"fmt"

	"github.com/TingYunAPM/go"
	"github.com/TingYunAPM/go/framework/gin"
	"github.com/gin-gonic/gin"
)

type Person struct {
	Name    string `form:"name"`
	Address string `form:"address"`
}

func main() {
	err := tingyun.AppInit("tingyun.json")
	if err != nil {
		fmt.Println(err)
	}
	defer tingyun.AppStop()
	//	route := gin.Default()
	route := tingyun_gin.Default()
	route.GET("/testing", startPage)
	route.POST("/testing", startPage)
	route.Run(":8032")
}

func startPage(c *gin.Context) {
	var person Person
	// If `GET`, only `Form` binding engine (`query`) used.
	// If `POST`, first checks the `content-type` for `JSON` or `XML`, then uses `Form` (`form-data`).
	// See more at https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48
	if c.Bind(&person) == nil {
		fmt.Println(person.Name)
		fmt.Println(person.Address)
	}

	c.String(200, "Success")
}

//curl -d "name=liqiang&address=tongzhou" http://127.0.0.1:8032/testing -v
