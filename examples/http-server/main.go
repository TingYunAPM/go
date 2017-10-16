package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/TingYunAPM/go"
)

func main() {
	//初始化tingyun: 应用名称、license等在tingyun.json中配置
	err := tingyun.AppInit("tingyun.json")
	if err != nil {
		fmt.Println(err)
	}
	defer tingyun.AppStop()

	tingyun.HandleFunc("/", webHandler)
	tingyun.HandleFunc("/login", loginHandler)
	tingyun.HandleFunc("/mysql", mysqlHandler)
	tingyun.HandleFunc("/redis", redisHandler)
	tingyun.HandleFunc("/external", externalHandler)
	tingyun.HandleFunc("/error", errorHandler)
	tingyun.HandleFunc("/ignore", ignoreHandler)
	tingyun.HandleFunc("/cross", crossHandler)
	tingyun.HandleFunc("/crossA", crossAHandler)
	tingyun.HandleFunc("/crossB", crossBHandler)

	port := flag.String("port", ":8080", "http listen port")
	flag.Parse()

	e := http.ListenAndServe(*port, nil)
	if e != nil {
		log.Fatalln("ListenAndServe: ", e)
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
	<a href=/cross>cross</a> 	
	<a href=/404>404</a> 
	<a href=/error>error</a> 
	<a href=/ignore>ignore</a> 
	<hr /> `)
	w.Write(body)
	time.Sleep(time.Second)
	submenu(w, r)
	time.Sleep(200 * time.Millisecond)
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
	time.Sleep(100 * time.Millisecond)
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
	c = tingyun.GetAction(w).CreateDBComponent(tingyun.ComponentMongo, "", "", "user", "Delete", "mongo.do")
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
	client := http.Client{
		Timeout: time.Duration(3 * time.Second),
	}
	resp, err := client.Get(url)
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			w.Write(body)
		}
	}
	c.Finish()

	//模拟thrift
	u := "thrift://www.baidu.com/thrift/get"
	c = tingyun.GetAction(w).CreateExternalComponent(u, "thrift.Get")
	c.Finish()

}

func crossHandler(w http.ResponseWriter, r *http.Request) {

	defer tingyun.GetAction(w).CreateComponent("main.crossHandler").Finish()
	//
	menu(w, r)

	time.Sleep(1 * time.Second)

	//cross app
	u := "http://127.0.0.1:8081/crossA"
	c := tingyun.GetAction(w).CreateExternalComponent(u, "http.Get")
	track := c.CreateTrackId()

	client := http.Client{
		Timeout: time.Duration(20 * time.Second),
	}
	resp, err := client.Post(u,
		"application/x-www-form-urlencoded",
		strings.NewReader("track="+url.QueryEscape(track)))
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			w.Write(body)
		}
	}
	c.Finish()
}

func crossAHandler(w http.ResponseWriter, r *http.Request) {
	track := r.FormValue("track")
	tingyun.GetAction(w).SetTrackId(track)
	defer tingyun.GetAction(w).CreateComponent("main.crossAHandler").Finish()
	//
	menu(w, r)

	w.Write([]byte("cross A"))

	time.Sleep(1 * time.Second)

	//cross self app
	u := "http://127.0.0.1:8081/crossB"
	c := tingyun.GetAction(w).CreateExternalComponent(u, "http.Get1")
	track = c.CreateTrackId()

	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
	resp, err := client.Post(u,
		"application/x-www-form-urlencoded",
		strings.NewReader("track="+url.QueryEscape(track)))
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			w.Write(body)
		}
	}
	c.Finish()
	//
	//cross another app
	u = "http://127.0.0.1:8082/crossB"
	c = tingyun.GetAction(w).CreateExternalComponent(strings.Replace(u, "http", "thrift", 1), "thrift.Call")
	track = c.CreateTrackId()

	client = http.Client{
		Timeout: time.Duration(6 * time.Second),
	}
	resp, err = client.Post(u,
		"application/x-www-form-urlencoded",
		strings.NewReader("track="+url.QueryEscape(track)))
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			w.Write(body)
		}
	}
	c.Finish()

}

func crossBHandler(w http.ResponseWriter, r *http.Request) {
	track := r.FormValue("track")
	tingyun.GetAction(w).SetTrackId(track)
	defer tingyun.GetAction(w).CreateComponent("main.crossBHandler").Finish()
	//
	menu(w, r)

	rand.Seed(int64(time.Now().Nanosecond()))
	time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)

	w.Write([]byte("cross B"))
}
