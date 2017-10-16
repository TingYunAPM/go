package main

import (
	"github.com/hprose/hprose-go/hprose"
	"github.com/TingYunAPM/go/examples/bee-rpc/models"

	"github.com/astaxie/beego"
	"github.com/TingYunAPM/go"
)

func main() {
	//初始化tingyun: 应用名称、license等在tingyun.json中配置
	tingyun.AppInit("tingyun.json")
	defer tingyun.AppStop()

	service := hprose.NewHttpService()
	service.AddFunction("AddOne", models.AddOne)
	service.AddFunction("GET", models.GetOne)
	beego.Handler("/", service)

	beego.Run()
}
