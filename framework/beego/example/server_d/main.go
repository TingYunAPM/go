package main

import (
	"fmt"
	"time"

	"github.com/TingYunAPM/go"
	"github.com/TingYunAPM/go/framework/beego"
	"github.com/astaxie/beego"
	"github.com/TingYunAPM/routinelocal"
//如果是通过修改golang代码内建支持协程局部存储的版本,请注释上一行,打开下面一行的注释
//	"github.com/TingYunAPM/routinelocal/native"
)

type MainController struct {
	tingyun_beego.Controller
}

func getTingyunAction() *tingyun.Action {
	return tingyun_beego.RoutineLocalGetAction()
}
func db_component() {
	action := getTingyunAction()
	component := action.CreateDBComponent(tingyun.ComponentMysql, "192.168.100.2", "mydb", "mytable", "select", "db_component")
	time.Sleep(2 * time.Second)
	component.Finish()
}
func (this *MainController) Get() {
	db_component()
	this.Ctx.WriteString("hello world")
}

func main() {
	err := tingyun.AppInit("tingyun.json")
	//注意: 如果要使用 tingyun_beego.RoutineLocalGetAction(),下边这行必须添加
	tingyun_beego.RoutineLocalInit(routinelocal.Get())
	if err != nil {
		fmt.Println(err)
	}
	defer tingyun.AppStop()

	beego.Router("/", &MainController{})
	tingyun_beego.Run()
}
