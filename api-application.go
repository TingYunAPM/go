// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

//听云性能采集探针(sdk)
package tingyun

//面向api使用者的接口实现部分

import (
	"errors"

	"github.com/TingYunAPM/go/utils/logger"
)

//初始化听云探针
//参数:
//    jsonFile: 听云配置文件路径，文件格式为json格式
func AppInit(jsonFile string) error {
	if app != nil {
		return errors.New("Agent already inited!")
	}
	if r, err := new(application).init(jsonFile); err != nil {
		return err
	} else {
		app = r
		return nil
	}
}

//检测探针是否启动(为Frameworks提供接口)
//返回值: bool
func Running() bool {
	return app != nil
}

//停止听云探针
func AppStop() {
	if app == nil {
		return
	}
	app.stop()
	app = nil
}

func ConfigRead(name string) (interface{}, bool) {
	if app == nil {
		return nil, false
	}
	return app.configs.Value(name)
}

func Log() *log.Logger {
	if app == nil {
		return nil
	}
	return app.logger
}

var config_disabled bool = false
var app *application = nil
