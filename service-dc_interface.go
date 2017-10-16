// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

//与server通信,login, upload

import (
	"errors"
	"fmt"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/TingYunAPM/go/utils/httprequest"
)

//登陆中..1.inredirect
//       2.inInit
type serviceDC struct {
	locked     int32 //login由redirect状态切换到init状态时提供保护，防止这时候Release
	configs    *configurations
	request    *postRequest.Request
	uploadHost string
	inLogin    bool
}

func (s *serviceDC) ReleaseRequest() {
	if s.request != nil {
		s.request.Release()
		s.request = nil
	}
}

//Login --启动登陆过程,如果已经在login中,返回error
func (s *serviceDC) Login(callback func(error, map[string]interface{})) error {
	if s.inLogin {
		return errors.New("Login already Startd.")
	}
	protocol := "https"
	if !s.configs.local.CBools.Read(configLocalBoolSSL, true) {
		protocol = "http"
	}
	appName := s.configs.local.CStrings.Read(configLocalStringNbsAppName, "GO_LANG")
	license := s.configs.local.CStrings.Read(configLocalStringNbsLicenseKey, "_")
	requrl := fmt.Sprintf("%s/getRedirectHost?app=%s&licenseKey=%s&request=entry&version=%s", getRedirectHost(s, protocol), url.QueryEscape(appName), license, "1.2.0")
	params := make(map[string]string)
	var err error = nil
	Log().Println(LevelInfo, "Redirect:", requrl)
	//post数据到redirect服务器
	s.request, err = postRequest.New(requrl, params, []byte("{}"), time.Second*10, func(data []byte, statusCode int, err error) {

		//完成回调,在另一个routine中触发
		use := atomic.AddInt32(&s.locked, 1)
		defer atomic.AddInt32(&s.locked, -1)
		if use != 1 {
			return
		}
		if err == nil {
			Log().Println(LevelInfo, "Redirect Status Code:", statusCode)
			if len(data) > 0 {
				Log().Println(LevelInfo|Audit, "Redirect Response Data:", string(data))
			}
		}
		var e error = nil
		for {
			var jsonData map[string]interface{}
			if e, jsonData = parseJson(data, statusCode, err); e != nil {
				break
			}
			if e, s.uploadHost = parseRedirectResult(jsonData); e != nil {
				break
			}
			var b []byte
			if b, e = makeLoginRequest(); e != nil {
				break
			}
			requrl := fmt.Sprintf("%s://%s/initAgentApp?app=%s&licenseKey=%s&request=login&version=%s", protocol, s.uploadHost, url.QueryEscape(appName), url.QueryEscape(license), "1.2.0")
			Log().Println(LevelInfo, "Login:", requrl)
			Log().Println(LevelInfo|Audit, "Login Request Data:", string(b))
			//post数据到login服务器,上行数据启用deflate压缩
			s.request, e = postRequest.New(requrl, map[string]string{"Content-Encoding": "deflate"}, b, time.Second*10, func(data []byte, statusCode int, err error) {
				use := atomic.AddInt32(&s.locked, 1)
				defer atomic.AddInt32(&s.locked, -1)
				if use != 1 {
					return
				}
				s.inLogin = false
				if err == nil {
					Log().Println(LevelInfo, "Login Status Code:", statusCode)
					if len(data) > 0 {
						Log().Println(LevelInfo|Audit, "Login Response Data:", string(data))
					}
				}
				er, r := parseJson(data, statusCode, err)
				callback(er, r)
			})
			break
		}
		if e != nil {
			s.inLogin = false
			callback(e, nil)
		}
	})
	if err == nil {
		s.inLogin = true
	}
	return err
}

//上传数据,如果inLogin, 返回false,否则创建request,
func (s *serviceDC) Upload(data []byte, callback func(err error, rCode int, httpStatus int)) error {
	if s.inLogin {
		return errors.New("server in login")
	}
	sessionKey := s.configs.server.CStrings.Read(configServerStringAppSessionKey, "")
	if sessionKey == "" {
		return errors.New("server session key not found")
	}
	protocol := "https"
	if !s.configs.local.CBools.Read(configLocalBoolSSL, true) {
		protocol = "http"
	}
	appName := s.configs.local.CStrings.Read(configLocalStringNbsAppName, "GO_LANG")
	license := s.configs.local.CStrings.Read(configLocalStringNbsLicenseKey, "_")

	requrl := fmt.Sprintf("%s://%s/upload?app=%s&licenseKey=%s&appSessionKey=%s&version=%s", protocol, s.uploadHost, url.QueryEscape(appName), url.QueryEscape(license), url.QueryEscape(sessionKey), "1.2.0")
	Log().Println(LevelInfo, "Upload", len(data), "bytes:", requrl)
	Log().Println(LevelInfo|Audit, "Upload Request Data:", string(data))
	//post数据到login服务器,上行数据启用deflate压缩
	var e error = nil
	s.request, e = postRequest.New(requrl, map[string]string{"Content-Encoding": "deflate"}, data, time.Second*10, func(data []byte, statusCode int, err error) {
		use := atomic.AddInt32(&s.locked, 1)
		defer atomic.AddInt32(&s.locked, -1)
		if use != 1 {
			return
		}
		if err == nil {
			Log().Println(LevelInfo, "Upload Status Code:", statusCode)
			if len(data) > 0 {
				Log().Println(LevelInfo|Audit, "Upload Response Data:", string(data))
			}
		}
		er, r := parseJson(data, statusCode, err)
		if er != nil {
			//发生网络错误（即http发送失败）
			//DC故障（即可以获取到http响应，但状态码不等于200或返回内容为非法json）
			Log().Println(LevelError, "Upload Error:", er, r)
			callback(er, -2, statusCode)
		} else if er, status := jsonReadString(r, "status"); er == nil && status == "success" {
			callback(nil, -1, statusCode)
		} else {
			Log().Println(LevelError, "Upload Result:", string(data))
			callback(errors.New(string(data)), -1, statusCode)
		}
	})
	return e
}
func (s *serviceDC) Release() {
	for { //serverObject对象的生命周期与app的生命周期同样长，
		//这里的等待只有在app停止时的瞬间，post恰好返回时,才可能会发生。所以这里的等待不会成为性能瓶颈
		use := atomic.AddInt32(&s.locked, 1)
		if use == 1 {
			break
		}
		atomic.AddInt32(&s.locked, -1)
		time.Sleep(1 * time.Millisecond)
	}
	s.ReleaseRequest()
	s.configs = nil
}
func (s *serviceDC) init(config *configurations) {
	s.configs = config
	//fmt.Println("serverObject.init")
	s.request = nil
	s.inLogin = false
	s.locked = 0
}
