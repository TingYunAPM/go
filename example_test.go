// Copyright 2016-2019 冯立强 fenglq@tingyun.com.  All rights reserved.

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
func ExampleComponent_CreateTrackId() {
	tingyun.AppInit("tingyun.json")
	defer tingyun.AppStop()
	tingyun_beego.Handler("/extern", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if track_id := r.Header.Get("X-Tingyun-Id"); len(track_id) > 0 {
			tingyun_beego.GetAction(w).SetTrackId(track_id)
		}
		url := "http://192.168.100.1/x.php"
		c := tingyun_beego.GetAction(w).CreateExternalComponent(url, "main.Extern")
		result, response, err := HttpGet(url, map[string]string{"X-Tingyun-Id": c.CreateTrackId()}) //由外部调用组件生成追踪id
		if response != nil {
			if tx_data := response.Header.Get("X-Tingyun-Tx-Data"); len(tx_data) > 0 {
				c.SetTxData(tx_data)
			}
			if tx_data := action.GetTxData(); len(tx_data) > 0 {
				w.Header().Set("X-Tingyun-Tx-Data", tx_data)
			}
			w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
			w.WriteHeader(response.StatusCode)
			w.Write(result)
		} else {
			fmt.Fprintf(w, "%q", "/extern error:"+err.Error())
		}
		c.Finish()
	}))
	tingyun_beego.Run()
}
func ExampleAction_GetTxData() {
	tingyun.AppInit("tingyun.json")
	defer tingyun.AppStop()
	tingyun_beego.Handler("/extern", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if track_id := r.Header.Get("X-Tingyun-Id"); len(track_id) > 0 {
			tingyun_beego.GetAction(w).SetTrackId(track_id)
		}
		url := "http://192.168.100.1/x.php"
		c := tingyun_beego.GetAction(w).CreateExternalComponent(url, "main.Extern")
		result, response, err := HttpGet(url, map[string]string{"X-Tingyun-Id": c.CreateTrackId()})
		if response != nil {
			if tx_data := response.Header.Get("X-Tingyun-Tx-Data"); len(tx_data) > 0 {
				c.SetTxData(tx_data)
			}
			//从被调用端action获取事务性能数据
			if tx_data := tingyun_beego.GetAction(w).GetTxData(); len(tx_data) > 0 {
				w.Header().Set("X-Tingyun-Tx-Data", tx_data)
			}
			w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
			w.WriteHeader(response.StatusCode)
			w.Write(result)
		} else {
			fmt.Fprintf(w, "%q", "/extern error:"+err.Error())
		}
		c.Finish()
	}))
	tingyun_beego.Run()
}
func ExampleAction_SetTrackId() {
	tingyun.AppInit("tingyun.json")
	defer tingyun.AppStop()
	tingyun_beego.Handler("/extern", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if track_id := r.Header.Get("X-Tingyun-Id"); len(track_id) > 0 {
			tingyun_beego.GetAction(w).SetTrackId(track_id) //将调用端传递来的追踪id保存到action
		}
		url := "http://192.168.100.1/x.php"
		c := tingyun_beego.GetAction(w).CreateExternalComponent(url, "main.Extern")
		result, response, err := HttpGet(url, map[string]string{"X-Tingyun-Id": c.CreateTrackId()})
		if response != nil {
			if tx_data := response.Header.Get("X-Tingyun-Tx-Data"); len(tx_data) > 0 {
				c.SetTxData(tx_data)
			}
			if tx_data := action.GetTxData(); len(tx_data) > 0 {
				w.Header().Set("X-Tingyun-Tx-Data", tx_data)
			}
			w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
			w.WriteHeader(response.StatusCode)
			w.Write(result)
		} else {
			fmt.Fprintf(w, "%q", "/extern error:"+err.Error())
		}
		c.Finish()
	}))
	tingyun_beego.Run()
}

func ExampleComponent_SetTxData() {
	tingyun.AppInit("tingyun.json")
	defer tingyun.AppStop()
	tingyun_beego.Handler("/extern", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if track_id := r.Header.Get("X-Tingyun-Id"); len(track_id) > 0 {
			tingyun_beego.GetAction(w).SetTrackId(track_id)
		}
		url := "http://192.168.100.1/x.php"
		c := tingyun_beego.GetAction(w).CreateExternalComponent(url, "main.Extern")
		result, response, err := HttpGet(url, map[string]string{"X-Tingyun-Id": c.CreateTrackId()})
		if response != nil {
			if tx_data := response.Header.Get("X-Tingyun-Tx-Data"); len(tx_data) > 0 {
				c.SetTxData(tx_data) //将被调用端返回的事务数据保存到外部调用组件
			}
			if tx_data := action.GetTxData(); len(tx_data) > 0 {
				w.Header().Set("X-Tingyun-Tx-Data", tx_data)
			}
			w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
			w.WriteHeader(response.StatusCode)
			w.Write(result)
		} else {
			fmt.Fprintf(w, "%q", "/extern error:"+err.Error())
		}
		c.Finish()
	}))
	tingyun_beego.Run()
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
func ExampleComponent_AppendSQL() {
	tingyun.AppInit("tingyun.json")
	defer tingyun.AppStop()
	action, _ := tingyun.CreateAction("/login", "main.ExampleCreateAction")
	dbComponent := action.CreateDBComponent(tingyun.ComponentPostgreSql, "postgredb.local:5432", "testdb", "testtable", "SELECT", "main")
	dbComponent.AppendSQL("SELECT id, name FROM testtable")
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
