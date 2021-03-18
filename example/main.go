package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yolong-lin/ratelimiter"
)

func main() {
	r := gin.Default()

	store := ratelimiter.NewRedisStore("redis:6379", "default", "", 0)
	config := ratelimiter.Config{
		TimeWindow:   time.Hour,
		RequestQuota: 1000,
		KeyPrefix:    "rl",
	}

	r.Use(ratelimiter.New(config, store))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run()
}
