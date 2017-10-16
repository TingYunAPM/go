// Copyright 2017 冯立强 fenglq@tingyun.com.  All rights reserved.
// Use of this source code is governed by a MIT license: https://opensource.org/licenses/MIT

package tingyun_beego

import (
	"github.com/TingYunAPM/go"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

//要使用tingyun_beego.Controller 替换 beego.Controller
type Controller struct {
	beego.Controller
}

func (c *Controller) Init(ctx *context.Context, controllerName, actionName string, app interface{}) {
	if tingyun.Running() {
		action := FindAction(ctx)
		action.SetName(controllerName, actionName)
	}
	c.Controller.Init(ctx, controllerName, actionName, app)
}
func (c *Controller) Finish() {
	if tingyun.Running() {
		finishAction(c.Ctx)
	}
	c.Controller.Finish()
}
