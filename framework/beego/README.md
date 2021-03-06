# tingyun_beego

## 安装

- 运行
```
go get github.com/TingYunAPM/go
```

## 使用
- ### 引用: 
 main函数文件中引入
```
import "github.com/TingYunAPM/go"
import "github.com/TingYunAPM/go/framework/beego"
```
- ### 函数和方法替换:
1. "beego.Run()" 替换为:=> "tingyun_beego.Run()"
例:
```
func main() {
  tingyun.AppInit("tingyun.json")
  defer tingyun.AppStop()
  //Do Other Things...
  tingyun_beego.Run()//replace beego.Run()
}
```
2. "beego.Handler" 替换为:=> "tingyun_beego.Handler"
例:
```
func main() {
    tingyun.AppInit("tingyun.json")
    defer tingyun.AppStop()
    //beego.Handler("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    tingyun_beego.Handler("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
     //Do Some Things
        fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
    }))
    tingyun_beego.Run()
}
```

3.  "beego.Controller" 替换为:=> "tingyun_beego.Controller"
例:
```
type MainController struct {
    //beego.Controller
    tingyun_beego.Controller
}
func (this *MainController) Get() {
    this.Ctx.WriteString("hello world")
}
func main() {
    tingyun.AppInit("tingyun.json")
    defer tingyun.AppStop()
    beego.Router("/", &MainController{})
    tingyun_beego.Run()
}
```
4. "beego.NSHandler" 替换为:=> "tingyun_beego.NSHandler"
例:
```
func main() {
    tingyun.AppInit("tingyun.json")
    defer tingyun.AppStop()
    beego.AddNamespace(beego.NewNamespace("/v1",
        //
        // beego.NSHandler("/handler", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tingyun_beego.NSHandler("/handler", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
        })),
    ))
    tingyun_beego.Run()
}
```
5. "beego.Namespace.Handler" 替换为:=> "tingyun_beego.NamespaceHandler"
例:
```
func main() {
    tingyun.AppInit("tingyun.json")
    defer tingyun.AppStop()
    nsobj := beego.NewNamespace("/v1")
    //
    //nsobj.Handler("/ttt", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    tingyun_beego.NamespaceHandler(ns, "/ttt", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "ttt Hello, %q", html.EscapeString(r.URL.Path))
    }))
    beego.AddNamespace(nsobj)
    tingyun_beego.Run()
}
```

## 获取Action,使用Component
- 在Controller中获取tingyun.Action对象
```
type MainController struct {
    tingyun_beego.Controller
}
func (this *MainController) Get() {
    action := tingyun_beego.FindAction(this.Ctx)
    componentCheck := action.CreateComponent("CheckJob")
    //Do Some Check Works
    componentCheck.Finish()
    componentWrite := action.CreateComponent("Get::out")
    this.Ctx.WriteString("hello world")
    componentWrite.Finish()
}
```
- 在http.Handler中获取tingyun.Action对象
```
tingyun_beego.Handler("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    action := tingyun_beego.GetAction(w)
    componentCheck := action.CreateComponent("CheckJob")
    //
    //Do Some Check Works
    componentCheck.Finish()
    componentWrite := action.CreateComponent("Handler::out")
    fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
    componentWrite.Finish()
}))
```

- 使用协程局部存储方式保存/获取tingyun.Action对象
```
//go get github.com/TingYunAPM/routinelocal 安装协程局部存储包
//在main函数文件 import "github.com/TingYunAPM/routinelocal"
//或者 import "github.com/TingYunAPM/routinelocal/native" //参考 github.com/TingYunAPM/routinelocal/native/Readme.md
//在main函数开始使用routinelocal.Storage初始化 tingyun_beego,实际用例参考 github.com/TingYunAPM/go/framework/beego/example/server_d
func main() {
	err := tingyun.AppInit("tingyun.json")
	//注意: 如果要使用 tingyun_beego.RoutineLocalGetAction(),下边这行必须添加
	tingyun_beego.RoutineLocalInit(routinelocal.Get()
	...
//在某个函数中要获取tingyun.Action, 使用: tingyun_beego.RoutineLocalGetAction()
```

## 其他
请参考 https://github.com/TingYunAPM/go/blob/master/README.md
