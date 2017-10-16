// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"errors"

	"github.com/TingYunAPM/go/utils/list"
	"github.com/TingYunAPM/go/utils/logger"
	"github.com/TingYunAPM/go/utils/pool"
	"github.com/TingYunAPM/go/utils/service"
)

func (a *application) stop() {
	if a == nil {
		return
	}
	a.svc.Stop()
	a.logger.Printf(LevelInfo, "Agent stoped\n")
	a.logger.Release()
	a.configs.Release()
}

//action结束,将Action事务对象抛给app处理
func (a *application) appendAction(action *Action) {
	action.time.End()
	a.actionPool.Put(action)
}

//给Action.Finish调用
func append_action(action *Action) {
	if app != nil {
		app.appendAction(action)
	} else {
		//释放Action对象
		action.destroy()
	}
}
func readServerConfigInt(id int, default_value int) int {
	if app == nil {
		return default_value
	}
	return int(app.configs.server_ext.CIntegers.Read(id, int64(default_value)))
}
func readServerConfigBool(id int, default_value bool) bool {
	if app == nil {
		return default_value
	}
	return app.configs.server_ext.CBools.Read(id, default_value)
}

func readServerConfigString(id int, defaultValue string) string {
	if app == nil {
		return defaultValue
	}
	return app.configs.server_ext.CStrings.Read(id, defaultValue)
}
func (a *application) init(configfile string) (*application, error) {
	err := a.configs.Init(configfile)
	if err != nil {
		return nil, err
	}
	if enabled := a.configs.local.CBools.Read(configLocalBoolAgentEnable, true); !enabled {
		config_disabled = true
		a.configs.Release()
		return nil, errors.New("Agent Is disabled by config file.")
	}
	if appname := a.configs.local.CStrings.Read(configLocalStringNbsAppName, ""); appname == "" {
		return nil, errors.New(configfile + ": nbs.app_name not found.")
	}

	if license := a.configs.local.CStrings.Read(configLocalStringNbsLicenseKey, ""); license == "" {
		return nil, errors.New(configfile + ": nbs.license_key not found.")
	}
	a.logger = log.New(&a.configs.local)
	a.actionPool.Init()
	a.server.init(&a.configs)
	a.serverCtrl.init()
	a.reportQueue.Init()
	a.dataBlock = nil
	a.Runtime.Init()
	a.svc.Start(a.loop)
	return a, nil
}
func (a *application) createAction(instance string, method string) (*Action, error) {
	if enabled := a.configs.server.CBools.Read(configServerBoolEnabled, true); !enabled {
		return nil, errors.New("Agent disabled by server config.")
	}
	ret := (&Action{}).init(formatActionName(instance, method))
	ret.method = method
	return ret, nil
}

type application struct {
	configs     configurations      "配置选项集合"
	actionPool  pool.SerialReadPool "完成事务消息池"
	logger      *log.Logger
	svc         service.Service
	server      serviceDC
	serverCtrl  serverControl
	reportQueue list.List
	dataBlock   *structAppData
	Runtime     runtimePerf
	//	actionTemps    *pool.Pool
	//	componentTemps *pool.Pool
}
