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

    store := ratelimiter.NewRedisStore("localhost:6379", "default", "", 0)
	config := ratelimiter.Config{
	    TimeWindow:   time.Hour,
	    RequestQuota: 1000,
	}

    r.Use(ratelimiter.New(config, store))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	
	r.Run()
}
```

