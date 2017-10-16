package routers

import (
	"github.com/astaxie/beego"
	"github.com/TingYunAPM/go/examples/bee-server/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/error", &controllers.ErrorController{})
}
