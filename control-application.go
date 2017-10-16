// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/TingYunAPM/go/utils/logger"
)

const (
	LevelOff      = log.LevelOff
	LevelCritical = log.LevelCritical
	LevelError    = log.LevelError
	LevelWarning  = log.LevelWarning
	LevelInfo     = log.LevelInfo
	LevelVerbos   = log.LevelVerbos
	LevelDebug    = log.LevelDebug
	LevelMask     = log.LevelMask
	Audit         = log.Audit
)

//action 结构化
func (a *application) parseAction(action *Action) {
	defer action.destroy()
	a.CreateReportBlock(time.Now())
	a.dataBlock.Append(action)
}

//actionPool里的请求归纳处理
func (a *application) parseActions(parse_max int) int {
	handle_count := 0
	for a.actionPool.Size() > 0 && handle_count < parse_max {
		if action := a.actionPool.Get(); action != nil {
			handle_count++
			if a.configs.HasLogin() {
				a.parseAction(action.(*Action))
			} else {
				//没登陆成功之前，没有任何服务器端配置信息,数据没有处理依据，丢弃。
				action.(*Action).destroy()
			}
		} else {
			break
		}
	}
	return handle_count
}

const (
	serverUnInited     = 0
	serverInLogin      = 1
	serverLoginSuccess = 2
	serverLoginFaild   = 3
	uploadSuccess      = 0
	uploadError        = 1
)

type serverControl struct {
	loginResultTime time.Time //前一次login返回时间(无论成功失败)
	lastUploadTime  time.Time //前一次上传时间，失败重传需要等待几秒钟间隔，据此计算
	loginState      uint8
	uploadState     uint8
	postIsReturn    bool
	reportTime      time.Time
	popcount        int
}

func (s *serverControl) OnReturn() {
	s.postIsReturn = true
	s.lastUploadTime = time.Now()
}

func (s *serverControl) init() {
	s.loginState = serverUnInited
	s.uploadState = uploadSuccess
	s.postIsReturn = false
	s.popcount = 0
}

//redirect, login, 处理
func (a *application) checkLogin() bool {
	//需要添加如下变量

	//appId
	//report计时器

	//处理过程

	//若login状态为serverLoginSuccess,则返回true
	//若login状态为serverUnInited,则开始login过程,置状态为serverInLogin, 返回false
	//若login状态为serverInLogin,则返回false
	//若状态为serverLoginFaild, 5 second waiting,则等待时间=time.Now()-loginResultTime
	//   若等待时间小于5秒,则返回false
	//   若等待时间大于等于5秒,则置login状态为serverUnInited,返回false

	//login完成回调(另一个routine)
	//    1.写日志
	//    2.更新loginResultTime
	//    2.成功,设置application Id, 设置配置项参数,设置login 状态为serverLoginSuccess,置位clear标志
	//    3.失败,设置login状态为 5 second waiting

	switch a.serverCtrl.loginState {
	case serverUnInited:
		a.startLogin()
		return false
	case serverInLogin:
		return false
	case serverLoginFaild:
		a.checkFaildState()
		return false
	case serverLoginSuccess:
		return true
	default:
		//未定义行为, 不该出现 log error
		return false
	}
}

func (a *application) ReadApdex(name string) int32 {
	return int32(a.configs.apdexs.Read(name))
}

func (a *application) setNextReportTime(value time.Time, reportInterval int) {
	unix := value.Unix()
	unix = int64(reportInterval) + unix - unix%int64(reportInterval)
	a.serverCtrl.reportTime = time.Unix(unix, 0)
}

//login返回, 验证application Id, session key
func (a *application) parseLoginResult(jsonData map[string]interface{}) error {
	//验证status
	//如果status不为success,则从result生成错误信息,返回error
	for {
		if err, status := jsonReadString(jsonData, "status"); err != nil {
			break
		} else if status != "success" {
			break
		} else if err, result := jsonReadObjects(jsonData, "result"); err != nil {
			break
		} else if !a.configs.UpdateServerConfig(result) {
			break
		}
		enabled := a.configs.server.CBools.Read(configServerBoolEnabled, true)
		if enabled {
			enabled = a.configs.server_ext.CBools.Read(configServerConfigBoolAgentEnabled, true)
		}
		if enabled {
			now := time.Now()
			if a.configs.NeverLogin() {
				a.logger.Println(LevelInfo, "ApplicationStart...")
				a.setNextReportTime(now, int(a.configs.server.CIntegers.Read(configServerIntegerDataSentInterval, 60)))
				a.CreateReportBlock(now)
			}
		}

		return nil
	}
	b, _ := json.Marshal(jsonData)
	return errors.New("server result json error : " + string(b))
}
func (a *application) startLogin() {
	a.serverCtrl.loginState = serverInLogin
	err := a.server.Login(func(err error, result map[string]interface{}) {
		defer a.serverCtrl.OnReturn()
		a.serverCtrl.loginResultTime = time.Now()
		if err != nil { //login过程有错误,置状态,写日志
			a.logger.Println(log.LevelError, err)
			a.serverCtrl.loginState = serverLoginFaild
		} else if err := a.parseLoginResult(result); err != nil { //login结果有错误,写日志,置状态
			a.logger.Println(log.LevelError, err)
			a.serverCtrl.loginState = serverLoginFaild
		} else { //login成功
			if !a.configs.HasLogin() {
				a.serverCtrl.lastUploadTime = time.Now()
			}
			a.serverCtrl.loginState = serverLoginSuccess
		}
	})
	if err != nil {
		a.serverCtrl.loginResultTime = time.Now()
		a.serverCtrl.loginState = serverLoginFaild
	}
}
func (a *application) checkFaildState() {
	if now := time.Now(); now.Sub(a.serverCtrl.loginResultTime) < 5*time.Second {
		return
	}
	a.serverCtrl.loginState = serverUnInited
}

func (a *application) timerCheck() {
	if a.serverCtrl.postIsReturn {
		a.server.ReleaseRequest()
		a.serverCtrl.postIsReturn = false
	}
	//不在LoginSuccess状态,返回
	if !a.checkLogin() {
		return
	}
	a.configs.UpdateConfig(a.Runtime.Init)
	a.upload()
}
func (a *application) popReported() {
	for a.serverCtrl.popcount > 0 {
		a.serverCtrl.popcount--
		data, _ := a.reportQueue.PopFront()
		data.(*structAppData).destroy()
	}
}
func (a *application) upload() {
	if a.server.request != nil {
		return
	}
	a.popReported()
	if a.serverCtrl.uploadState != uploadSuccess {
		if time.Now().Sub(a.serverCtrl.lastUploadTime) < 5*time.Second {
			return
		}
	}
	if a.reportQueue.Size() > 0 {
		//fmt.Printf("goroutine:%d\n", runtime.NumGoroutine())
		a.mergeReport()
		data, _ := a.reportQueue.Front().Value()
		b, _ := data.(*structAppData).Serialize()
		err := a.server.Upload(b, func(err error, rcode int, statusCode int) {
			defer a.serverCtrl.OnReturn()
			if rcode == -2 {
				//需要merge
				a.serverCtrl.uploadState = uploadError
				//fmt.Println("Upload error :", statusCode, ",", err)
				return
			}
			//返回http 200 ,正常json,说明上传没问题,缓存在队列的数据可以pop出去了
			a.serverCtrl.popcount++
			if err == nil {
				//fmt.Println("Upload success.")
			} else {
				//fmt.Println("Upload status: ", err)
				//上传失败,重新启动login过程
				a.serverCtrl.loginState = serverLoginFaild
			}
		})
		if err != nil {
			a.logger.Println(LevelError, "App.", "Upload Error :", err)
			a.serverCtrl.uploadState = uploadError
			a.serverCtrl.lastUploadTime = time.Now()
		}
	}
}

//1.数据按发送间隔拆分对象
//   取当前时间,计算与lastCommitTime计时器的差值,若小于发送间隔,则return
//   大于等于发送间隔,取下数据对象push_back到发送队列.
//2.从发送队列取数据发送，若发送队列容量大于10，则合并front数据。

func (a *application) CommitData() {
	if !a.configs.HasLogin() {
		return
	}
	now := time.Now()
	a.CreateReportBlock(now)
	//结构数据提交到发送队列，超出队列长度的pop_front, 合并到front
	if now.After(a.serverCtrl.reportTime) {
		interval := a.configs.server.CIntegers.Read(configServerIntegerDataSentInterval, 60)
		a.setNextReportTime(now, int(interval))
		a.dataBlock.end(&a.Runtime)
		a.dataBlock.endTime = now
		a.logger.Println(LevelInfo, "Report Commit:", a.dataBlock.startTime.Unix(), "-", a.dataBlock.endTime.Unix())
		a.reportQueue.PushBack(a.dataBlock)
		a.dataBlock = nil
		a.CreateReportBlock(now)
	}
}
func (a *application) mergeReport() {
	saveCount := a.configs.local.CIntegers.Read(configLocalIntegerNbsSaveCount, 10)
	if saveCount < 1 {
		saveCount = 1
	}
	for a.reportQueue.Size() > int(saveCount) {
		data, _ := a.reportQueue.PopFront()
		t := data.(*structAppData)
		data, _ = a.reportQueue.Front().Value()
		data.(*structAppData).Merge(t)
		t.destroy()
	}
}
func (a *application) loop(running func() bool) {
	last_count := 1
	for running() {
		//1.action事件处理
		handle_count := a.parseActions(10000)
		//2.提交数据
		a.CommitData()
		//3.server通信事件处理
		a.timerCheck()
		if handle_count == 0 {
			if last_count == 0 {
				time.Sleep(10 * time.Millisecond)
			} else {
				time.Sleep(1 * time.Millisecond)
			}
		}
		last_count = handle_count
	}
}

func makeLoginRequest() ([]byte, error) {
	//用参数初始化头信息
	//构造login请求包
	//序列化数据
	get_host := func() string {
		if host, e := os.Hostname(); e == nil {
			return host
		}
		return "Unknown"
	}
	get_path := func() string {
		file, _ := exec.LookPath(os.Args[0])
		path, _ := filepath.Abs(file)
		return path
	}

	return json.Marshal(map[string]interface{}{
		"host":         get_host(),
		"appName":      [1]string{app.configs.local.CStrings.Read(configLocalStringNbsAppName, "TingYunDefault")},
		"language":     "Go",
		"agentVersion": TINGYUN_GO_AGENT_VERSION,
		"pid":          os.Getpid(),
		"config": map[string]interface{}{
			"nbs.app_name":      app.configs.local.CStrings.Read(configLocalStringNbsAppName, "GO_LANG"),
			"nbs.license_key":   app.configs.local.CStrings.Read(configLocalStringNbsLicenseKey, ""),
			"nbs.log_file_name": app.configs.local.CStrings.Read(configLocalStringNbsLogFileName, "agent.log"),
			"nbs.audit":         app.configs.local.CBools.Read(configLocalBoolAudit, false),
			"nbs.max_log_count": app.configs.local.CIntegers.Read(configLocalIntegerNbsMaxLogCount, 3),
			"nbs.max_log_size":  app.configs.local.CIntegers.Read(configLocalIntegerNbsMaxLogSize, 10),
			"nbs.ssl":           app.configs.local.CBools.Read(configLocalBoolSSL, true),
		},
		"env": map[string]string{
			"cmdline":    get_path(),
			"OS":         runtime.GOOS,
			"ARCH":       runtime.GOARCH,
			"Compiler":   runtime.Compiler,
			"Go-Version": runtime.Version(),
			"GOROOT":     runtime.GOROOT(),
		},
	})
}
