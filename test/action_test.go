package main

import (
	"testing"
	"time"

	"github.com/TingYunAPM/go"
)

func TestActionIgnore(t *testing.T) {
	tingyun.AppInit("app.json")
	defer tingyun.AppStop()

	c := make(chan int)
	go func() {
		action, _ := tingyun.CreateAction("route", "login")
		defer action.Finish()

		func() {
			defer action.CreateComponent("A").Finish()
			func() {
				defer action.CreateComponent("A,C").Finish()
				time.Sleep(10 * time.Millisecond)

			}()
			time.Sleep(15 * time.Millisecond)
		}()

		c <- 1
	}()
	<-c

}
