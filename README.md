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
# 框架支持
## gin
参见 https://github.com/TingYunAPM/go/framework/gin/README.md
## beego
参见 https://github.com/TingYunAPM/go/framework/beego/README.md
