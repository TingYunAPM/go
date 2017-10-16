package main

import (
	"time"
	//"time"
	//	"fmt"
	"os"
	"testing"

	"github.com/TingYunAPM/go/utils/logger"
)

import "io/ioutil"
import "encoding/json"
import "github.com/TingYunAPM/go/utils/config"

func initConfigFile(configfile string) (*configBase.ConfigBase, error) {
	res := configBase.CreateConfig()
	bytes, err := ioutil.ReadFile(configfile)
	if err != nil {
		return nil, err
	}
	jsonData := map[string]interface{}{}
	if err := json.Unmarshal(bytes, &jsonData); err != nil {
		return nil, err
	}
	for k, v := range jsonData {
		res.Update(k, v)
	}
	res.Commit()
	return res, nil
}

func FileSize(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}

func FileRemove(file string) error {
	return os.Remove(file)
}

func TestSmallLog(t *testing.T) {
	conf, err := initConfigFile("log-small.json")
	if err != nil {
		t.Error("json fail")
		return
	}

	c := make(chan int)
	l := log.New(conf)

	r := func() {
		for i := 0; i < 1000000; i++ {
			if i%1000 == 0 {
				time.Sleep(time.Millisecond)
			}
			l.Printf(log.LevelInfo|log.Audit, "loop count %d\n", i)
		}
		c <- 1
	}
	go r()
	go r()
	<-c
	<-c
	l.Release()
	//
	size, err := FileSize("log/log-small.log")
	if err != nil {
		t.Error(err)
	}
	if size > 1024*1024 {
		t.Error("log size > %d", size)
	}

}
