package main

import (
	"errors"
	"fmt"
	"html"
	"net/http"

	"github.com/TingYunAPM/go"
	"github.com/TingYunAPM/go/framework/beego"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type AdminController struct {
	tingyun_beego.Controller
}

func (c *AdminController) Init(ctx *context.Context, controllerName, actionName string, app interface{}) {
	fmt.Printf("AdminController.Init\n")
	c.Controller.Init(ctx, controllerName, actionName, app)
}
func (this *AdminController) Get() {
	this.Ctx.WriteString("AdminController")
}
func (this *AdminController) ShowAPIVersion() {
	this.Ctx.WriteString("AdminController ShowAPIVersion")
}

type UserController struct {
	tingyun_beego.Controller
}

func (this *UserController) Get() {
	this.Ctx.WriteString("UserController")
}

type MainController struct {
	tingyun_beego.Controller
}

func (this *MainController) Get() {
	this.Ctx.WriteString("MainController")
}

type CMSController struct {
	tingyun_beego.Controller
}

func (this *CMSController) Get() {
	this.Ctx.WriteString("CMSController")
}

type BlockController struct {
	tingyun_beego.Controller
}

func (this *BlockController) Get() {
	this.Ctx.WriteString("BlockController")
}

func main() {
	err := tingyun.AppInit("tingyun.json")
	if err != nil {
		fmt.Println(err)
	}
	defer tingyun.AppStop()
	ns := beego.NewNamespace("/v1",
		beego.NSCond(func(ctx *context.Context) bool {
			if ctx.Input.Domain() == "127.0.0.1" {
				return true
			}
			return false
		}),
		beego.NSGet("/notallowed", func(ctx *context.Context) {
			ctx.Output.Body([]byte("notAllowed"))
		}),
		tingyun_beego.NSHandler("/handler", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		})),
		beego.NSRouter("/version", &AdminController{}, "get:ShowAPIVersion"),
		beego.NSRouter("/changepassword", &UserController{}),
		beego.NSNamespace("/shop",
			beego.NSGet("/:id", func(ctx *context.Context) {
				ctx.Output.Body([]byte("notAllowed"))
			}),
		),
		beego.NSNamespace("/cms",
			beego.NSInclude(
				&MainController{},
				&CMSController{},
				&BlockController{},
			),
		),
	)
	tingyun_beego.NamespaceHandler(ns, "/ttt", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ttt Hello, %q", html.EscapeString(r.URL.Path))
		panic(errors.New("Panic Test"))
	}))
	beego.AddNamespace(ns)
	tingyun_beego.Run()
}
