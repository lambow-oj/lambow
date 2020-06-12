package main

import (
	"git/lambow-oj/lambow/cinex"
	"git/lambow-oj/lambow/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	handler.Init()
	cinex.InitHandler(router)
	router.Run() // listen and serve on 0.0.0.0:8080
}
