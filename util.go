// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func round(x float64) int {
	if x < 0.0 {
		return -round(-x)
	}
	r := int(x)
	if x-float64(r) < 0.5 {
		return r
	}
	return r + 1
}
func binarySearch(p []float64, value float64) int {
	Begin, Len := 0, len(p)
	for Len > 0 {
		middle := Len / 2
		if value < p[Begin+middle] {
			Len = middle
		} else if p[Begin+middle] < value {
			Begin, Len = Begin+middle+1, Len-middle-1
		} else {
			return Begin + middle
		}
	}
	return -1 - Begin
}

func callStack(skip int) interface{} {
	var slice []string
	slice = make([]string, 0, 15)
	opc := uintptr(0)
	for i := skip + 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		/*
			过滤包含tingyun字样的堆栈
			合并可能的连续相同的地址：增加wrapper并过滤tingyun函数导致 或 递归调用 导致
			net/http.HandlerFunc.ServeHTTP  -> tingyun.wrapper -> net/http.HandlerFunc.ServeHTTP
			f -> f
		*/
		if opc == pc {
			continue
		}
		fname := getnameByAddr(pc)
		index := strings.Index(fname, "/tingyun/")
		if index > 0 {
			continue
		}
		opc = pc
		//截断源文件名
		index = strings.Index(file, "/src/")
		if index > 0 {
			file = file[index+5 : len(file)]
		}
		//fmt.Println(fmt.Sprintf("%s(%s:%d)", fname, file, line))
		slice = append(slice, fmt.Sprintf("%s(%s:%d)", fname, file, line))
	}
	return slice
}
func getnameByAddr(p interface{}) string {
	ptr, _ := strconv.ParseInt(fmt.Sprintf("%x", p), 16, 64)
	return runtime.FuncForPC(uintptr(ptr)).Name()
}
func unicId(t time.Time, p interface{}) string {
	return strings.Replace(fmt.Sprintf("%x-%p", t.UnixNano(), p), "0x", "", -1)
}

func md5sum(src string) string {
	val := md5.New()
	val.Write([]byte(src))
	return fmt.Sprintf("%x", val.Sum(nil))
}

type timeRange struct {
	begin    time.Time
	duration time.Duration
}

func (t *timeRange) End() {

	t.duration = time.Now().Sub(t.begin)
}
func (t *timeRange) Init() {
	t.begin = time.Now()
	t.duration = -1
}
func (r *timeRange) EndTime() time.Time {
	ret := r.begin
	return ret.Add(r.duration)
}
func (r *timeRange) Inside(t *timeRange) bool {
	if r.begin.Before(t.begin) || t.duration < r.duration {
		return false
	}
	return !r.EndTime().After(t.EndTime())
}

func jsonReadString(jsonData map[string]interface{}, name string) (error, string) {
	if r, ok := jsonData[name]; !ok { //验证是否有name
		return errors.New("Has no " + name), ""
	} else if v, ok := r.(string); !ok { //类型验证
		return errors.New("json \"" + name + "\" not string."), ""
	} else {
		return nil, v
	}
}
func jsonReadObjects(jsonData map[string]interface{}, name string) (error, map[string]interface{}) {
	if r, ok := jsonData[name]; !ok { //验证是否有name
		return errors.New("Has no " + name), nil
	} else if v, ok := r.(map[string]interface{}); !ok { //类型验证
		return errors.New("json \"" + name + "\" not objects."), nil
	} else {
		return nil, v
	}
}
func jsonReadBool(jsonData map[string]interface{}, name string) (error, bool) {
	if r, ok := jsonData[name]; !ok { //验证是否有name
		return errors.New("Has no " + name), false
	} else if v, ok := r.(bool); !ok { //类型验证
		return errors.New("json \"" + name + "\" not bool."), false
	} else {
		return nil, v
	}
}
func readInt(v interface{}) (error, int) {
	switch r := v.(type) {
	case float64:
		return nil, int(r)
	case float32:
		return nil, int(r)
	case int:
		return nil, r
	case int32:
		return nil, int(r)
	case int64:
		return nil, int(r)
	case uint32:
		return nil, int(r)
	case uint64:
		return nil, int(r)
	default:
		return errors.New(fmt.Sprint(v, ":  not int value.")), 0
	}
}
func readFloat(v interface{}) (error, float64) {
	switch r := v.(type) {
	case float64:
		return nil, r
	case float32:
		return nil, float64(r)
	case int:
		return nil, float64(r)
	case int32:
		return nil, float64(r)
	case int64:
		return nil, float64(r)
	case uint32:
		return nil, float64(r)
	case uint64:
		return nil, float64(r)
	default:
		return errors.New(fmt.Sprint(v, ":  not float value.")), 0.0
	}
}
func jsonReadInt(jsonData map[string]interface{}, name string) (error, int) {
	if t, ok := jsonData[name]; !ok { //验证是否有name
		return errors.New("Has no " + name), 0
	} else {
		return readInt(t)
	}
}
func jsonReadFloat(jsonData map[string]interface{}, name string) (error, float64) {
	if t, ok := jsonData[name]; !ok { //验证是否有name
		return errors.New("Has no " + name), 0.0
	} else {
		switch r := t.(type) {
		case float64:
			return nil, t.(float64)
		case float32:
			return nil, float64(r)
		case int:
			return nil, float64(r)
		case int32:
			return nil, float64(r)
		case int64:
			return nil, float64(r)
		case uint32:
			return nil, float64(r)
		case uint64:
			return nil, float64(r)
		default:
			return errors.New(fmt.Sprint(name+": ", t, " not float value.")), 0.0
		}
	}
}
func jsonToString(jsonData map[string]interface{}, name string) (error, string) {
	if r, ok := jsonData[name]; !ok { //验证是否有name
		return errors.New("Has no " + name), ""
	} else if v, ok := r.(string); ok { //类型验证
		return nil, v
	} else {
		switch t := r.(type) {
		case float64:
			return nil, fmt.Sprintf("%d", int64(t))
		case float32:
			return nil, fmt.Sprintf("%d", int64(t))
		case int:
			return nil, fmt.Sprintf("%d", int64(t))
		case int32:
			return nil, fmt.Sprintf("%d", int64(t))
		case int64:
			return nil, fmt.Sprintf("%d", t)
		case uint32:
			return nil, fmt.Sprintf("%d", int64(t))
		case uint64:
			return nil, fmt.Sprintf("%d", t)
		default:
			return errors.New(fmt.Sprint(name+": ", t, " not string or int value.")), ""
		}
	}
}
func parseMethod(method string) (string, string) {
	array := strings.Split(method, "::")
	arrayLen := len(array)
	if arrayLen > 1 {
		classRet := array[0]
		for i := 1; i < arrayLen-1; i++ {
			classRet = classRet + "::" + array[i]
		}
		return classRet, array[arrayLen-1]
	}
	array = strings.Split(method, ".")
	arrayLen = len(array)
	if arrayLen == 1 {
		return "", method
	}
	classRet := array[0]
	for i := 1; i < arrayLen-1; i++ {
		classRet = classRet + "." + array[i]
	}
	return classRet, array[arrayLen-1]
}
func jsonDecodeArray(src string) []interface{} {
	ret := make([]interface{}, 0)
	err := json.Unmarshal([]byte(src), &ret)
	if err != nil {
		return nil
	}
	return ret
}
