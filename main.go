package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	pongCount := 0
	pingCount := 0

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Content-Length"},
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"count": pingCount,
		})
	})
	r.GET("/pong", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"count": pongCount,
		})
	})

	r.POST("/ping", func(c *gin.Context) {
		pingCount++
		c.JSON(200, gin.H{
			"count": pingCount,
		})
	})
	r.POST("/pong", func(c *gin.Context) {
		pongCount++
		c.JSON(200, gin.H{
			"count": pongCount,
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
