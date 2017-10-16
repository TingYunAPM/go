// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func parseJson(data []byte, statusCode int, err error) (error, map[string]interface{}) {

	if err != nil { //http过程有错误
		return err, nil
	}
	if statusCode != 200 {
		return errors.New(fmt.Sprintf("server response status %d", statusCode)), nil
	}
	//拆解返回的json
	jsonData := make(map[string]interface{})
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err, nil
	}
	return nil, jsonData
}

//解析redirect服务器返回值，成功则返回loginhost,失败则返回error
func parseRedirectResult(jsonData map[string]interface{}) (error, string) {
	result, resok := jsonData["result"]
	var resString string = "has no result"
	if resok {
		resString = fmt.Sprint(result)
	}
	if status, ok := jsonData["status"]; !ok { //验证是否有status
		return errors.New("server result have no status."), ""
	} else if v, ok := status.(string); !ok { //类型验证
		return errors.New("server result status format error."), ""
	} else if v != "success" { //值验证
		return errors.New("server result not success: " + resString), ""
	}
	if !resok {
		return errors.New("Redirect server status is success, no result"), ""
	}
	return nil, resString
}

func getRedirectHost(s *serviceDC, protocol string) string {
	host := s.configs.local.CStrings.Read(configLocalStringNbsHost, "redirect.networkbench.com")
	array := strings.Split(host, "://")
	if len(array) > 1 {
		host = array[1]
	}
	array = strings.Split(host, ":")
	if len(array) > 1 {
		host = array[0]
	}
	port := 80
	if protocol != "http" {
		port = 443
	}
	port = int(s.configs.local.CIntegers.Read(configLocalIntegerNbsPort, int64(port)))
	return fmt.Sprintf("%s://%s:%d", protocol, host, port)
}
