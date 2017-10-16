package main

import (
	//	"errors"
	"fmt"
	"html"
	"net/http"

	"github.com/TingYunAPM/go"
	"github.com/TingYunAPM/go/framework/beego"
	//	"github.com/astaxie/beego"
)

func main() {
	err := tingyun.AppInit("tingyun.json")
	if err != nil {
		fmt.Println(err)
	}
	defer tingyun.AppStop()
	tingyun_beego.Handler("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		//		panic(errors.New("Panic Test"))
	}))
	tingyun_beego.Run()
}
