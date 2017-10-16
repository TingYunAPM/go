// Copyright 2017 冯立强 fenglq@tingyun.com.  All rights reserved.
// Use of this source code is governed by a MIT license: https://opensource.org/licenses/MIT

package tingyun_beego

import (
	//	"errors"
	"fmt"
	"html"
	"net/http"

	"github.com/TingYunAPM/go"
	"github.com/TingYunAPM/go/framework/beego"
	//	"github.com/astaxie/beego"
)

func ExampleHandler() {
	err := tingyun.AppInit("tingyun.json")
	if err != nil {
		fmt.Println(err)
	}
	defer tingyun.AppStop()
	//"beego.Handler" 替换为:=> "tingyun_beego.Handler"
	tingyun_beego.Handler("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	}))
	//"beego.Run" 替换为:=> "tingyun_beego.Run"
	tingyun_beego.Run()
}

type MainController struct {
	tingyun_beego.Controller
}

func (this *MainController) Get() {
	this.Ctx.WriteString("hello world")
}
func ExampleController() {
	err := tingyun.AppInit("tingyun.json")
	if err != nil {
		fmt.Println(err)
	}
	defer tingyun.AppStop()
	//"beego.Controller" 替换为:=> "tingyun_beego.Controller"
	//type MainController struct {
	//	tingyun_beego.Controller
	//}
	//func (this *MainController) Get() {
	//	this.Ctx.WriteString("hello world")
	//}
	beego.Router("/", &MainController{})
	//"beego.Run" 替换为:=> "tingyun_beego.Run"
	tingyun_beego.Run()
}
func ExampleNSHandler() {
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
		//"beego.NSHandler" 替换为:=> "tingyun_beego.NSHandler"
		tingyun_beego.NSHandler("/handler", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		})),
	)
	beego.AddNamespace(ns)
	//"beego.Run" 替换为:=> "tingyun_beego.Run"
	tingyun_beego.Run()
}
func ExampleNamespaceHandler() {
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
	)
	//"beego.Namespace.Handler" 替换为:=> "tingyun_beego.NamespaceHandler"
	//ns.Handler("/ttt", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	tingyun_beego.NamespaceHandler(ns, "/ttt", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ttt Hello, %q", html.EscapeString(r.URL.Path))
	}))

	beego.AddNamespace(ns)
	//"beego.Run" 替换为:=> "tingyun_beego.Run"
	tingyun_beego.Run()
}
func Example() {
	err := tingyun.AppInit("tingyun.json")
	if err != nil {
		fmt.Println(err)
	}
	defer tingyun.AppStop()
	//"beego.Run" 替换为:=> "tingyun_beego.Run"
	tingyun_beego.Run()
}
