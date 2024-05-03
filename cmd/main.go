package main

import (
	"FeedCraft/internal/router"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()

	router.RegisterRouters(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	listenAddr := os.Getenv("LISTEN_ADDR")
	_ = r.Run(listenAddr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
