package controllers

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/TingYunAPM/go"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {

	c.TplName = "login.tpl"
}

func (c *LoginController) Post() {
	defer tingyun.GetAction(c.Ctx.ResponseWriter.ResponseWriter).CreateComponent("post").Finish()
	time.Sleep(3 * time.Second)
	c.Data["username"] = c.Input().Get("fname") + " " + c.Input().Get("lname")
	c.TplName = "index.tpl"
}
