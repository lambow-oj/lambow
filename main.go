package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.New()
	// Init(r)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run() // listen and serve on 0.0.0.0:8080
}
