# RateLimiter

Gin middleware to restrict request rate.

## Usage

Install it:

```bash
go get github.com/yolong-lin/ratelimiter
```

import it:

```go
import "github.com/yolong-lin/ratelimiter"
```

## Example

```go
package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yolong-lin/ratelimiter"
)

func main() {
	r := gin.Default()
	r.Use(ratelimiter.New(ratelimiter.Config{
		TimeWindow:    time.Hour,
		RequestQuota:  1000,
		RedisIP:       "localhost",
		RedisPort:     "6379",
		RedisUsername: "default",
		RedisPassword: "",
		RedisDB:       0,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	
	r.Run()
}
```

