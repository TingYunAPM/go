package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/TingYunAPM/go/utils/httprequest"
)

func TestHttpRequestWithoutData(t *testing.T) {
	url := "http://redirect.networkbench.com/getRedirectHost?licenseKey=fd3a08224ec843552684e8a641af0c2f&version=1.2.0"
	//预期返回结果内容包含"dcs1.networkbench.com"
	c := make(chan int)
	_, err := postRequest.New(url, nil, nil, time.Second*10,
		func(d []byte, statusCode int, e error) {
			if statusCode != 200 {
				t.Error(fmt.Sprint(url, ", status=", statusCode))
			}
			if e != nil {
				t.Error(e)
			}
			s := string(d)
			t.Log(s)
			if (strings.Index(s, ".networkbench.com") == -1) && (strings.Index(s, ".tingyun.com") == -1) {
				t.Error("not find '.networkbench.com' or '.tingyun.com' in " + s)
			}
			c <- 1
		})

	if err != nil {
		t.Error(err)
		c <- 1
	}
	//等待服务器错误或响应后才能结束
	<-c
}

func TestHttpRequestUndeflate(t *testing.T) {
	url := "http://dcs1.networkbench.com/initAgentApp?licenseKey=fd3a08224ec843552684e8a641af0c2f&version=1.2.0"
	params := map[string]string{"Accept-Encoding": "none"}
	data := []byte(`{"host":"centos","appName":["go-request-Undeflate"],"language":"Go","agentVersion":"0.0.1","config":{},"env":{"Go Version":"0.1"}}`)
	timeout := time.Second * 10
	//预期返回结果内容包含"tingyunIdSecret"
	c := make(chan int)
	_, err := postRequest.New(url, params, data, timeout,
		func(d []byte, statusCode int, e error) {
			if statusCode != 200 {
				t.Error(fmt.Sprint(url, ", status=", statusCode))
			}
			if e != nil {
				t.Error(e)
			}
			s := string(d)
			if strings.Index(s, "tingyunIdSecret") == -1 {
				t.Error("not find 'tingyunIdSecret' in " + s)
			}
			c <- 1
		})

	if err != nil {
		t.Error(err)
		c <- 1
	}
	//等待服务器错误或响应后才能结束
	<-c
}

func TestHttpRequestDeflate(t *testing.T) {
	url := "http://dcs1.networkbench.com/initAgentApp?licenseKey=fd3a08224ec843552684e8a641af0c2f&version=1.2.0"
	params := map[string]string{}
	data := []byte(`{"host":"centos","appName":["go-request-deflate"],"language":"Go","agentVersion":"0.0.1","config":{},"env":{"Go Version":"0.1"}}`)
	timeout := time.Second * 10
	//预期返回结果内容包含"tingyunIdSecret"
	c := make(chan int)
	_, err := postRequest.New(url, params, data, timeout,
		func(d []byte, statusCode int, e error) {
			if statusCode != 200 {
				t.Error(fmt.Sprint(url, ", status=", statusCode))
			}
			if e != nil {
				t.Error(e)
			}
			s := string(d)
			if strings.Index(s, "tingyunIdSecret") == -1 {
				t.Error("not find 'tingyunIdSecret' in " + s)
			}
			c <- 1
		})

	if err != nil {
		t.Error(err)
		c <- 1
	}
	//等待服务器错误或响应后才能结束
	<-c
}

func TestHttpRequestHttps(t *testing.T) {
	url := "https://dcs1.networkbench.com/initAgentApp?licenseKey=fd3a08224ec843552684e8a641af0c2f&version=1.2.0"
	data := []byte(`{"host":"centos","appName":["go-request-https"],"language":"Go","agentVersion":"0.0.1","config":{},"env":{"Go Version":"0.1"}}`)
	timeout := time.Second * 10
	//预期返回结果内容包含"tingyunIdSecret"
	c := make(chan int)
	_, err := postRequest.New(url, nil, data, timeout,
		func(d []byte, statusCode int, e error) {
			if statusCode != 200 {
				t.Error(fmt.Sprint(url, ", status=", statusCode))
			}
			if e != nil {
				t.Error(e)
			}
			s := string(d)
			if strings.Index(s, "tingyunIdSecret") == -1 {
				t.Error("not find 'tingyunIdSecret' in " + s)
			}
			c <- 1
		})

	if err != nil {
		t.Error(err)
		c <- 1
	}
	//等待服务器错误或响应后才能结束
	<-c
}

func TestHttpRequestExceptionParse(t *testing.T) {
	url := "http://dcs1.networkbench.com/%gh&%ij"
	params := map[string]string{"Accept-Encoding": "none"}
	data := []byte(`{"host":"centos","appName":["go-request-Undeflate"],"language":"Go","agentVersion":"0.0.1","config":{},"env":{"Go Version":"0.1"}}`)
	timeout := time.Second * 10
	//预期url解析失败
	_, err := postRequest.New(url, params, data, timeout,
		func(d []byte, statusCode int, e error) {
			t.Log(e)
			t.Log(string(statusCode))
			t.Log(string(d))
		})
	if err == nil {
		t.Error("预期url解析失败")
	} else {
		t.Log(err)
	}
}

func TestHttpRequestExceptionPort(t *testing.T) {
	url := "http://dcs1.networkbench.com:2233/initAgentApp?licenseKey=fd3a08224ec843552684e8a641af0c2f&version=1.2.0"
	params := map[string]string{"Accept-Encoding": "none"}
	data := []byte(`{"host":"centos","appName":["go-request-Undeflate"],"language":"Go","agentVersion":"0.0.1","config":{},"env":{"Go Version":"0.1"}}`)
	timeout := time.Second * 10
	//预期无法连接服务器2233端口
	c := make(chan int)
	_, err := postRequest.New(url, params, data, timeout,
		func(d []byte, statusCode int, e error) {
			if e == nil {
				t.Error("预期无法连接服务器2233端口")
			} else {
				t.Log(e)
			}
			c <- 1
		})

	if err != nil {
		t.Error(err)
		c <- 1
	}

	//等待服务器错误或响应后才能结束
	<-c
}

func TestHttpRequestExceptionDomain(t *testing.T) {
	url := "http://xml.networkbench.com/initAgentApp?licenseKey=fd3a08224ec843552684e8a641af0c2f&version=1.2.0"
	params := map[string]string{"Accept-Encoding": "none"}
	data := []byte(`{"host":"centos","appName":["go-request-Undeflate"],"language":"Go","agentVersion":"0.0.1","config":{},"env":{"Go Version":"0.1"}}`)
	timeout := time.Second * 10
	//预期无法解析域名xml.networkbench.com
	c := make(chan int)
	_, err := postRequest.New(url, params, data, timeout,
		func(d []byte, statusCode int, e error) {
			if e == nil {
				t.Error("预期无法解析域名xml.networkbench.com")
			} else {
				t.Log(e)
			}
			c <- 1
		})

	if err != nil {
		t.Error(err)
		c <- 1
	}
	//等待服务器错误或响应后才能结束
	<-c
}

func TestHttpRequestExceptionHost(t *testing.T) {
	url := "http://192.168.8.88/initAgentApp?licenseKey=fd3a08224ec843552684e8a641af0c2f&version=1.2.0"
	params := map[string]string{"Accept-Encoding": "none"}
	data := []byte(`{"host":"centos","appName":["go-request-Undeflate"],"language":"Go","agentVersion":"0.0.1","config":{},"env":{"Go Version":"0.1"}}`)
	timeout := time.Second * 6
	//预期无法与服务器建立连接
	c := make(chan int)
	_, err := postRequest.New(url, params, data, timeout,
		func(d []byte, statusCode int, e error) {
			if e == nil {
				t.Error("预期无法与服务器建立连接")
			} else {
				t.Log(e)
			}
			c <- 1
		})

	if err != nil {
		t.Error(err)
		c <- 1
	}
	//等待服务器错误或响应后才能结束
	<-c
}

func TestHttpRequestException404(t *testing.T) {
	url := "http://dcs1.networkbench.com/initAgen"
	params := map[string]string{"Accept-Encoding": "none"}
	data := []byte(`{"host":"centos","appName":["go-request-Undeflate"],"language":"Go","agentVersion":"0.0.1","config":{},"env":{"Go Version":"0.1"}}`)
	timeout := time.Second * 10
	//预期无法与服务器建立连接
	c := make(chan int)
	_, err := postRequest.New(url, params, data, timeout,
		func(d []byte, statusCode int, e error) {
			if statusCode == 404 {
				t.Log(fmt.Sprint("statusCode=", statusCode))
			} else {
				t.Error(fmt.Sprint("expect 404 but found ", statusCode))
			}
			c <- 1
		})

	if err != nil {
		t.Error(err)
		c <- 1
	}
	//等待服务器错误或响应后才能结束
	<-c
}
