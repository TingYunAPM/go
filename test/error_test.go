package main

import (
	"errors"
	"testing"

	"github.com/TingYunAPM/go"
)

func TestErrorStatus(t *testing.T) {
	tingyun.AppInit("app.json")
	defer tingyun.AppStop()

	action, _ := tingyun.CreateAction("user", "logerror")
	defer action.Finish()

	action.SetStatusCode(502)
	if action.Slow() || action.HasError() {
		action.AddCustomParam("status_code", "502")
	}

}

func TestErrorClass(t *testing.T) {
	tingyun.AppInit("app.json")
	defer tingyun.AppStop()

	action, _ := tingyun.CreateAction("user", "logerror")
	defer action.Finish()

	action.SetError(errors.New("user invalid"))
	if action.Slow() || action.HasError() {
		action.AddCustomParam("class", "TestErrorClass")
	}

}
