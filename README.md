# go
听云 Go Agent SDK

# 安装
go get 下载源码
```
go get github.com/TingYunAPM/go
```
# 使用探针
## 引用探针
```
import "github.com/TingYunAPM/go"
```
## 探针初始化
  Go探针使用json文件获取配置信息，在main函数开始处初始化探针。
  ```
func main() {
	//初始化tingyun: 应用名称、帐号等在tingyun.json中配置
	tingyun.AppInit("tingyun.json")
	defer tingyun.AppStop()
	//原有业务逻辑
	...
}
```
## 创建Action
### Action说明
  应用性能分解过程中,我们使用Action定义一个完整事务，通常它对应的是一个完整的http请求过程。
  ### 代码
  ```
action, _ := tingyun.CreateAction("URI", "/index")
defer action.Finish()
```
## 创建Component
### Component说明
  一个事务通常会包含多个子过程，子过程还可能由其他子过程组成。我们将这样的子过程定义为Component,通过对Component树的耗时分析来定位事务执行过程中的性能瓶颈。
  ### 从Action创建Component
  ```
  component_mysub := action.CreateComponent("my_submethod")
  ```
  创建一个数据库的Component
  ```
  mytable_select := action.CreateDBComponent(tingyun.ComponentMysql, "", "mydatabase", "mytable", "select", "database_method")
  ```
  创建一个外部调用(rpc,http等)的Component
  ```
  external := action.CreateExternalComponent("http://tingyun.com/", "rpc_method")
  ```
### 从Component创建Component  
```
component_next := component_mysub.CreateComponent("my_nextsub")
```


子过程结束时,需要调用对应的Component.Finish(),才能达到采集数据的目的

## 代码示例
```
package main
import (
	"io"
	"net/http"
	"github.com/TingYunAPM/go"
)
func handler(w http.ResponseWriter, r *http.Request) {
	action, _ := tingyun.CreateAction("URI", r.URL.Path)
	defer action.Finish()
	header := w.Header()
	headerComponent := action.CreateComponent("header")
	header.Set("Cache-Control", "no-cache")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	action.SetStatusCode(uint16(http.StatusOK))
	headerComponent.Finish()
	bodyComponent := action.CreateComponent("body")
	io.WriteString(w, "helloworld.")
	bodyComponent.Finish()
}
func main() {
	tingyun.AppInit("tingyun.json")
	defer tingyun.AppStop()
	http.HandleFunc("/pf", handler)
	http.ListenAndServe(":8000", nil)
}
```

## 配置项说明
```
{
  "nbs.app_name" : "go app",
  "nbs.agent_enabled" : true,
  "nbs.license_key" : "999-999-999",
  "nbs.log_file_name" : "agent.log",
  "nbs.audit" : true,
  "nbs.max_log_count": 5,
  "nbs.max_log_size": 10,
  "nbs.ssl" : true,
  "nbs.savecount" : 2
}
```
### "nbs.app_name:
    由探针使用者自定义的监控APP的名字
###  "nbs.agent_enabled" :
    探针禁用标志
### "nbs.license_key":
	使用者的license,登陆tingyun.com获得
### "nbs.log_file_name" :
	日志文件路径
### "nbs.audit":
	审计模式,开启审计模式,会在日志文件中输出到日志分析服务器的上行和下行数据。
### "nbs.max_log_count":
	日志文件最大个数
### "nbs.max_log_size":
	日志文件最大M字节数,此处为10M，超过10M日志就滚动打包成.tar.gz文件
### "nbs.ssl" :
	上传数据是否启用安全套接字
### "nbs.savecount" :
	采样数据上传失败时,在探针端保留的采样数据个数,用于应对网络故障

## 跨应用
* 应用拓扑
当一个帐号下存在多个应用的相互调用关系时, 可以利用API追踪应用之间的调用关系。
调用者CreateTrackId, 被调用者SetTrackId, 在报表内就会产生“调用者” -> "被调用者" 的拓扑图。
thrift 客户端
```
url := "thrift://192.168.1.5/login"
c := tingyun.GetAction(w).CreateExternalComponent(url, "thrift.login")
//调用CreateTrackId生成调用信息
track := c.CreateTrackId()
//将track传递给thrift服务器
thriftLogin(url, track)
//
c.Finish()
```
thrift 服务器端
```
//从thrift数据内获取track, 并SetTrackId, 生成调用关系
track := getTrack(r)
tingyun.GetAction(w).SetTrackId(track)
```
* 跨应用追踪
当产生拓扑关系的应用过程性能超过阈值时，会产生慢过程跟踪数据，同时在慢过程跟踪数据内会记录调用者和被调用者的详细追踪信息。
通过点击慢过程跟踪图表内的链接，可以跳转到被调用者的详细追踪数据。

# 框架支持

## http标准库
如果您使用了`http`标准库, 听云SDK提供了一个封装库来简化您的嵌码工作量, 
封装库在HTTP request进入时自动创建并对应唯一一个应用过程(Action),
在HTTP response结束时自动结束对应的应用过程(Action)。

* 如果您使用`HandleFunc`方式, 请做如下替换:
```
http.HandleFunc("/login", loginHandler)
```
替换为
```
tingyun.HandleFunc("/login", loginHandler)
```
* 如果您使用`Handle`方式, 请做如下替换:
```
http.Handle("/login", http.HandlerFunc(loginHandler))
```
替换为
```
tingyun.Handle("/login", http.HandlerFunc(loginHandler))
```
* 如果您使用`Handler`方式, 请做如下替换:
```
Server.Handler = app.Handlers
```
替换为
```
Server.Handler = tingyun.WrapHandler(app.Handlers)
```
在loginHandler内, 您可以通过`tingyun.GetAction`获取HTTP请求上下文内的应用过程
```
func loginHandler(w http.ResponseWriter, r *http.Request) {
	//增加下面一行代码（可选，根据您的需求变化可能是其他代码）
	defer tingyun.GetAction(w).CreateComponent("loginHandler").Finish()
	//原有业务逻辑
	...
}
```

## gin
参见 https://github.com/TingYunAPM/go/blob/master/framework/gin/README.md
## beego
参见 https://github.com/TingYunAPM/go/blob/master/framework/beego/README.md

# API
## func AppInit(jsonFile string) error 
	探针初始化,参数jsonFile是配置文件路径。成功返回 nil,失败返回error
## func AppStop() 
	探针停止,结束数据采集。之后的除AppInit的所有API调用都将不会产生数据，同时也不会有性能损失。
## func ConfigRead(name string) (interface{}, bool) 
	读取探针配置项,如果应用需要关注探针的设置，则调用这个方法。一般情况下不需要。
## func Log() *log.Logger 
	返回一个具有日志输出功能的对象,在这个对象上调用日志输出方法,将输出到配置文件里设定的日志文件里。一般情况下应用应该有自己的日志模块，不需要调用这个方法。
## func Running() bool 
	检测探针是否在运行
## type Action struct
### func (a *Action) AddCustomParam(k string, v string) 
	将uri携带的key-value参数传递给Action对象,协助应用在报表端分析慢过程追踪和错误追踪栈
### func (a *Action) AddRequestParam(k string, v string) 
	将http头携带的key-value 参数传递给Action对象,协助应用在报表端分析慢过程追踪和错误追踪栈
### func (a *Action) CreateComponent(method string) *Component 
	创建性能分解组件，作用为将一个HTTP请求拆分为多个可以度量的组件。 参数一般设置为对应的类名+方法名,格式="ClassName.MethodName",或者用户自己定义的过程名。
### func (a *Action) CreateDBComponent(dbType uint8, host string, dbname string, table string, op string, method string) *Component
	创建数据库或NOSQL性能分解组件 参数: 
		dbType : 组件类型 (ComponentMysql, ComponentPostgreSql, ComponentMongo, ComponentMemCache, ComponentRedis)
		host   : 主机地址，可空
		dbname : 数据库名称，可空
		table  : 数据库表名
		op     : 操作类型, 关系型数据库("SELECT", "INSERT", "UPDATE", "DELETE" ...), NOSQL("GET", "SET" ...)
		method : 发起这个数据库调用的类名.方法名, 例如 db.query redis.get 
### func (a *Action) CreateExternalComponent(url string, method string) *Component 
	创建Web Service性能分解组件 参数: 
		url    : 调用Web Service的url,格式: http(s)://host/uri, 例如 http://www.tingyun.com/
		method : 发起这个Web Service调用的类名.方法名, 例如 http.Get 
### func (a *Action) GetName() string
	读取Action格式化过的名字,通常不需要使用者调用
### func (a *Action) GetUrl() string 
	读取Action对应的URI
### func (a *Action) HasError() bool 
	一般用来检测是否调用过Action.SetError,或者是Action.SetStatusCode()调用传递大于等于400的状态码后,也会返回true
### func (a *Action) Ignore() 
	忽略这次Action的性能采集,只在Action.Finish()调用之前有效。
### func (a *Action) SetError(e interface{}) 
	运行时记录Action事务过程中的错误信息,将在报表端展示
### func (a *Action) SetName(instance string, method string) 
	修改Action的友好名称 参数: 
		instance   : 分类, 例如 loginController
		method : 方法, 例如 POST 
### func (a *Action) SetStatusCode(code uint16) 
	设置Action对应的http事务应答的http状态码,缺省值200
### func (a *Action) SetTrackId(id string) 
	用于rpc调用或者http外部调用的跨应用追踪,id 由使用者从调用端传过来。
### func (a *Action) SetUrl(name string) 
	设置对应Action的 uri, 用于追踪慢过程和错误分析。
### func (a *Action) Slow() bool
	根据用户的配置和当前事务的值行时间,判定对应的Action是否为慢过程.
### func (a *Action) Finish() 
	当前事务数据采集结束
## func CreateAction(instance string, method string) (*Action, error) 
	创建一个事务对应的Action对象
		instance   : 分类, 例如 loginController
		method : 方法, 例如 POST
	举例: CreateAction("URI", "/index")
		或 CreateAction("MyController", "POST")
## type Component struct
	性能分解组件
### func (c *Component) CreateComponent(method string) *Component
	对本组件再进行性能分解，创建下层性能分解组件，参数同Action.CreateComponent
### func (c *Component) CreateTrackId() string 
	用于跨应用追踪,本组件内调用了外部应用过程或者发起了rpc调用,由应用此方法返回的结果携带到server端,server端通过Action.SetTrackId使用这个结果。最终在报表端生成跨应用追踪图表
### func (c *Component) GetAction() *Action 
	取本组件所属的事务对象Action
### func (c *Component) Finish() 
	停止性能分解组件计时 性能分解组件时长 = Finish时刻 - CreateComponent时刻 当时长超出堆栈阈值时，记录当前组件的代码堆栈 
## func Handle(pattern string, handler http.Handler) 
	对http.Handle方法的包装,内部调用了CreateAction和 Action.Finish
## func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) 
	对http.HandleFunc的包装,内部调用了CreateAction和 Action.Finish
## func WrapHandler(h http.Handler) http.Handler 
	封装http.Handler对象 例如： 初始代码为 app.Server.Handler = app.Handlers 应该替换为 app.Server.Handler = tingyun.WrapHandler(app.Handlers)
## func GetAction(w http.ResponseWriter) *Action
	对于 tingyun.Handle或tingyun.HandleFunc包装过的调用,应用处理过程的http.RewponseWriter参数已经被重新包装,通过此方法可以获取对应的Action,方便调用Action其他方法。
