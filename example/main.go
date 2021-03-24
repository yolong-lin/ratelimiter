package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yolong-lin/ratelimiter"
)

func main() {
	r := gin.Default()

	store := ratelimiter.NewRedisStore("redis:6379", "default", "", 0)
	v1Config := ratelimiter.Config{
		TimeWindow:   time.Hour,
		RequestQuota: 1000,
		KeyPrefix:    "v1",
	}
	v2Config := ratelimiter.Config{
		TimeWindow:   time.Hour,
		RequestQuota: 1000,
		KeyPrefix:    "v2",
	}

	v1 := r.Group("v1", ratelimiter.New(v1Config, store))
	{
		v1.GET("/hello", func(c *gin.Context) {
			c.String(200, "Hello, World")
		})
	}

	v2 := r.Group("v2", ratelimiter.New(v2Config, store))
	{
		v2.GET("/hello", func(c *gin.Context) {
			c.String(200, "Hello, Dcard")
		})
	}

	r.Run()
}
