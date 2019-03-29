package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/TingYunAPM/go"
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
	"github.com/martini-contrib/render"
)

var (
	id   int
	name string
)

func main() {
	if err := tingyun.AppInit("tingyun.json"); err != nil {
		fmt.Println(err)
	}
	defer tingyun.AppStop()

	m := martini.Classic()
	m.Use(render.Renderer())
	//普通的GET方式路由
	m.Get("/test", func() string {
		action, _ := tingyun.CreateAction("URI", "/test")
		defer action.Finish()
		connStr := "postgres://dbusername:password@postgredb.local:5432/testdb?sslmode=disable"
		open_component := action.CreateComponent("postgresql.Open")
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			action.SetError(err)
			log.Fatal(err)
		}
		open_component.Finish()
		defer db.Close()
		query_component := action.CreateDBComponent(tingyun.ComponentPostgreSql, "postgredb.local:5432", "testdb", "testtable", "SELECT", "main")
		rows, err := db.Query("SELECT id, name FROM testtable")
		if err != nil {
			action.SetError(err)
			log.Fatal(err)
		}
		var array = make([]string, 0)
		for rows.Next() {
			err := rows.Scan(&id, &name)
			if err != nil {
				action.SetError(err)
				log.Fatal(err)
			}
			array = append(array, "("+strconv.Itoa(id)+","+name+")")
		}
		query_component.Finish()
		ret := "[" + strings.Join(array, ",") + "]"
		return ret
	})
	m.RunOnAddr(":8080")
}
