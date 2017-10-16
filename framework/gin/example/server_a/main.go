package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/TingYunAPM/go"
	"github.com/TingYunAPM/go/framework/gin"
	//	"github.com/gin-gonic/gin"
)

func main() {
	err := tingyun.AppInit("tingyun.json")
	if err != nil {
		fmt.Println(err)
	}
	defer tingyun.AppStop()
	//	router := gin.Default()
	router := tingyun_gin.Default()
	//	http.ListenAndServe(":8045", router)
	s := &http.Server{
		Addr:           ":8045",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
