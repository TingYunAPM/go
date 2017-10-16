package main

import (
	"fmt"

	"net/http"

	"github.com/TingYunAPM/go"
	"github.com/TingYunAPM/go/framework/gin"
	"github.com/gin-gonic/gin"
)

func main() {

	err := tingyun.AppInit("tingyun.json")
	if err != nil {
		fmt.Println(err)
	}
	defer tingyun.AppStop()
	fmt.Println("Gin Test case C:")
	//	router := gin.Default()
	router := tingyun_gin.Default()

	router.GET("/someGet", getting)
	router.POST("/somePost", posting)
	router.PUT("/somePut", putting)
	router.DELETE("/someDelete", deleting)
	router.PATCH("/somePatch", patching)
	router.HEAD("/someHead", head)
	router.OPTIONS("/someOptions", options)
	router.Run(":8021")
}

func getting(c *gin.Context) {
	c.String(http.StatusOK, "/someGet, getting")
}
func posting(c *gin.Context) {
	c.String(http.StatusOK, "/somePost, posting")
}

func putting(c *gin.Context) {
	c.String(http.StatusOK, "/somePut, putting")
}
func deleting(c *gin.Context) {
	c.String(http.StatusOK, "/someDelete, deleting")
}

func patching(c *gin.Context) {
	c.String(http.StatusOK, "/somePatch, patching")
}
func head(c *gin.Context) {
	c.String(http.StatusOK, "/someHead, head")
}
func options(c *gin.Context) {
	c.String(http.StatusOK, "/someOptions, options")
}
