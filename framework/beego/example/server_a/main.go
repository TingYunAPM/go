package main

import (
	"fmt"

	"github.com/TingYunAPM/go"
	"github.com/TingYunAPM/go/framework/beego"
	"github.com/astaxie/beego"
)

type MainController struct {
	tingyun_beego.Controller
}

func (this *MainController) Get() {
	this.Ctx.WriteString("hello world")
}

func main() {
	err := tingyun.AppInit("tingyun.json")
	if err != nil {
		fmt.Println(err)
	}
	defer tingyun.AppStop()

	beego.Router("/", &MainController{})
	tingyun_beego.Run()
}
