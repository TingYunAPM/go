package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/TingYunAPM/go"
)

func main() {
	//初始化tingyun: 应用名称、license等在tingyun.json中配置
	tingyun.AppInit("tingyun.json")
	defer tingyun.AppStop()
	//
	tingyun.Handle("/", http.HandlerFunc(webHandler))
	tingyun.Handle("/login", http.HandlerFunc(loginHandler))
	tingyun.Handle("/mysql", http.HandlerFunc(mysqlHandler))
	tingyun.Handle("/redis", http.HandlerFunc(redisHandler))
	tingyun.Handle("/external", http.HandlerFunc(externalHandler))
	tingyun.Handle("/error", http.HandlerFunc(errorHandler))
	tingyun.Handle("/ignore", http.HandlerFunc(ignoreHandler))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func submenu(w http.ResponseWriter, r *http.Request) {
	defer tingyun.GetAction(w).CreateComponent("submenu").Finish()
	time.Sleep(300 * time.Millisecond)
}

func menu(w http.ResponseWriter, r *http.Request) {
	defer tingyun.GetAction(w).CreateComponent("menu").Finish()
	//
	body := []byte(
		`<a href=/>home</a> 
	<a href=/login>login</a> 
	<a href=/mysql>mysql</a> 
	<a href=/redis>redis</a> 
	<a href=/external>external</a> 
	<a href=/404>404</a> 
	<a href=/error>error</a> 
	<a href=/ignore>ignore</a> 
	<hr /> `)
	w.Write(body)
	time.Sleep(time.Second)
	submenu(w, r)
}

func ignoreHandler(w http.ResponseWriter, r *http.Request) {
	tingyun.GetAction(w).Ignore()
	//
	menu(w, r)
	w.Write([]byte("ignored"))
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	menu(w, r)
	db, err := sql.Open("mysql", "admin:xxx@tcp(tingyun-mysql-tingyunorg.myalauda.cn:33306)/mysql")
	defer db.Close()
	rows, err := db.Query("SELECT distinct u.user as username FROM user u")
	if err != nil {
		w.Write([]byte(err.Error()))
		tingyun.GetAction(w).SetError(err)
	} else {
		w.Write([]byte("<br/>list:"))
		var username string
		for rows.Next() {
			rows.Scan(&username)
			w.Write([]byte(" " + username))
		}
	}
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	defer tingyun.GetAction(w).CreateComponent("main.webHandler").Finish()
	//
	if r.URL.Path == "/404" {
		w.WriteHeader(404)
		return
	}
	//
	menu(w, r)
	w.Write([]byte("wellcome"))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	defer tingyun.GetAction(w).CreateComponent("main.loginhandler").Finish()
	//
	menu(w, r)
	body := []byte(`
		<form action="/mysql" method="POST">
  		First name: <input type="text" name="fname" />
		Last name: <input type="text" name="lname" />
  		<input type="submit" value="Submit" />
		</form>`)
	w.Write(body)
}

func mysqlHandler(w http.ResponseWriter, r *http.Request) {
	defer tingyun.GetAction(w).CreateComponent("main.mysqlhandler").Finish()
	//
	menu(w, r)
	r.ParseForm()
	body := "User input:"
	for k, v := range r.Form {
		body += " " + k + ": " + strings.Join(v, ",")
	}
	w.Write([]byte(body))

	o := tingyun.GetAction(w).CreateComponent("mysql.open")
	db, err := sql.Open("mysql", "admin:tingyun20161109@tcp(tingyun-mysql-tingyunorg.myalauda.cn:33306)/mysql")
	if err != nil {
		w.Write([]byte(err.Error()))
		tingyun.GetAction(w).SetError(err)
	}
	time.Sleep(20 * time.Millisecond)
	o.Finish()
	if err != nil {
		w.Write([]byte(err.Error()))
		tingyun.GetAction(w).SetError(err)
	} else {
		defer db.Close()
		c := tingyun.GetAction(w).CreateDBComponent(tingyun.ComponentMysql, "", "", "user", "SELECT", "db.Query")
		rows, err := db.Query("SELECT distinct u.user as username FROM user u")
		if err != nil {
			w.Write([]byte(err.Error()))
			tingyun.GetAction(w).SetError(err)
		} else {
			w.Write([]byte("<br/>list:"))
			var username string
			for rows.Next() {
				rows.Scan(&username)
				w.Write([]byte(" " + username))
			}
		}
		time.Sleep(1000 * time.Millisecond)
		c.Finish()
	}
}

func redisHandler(w http.ResponseWriter, r *http.Request) {
	defer tingyun.GetAction(w).CreateComponent("main.redis").Finish()
	menu(w, r)
	c := tingyun.GetAction(w).CreateDBComponent(tingyun.ComponentRedis, "", "", "", "Get", "redis.do")
	time.Sleep(8 * time.Millisecond)
	c.Finish()
	c = tingyun.GetAction(w).CreateDBComponent(tingyun.ComponentRedis, "", "", "", "Set", "redis.do")
	time.Sleep(6 * time.Millisecond)
	c.Finish()
	c = tingyun.GetAction(w).CreateDBComponent(tingyun.ComponentRedis, "", "", "", "Insert", "redis.do")
	time.Sleep(24 * time.Millisecond)
	c.Finish()
	c = tingyun.GetAction(w).CreateDBComponent(tingyun.ComponentRedis, "", "", "", "Delete", "redis.do")
	time.Sleep(15 * time.Millisecond)
	c.Finish()
	w.Write([]byte("redis"))
}

func externalHandler(w http.ResponseWriter, r *http.Request) {
	defer tingyun.GetAction(w).CreateComponent("main.externalhandler").Finish()
	//
	menu(w, r)
	url := "http://www.baidu.com/"
	c := tingyun.GetAction(w).CreateExternalComponent(url, "http.Get")
	resp, err := http.Get(url)
	c.Finish()
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			w.Write(body)
		}
	}
}
