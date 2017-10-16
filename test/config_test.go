package main

import (
	"testing"
	"time"

	"github.com/TingYunAPM/go"
)

type application struct {
	_config config "配置选项集合"
}
type configInfo map[string]interface{}
type config struct {
	array [2]configInfo
	index int
}

func (c *config) value(name string) (interface{}, bool) {
	v, ok := c.array[c.index][name]
	return v, ok
}
func (c *config) update(name string, value interface{}) {
	c.array[c.index^1][name] = value
}
func (c *config) commit() {
	for k, v := range c.array[c.index] {
		if _, exist := c.array[c.index^1][k]; !exist {
			c.array[c.index^1][k] = v
		}
	}
	c.index = c.index ^ 1
}
func createConfig(configfile string) *config {
	res := &config{}
	res.array[0] = make(configInfo)
	res.array[1] = make(configInfo)
	res.index = 0
	return res
}

func TestConfig(t *testing.T) {
	config := createConfig("ls")
	config.update("a", 1)
	config.update("b", "test")
	config.commit()
	{
		_, ok := config.value("a")
		if ok {
			//t.Log("a : %d\n", x)
		} else {
			t.Error("a not exist\n")
		}
		_, ok = config.value("b")
		if ok {
			//t.Log("b : %s\n", y)
		} else {
			t.Error("b not exist\n")
		}

	}
	config.update("c", "testc")
	config.commit()
	{
		_, ok := config.value("a")
		if ok {
			//t.Log("a : %d\n", x)
		} else {
			t.Error("a not exist\n")
		}
		_, ok = config.value("b")
		if ok {
			//t.Log("b : %s\n", y)
		} else {
			t.Error("b not exist\n")
		}
		_, ok = config.value("c")
		if ok {
			//t.Log("c : %s\n", z)
		} else {
			t.Error("c not exist\n")
		}

	}
}

func TestAppConfig(t *testing.T) {

	err := tingyun.AppInit("agent-disabled.json")
	if err != nil {
		//探针本地禁用，预期返回err
		t.Log(err)
	} else {
		t.Error("探针未被禁用")
	}
	defer tingyun.AppStop()

	c := make(chan int)

	go func() {
		action, err := tingyun.CreateAction("user", "login")
		//探针本地禁用，预期返回err
		if err != nil {
			t.Log(err)
		} else {
			t.Error("探针未被禁用")
		}
		defer action.Finish()
		time.Sleep(3 * time.Second)
		c <- 1
	}()

	go func() {
		action, _ := tingyun.CreateAction("user", "logout")
		defer action.Finish()
		time.Sleep(2 * time.Second)
		c <- 1
	}()

	<-c
	<-c

}

func TestInvalidLicense(t *testing.T) {

	tingyun.AppInit("invalid-license.json")
	defer tingyun.AppStop()

	time.Sleep(10 * time.Second)

	c := make(chan int)

	go func() {
		action, _ := tingyun.CreateAction("user", "login")
		defer action.Finish()
		time.Sleep(3 * time.Second)
		c <- 1
	}()

	go func() {
		action, _ := tingyun.CreateAction("user", "logout")
		defer action.Finish()
		time.Sleep(2 * time.Second)
		c <- 1
	}()

	<-c
	<-c

}
