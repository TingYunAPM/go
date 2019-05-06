package main

import (
	"fmt"
	"time"

	"github.com/TingYunAPM/go"
	"github.com/TingYunAPM/go/framework/beego"
	"github.com/astaxie/beego"
	"github.com/TingYunAPM/routinelocal"
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
	tingyun_beego.RoutineLocalInit(routinelocal.Get())
	if err != nil {
		fmt.Println(err)
	}
	defer tingyun.AppStop()

	beego.Router("/", &MainController{})
	tingyun_beego.Run()
}
