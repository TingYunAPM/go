# tingyun_gin

## 安装

- 运行

```
go get github.com/TingYunAPM/go
```

## 使用

- 引用: 
main函数文件中引入
```
import "github.com/TingYunAPM/go"
import "github.com/TingYunAPM/go/framework/gin"
```
- router初始化: 
在main函数中初始化tingyun agent,用tingyun_gin.New 和tingyun_gin.Default 替换gin.New 和gin.Default
```
tingyun.AppInit("tingyun.json")
defer tingyun.AppStop()
//router := gin.Default()
router := tingyun_gin.Default()
```
## 获取Action
- 在handler函数中获取tingyun.Action对象
```
func handler(c *gin.Context){
    action := tingyn_gin.FindAction(c)
    //...
}
```

## 使用Component
```
func handler(c *gin.Context){
    action := tingyn_gin.FindAction(c)
    componentCheck := action.CreateComponent("CheckJob")
    //do check works
    componentCheck.Finish()
    componentJSON := action.CreateComponent("gin.Context::JSON")
    c.JSON(200, gin.H{"message":"pong",})
    componentJSON.Finish()
}

```
## 其他
请参考 https://github.com/TingYunAPM/go/blob/master/README.md
