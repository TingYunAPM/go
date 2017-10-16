package controllers

import (
	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Get() {
	panic("throw a exception")
	c.TplName = "e500.tpl"
}
