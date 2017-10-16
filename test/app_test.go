package main

import (
	"testing"
	"time"

	"github.com/TingYunAPM/go"
)

func TestApp(t *testing.T) {
	err := tingyun.AppInit("app.json")
	if err != nil {
		t.Error(err)
	}
	defer tingyun.AppStop()

	c := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		c <- 1
	}()
	<-c

}

func TestMissConfig(t *testing.T) {
	err := tingyun.AppInit("miss.json")
	if err != nil {
		t.Log(err)
	} else {
		t.Error("没有配置文件时，创建app应该失败")
	}
	defer tingyun.AppStop()

	c := make(chan int)
	go func() {
		action, err := tingyun.CreateAction("user", "login")
		if err != nil {
			t.Log(err)
		} else {
			t.Error("没有配置文件时，创建action应该失败")
		}
		defer action.Finish()
		time.Sleep(3 * time.Second)
		c <- 1
	}()
	<-c

}

func TestMissConfig2(t *testing.T) {
	tingyun.AppInit("miss.json")
	defer tingyun.AppStop()

	c := make(chan int)
	go func() {
		action, _ := tingyun.CreateAction("user", "login")
		defer action.Finish()
		time.Sleep(3 * time.Second)
		c <- 1
	}()
	<-c

}
