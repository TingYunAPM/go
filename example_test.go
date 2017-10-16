// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"time"

	"github.com/TingYunAPM/go"
)

func ExampleAppInit() {
	tingyun.AppInit("tingyun.json")
	tingyun.AppStop()
}
func ExampleCreateAction() {
	tingyun.AppInit("tingyun.json")
	defer tingyun.AppStop()
	action, _ := tingyun.CreateAction("/login", "main.ExampleCreateAction")
	time.Sleep(time.Millisecond * 100)
	action.Finish()
}
func ExampleAction_CreateComponent() {
	tingyun.AppInit("tingyun.json")
	defer tingyun.AppStop()
	action, _ := tingyun.CreateAction("/login", "main.ExampleCreateAction")
	component := action.CreateComponent("normalExample")
	subComponent := component.CreateComponent("subcomponent")
	time.Sleep(time.Millisecond * 100)
	subComponent.Finish()
	component.Finish()
	action.Finish()
}

func ExampleAction_CreateDBComponent() {
	tingyun.AppInit("tingyun.json")
	defer tingyun.AppStop()
	action, _ := tingyun.CreateAction("/login", "main.ExampleCreateAction")
	dbComponent := action.CreateDBComponent(tingyun.ComponentMysql, "", "mydatabase", "mytable", "select", "ExampleCreateDBComponent")
	time.Sleep(time.Millisecond * 100)
	dbComponent.Finish()
	action.Finish()
}
func ExampleAction_CreateExternalComponent() {
	tingyun.AppInit("tingyun.json")
	defer tingyun.AppStop()
	action, _ := tingyun.CreateAction("/login", "main.ExampleCreateAction")
	externalComponent := action.CreateExternalComponent("/external", "main.ExampleCreateExternalComponent")
	time.Sleep(time.Millisecond * 100)
	externalComponent.Finish()
	action.Finish()
}
