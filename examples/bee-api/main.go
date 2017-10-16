package main

import (
	_ "github.com/TingYunAPM/go/examples/bee-api/routers"

	"github.com/astaxie/beego"

	"github.com/TingYunAPM/go"
)

//http://127.0.0.1:8080/v1/user/user_11111

func main() {
	//初始化tingyun: 应用名称、license等在tingyun.json中配置
	tingyun.AppInit("tingyun.json")
	defer tingyun.AppStop()

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
