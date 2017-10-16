package main

import (
	"testing"
	"time"

	"github.com/TingYunAPM/go"
)

func TestTrace(t *testing.T) {

	//
	tingyun.AppInit("app.json")
	defer tingyun.AppStop()
	//

	c := make(chan int)

	go func() {
		//
		action, _ := tingyun.CreateAction("user", "login")
		defer action.Finish()
		//

		func() {
			defer action.CreateComponent("A:B").Finish()

			func() {
				defer action.CreateComponent("A,C").Finish()

				time.Sleep(10 * time.Millisecond)
				func() {
					defer action.CreateComponent("\"X'").Finish()

					time.Sleep(5 * time.Second)
				}()
			}()
			time.Sleep(15 * time.Millisecond)
		}()
		time.Sleep(20 * time.Millisecond)

		c <- 1
	}()

	<-c

}
