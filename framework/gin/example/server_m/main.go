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
	route.Any("/testing", startPage)
	route.Run(":8031")
}

func startPage(c *gin.Context) {
	var person Person
	if c.BindQuery(&person) == nil {
		fmt.Println("====== Only Bind By Query String ======")
		fmt.Println(person.Name)
		fmt.Println(person.Address)
	}
	c.String(200, "Success")
}
