// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

//Post请求异步封装
package postRequest

import (
	"io/ioutil"
	"net/http"
)
import "time"
import "github.com/TingYunAPM/go/utils/zip"
import "bytes"
import "sync/atomic"

type Request struct {
	used     int32 //初始化为0,调用callback前原子 +1,不为1则-1后return.调用callback后+1；release前原子+1判定是否为1，否则-1,等为2时return
	callback func(data []byte, statusCode int, err error)
}

func (r *Request) answer(data []byte, statusCode int, err error) {
	use := atomic.AddInt32(&r.used, 1)
	if use != 1 {
		return
	}
	defer atomic.AddInt32(&r.used, 1)
	r.callback(data, statusCode, err)
}

//释放请求对象，不管返回结果
func (r *Request) Release() {
	defer func() { r.callback = nil }()
	if atomic.AddInt32(&r.used, 1) != 1 {
		atomic.AddInt32(&r.used, -1)
		for r.used != 2 {
			time.Sleep(1 * time.Millisecond)
		}
	}
}

//发起一个post请求,返回请求对象
func New(url string, params map[string]string, data []byte, duration time.Duration, callback func(data []byte, statusCode int, err error)) (*Request, error) {
	var body []byte
	var err error = nil
	if v, ok := params["Content-Encoding"]; ok && v == "deflate" {
		body, err = zip.Deflate(data)
	} else {
		body = data
	}
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if nil != err {
		return nil, err
	}
	useParams := make(map[string]string)
	useParams["Accept-Encoding"] = "identity, deflate"
	useParams["Content-Type"] = "Application/json;charset=UTF-8"
	useParams["User-Agent"] = "TingYun-Agent/GoLang"
	for k, v := range params {
		useParams[k] = v
	}
	res := &Request{0, callback}
	for k, v := range useParams {
		request.Header.Add(k, v)
	}
	go func() {
		client := &http.Client{Timeout: duration}
		response, err := client.Do(request)
		if err != nil {
			res.answer(nil, -1, err)
			return
		}
		defer response.Body.Close()
		if response.StatusCode == 200 {
			if b, err := ioutil.ReadAll(response.Body); err != nil { //server返回200，然后读数据失败....
				res.answer(nil, 200, err)
			} else {
				encoding := response.Header.Get("Content-Encoding")
				if encoding == "gzip" || encoding == "deflate" {
					d, err := zip.Inflate(b)
					if err == nil {
						res.answer(d, 200, nil)
						return
					}
				}
				res.answer(b, 200, nil)
			}
		} else {
			res.answer(nil, response.StatusCode, nil)
		}
	}()
	return res, nil
}
