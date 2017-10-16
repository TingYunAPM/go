package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/TingYunAPM/go"
)

func childComponent(t *testing.T, level int, name string, action *tingyun.Action) {
	defer action.CreateComponent(name).Finish()

	if level == 0 {
		time.Sleep(time.Millisecond)
		return
	} else {
		childComponent(t, level-1, fmt.Sprintf("%s-%d", name, level), action)
		childComponent(t, level-1, fmt.Sprintf("%s+%d", name, level), action)
	}
}

func TestParentComponent(t *testing.T) {

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
		childComponent(t, 3, "parse", action)
		childComponent(t, 10, "render", action)
		c <- 1
	}()

	<-c

}

func TestComponent(t *testing.T) {

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
			defer action.CreateComponent("root").Finish()

			func() {
				defer action.CreateComponent("child").Finish()

				func() {
					defer action.CreateComponent("leaf").Finish()
					time.Sleep(5 * time.Second)
				}()

				func() {
					defer action.CreateComponent("short").Finish()
				}()
			}()
			time.Sleep(15 * time.Millisecond)
		}()
		time.Sleep(20 * time.Millisecond)

		c <- 1
	}()

	<-c

}

func TestDisableComponent(t *testing.T) {

	//
	tingyun.AppInit("miss.json")
	defer tingyun.AppStop()
	//

	c := make(chan int)

	go func() {
		//
		action, _ := tingyun.CreateAction("user", "login")
		defer action.Finish()
		//

		func() {
			defer action.CreateComponent("root").Finish()

			func() {
				defer action.CreateComponent("child").Finish()

				time.Sleep(10 * time.Millisecond)
				func() {
					defer action.CreateComponent("leaf").Finish()

					time.Sleep(5 * time.Millisecond)
				}()
			}()
			time.Sleep(15 * time.Millisecond)
		}()
		time.Sleep(20 * time.Millisecond)

		c <- 1
	}()

	<-c

}
