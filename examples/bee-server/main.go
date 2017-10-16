package main

import (
	"github.com/astaxie/beego"
	"github.com/TingYunAPM/go"
	_ "github.com/TingYunAPM/go/examples/bee-server/routers"
)

func main() {
	tingyun.AppInit("tingyun.json")
	defer tingyun.AppStop()

	beego.Run()
}
