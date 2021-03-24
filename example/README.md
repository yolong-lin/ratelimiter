## 範例程式

```go
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
```

啟動 Redis 並執行 `main.go`

```shell
$ docker-compose up -d
```

檢查 Header

```shell
$ curl -sI -X GET 127.0.0.1:8080/v1/hello | grep -iE 'X-Ratelimit|HTTP/1.1'
HTTP/1.1 200 OK
X-Ratelimit-Remaining: 999
X-Ratelimit-Reset: 3600

$ curl -sI -X GET 127.0.0.1:8080/v2/hello | grep -iE 'X-Ratelimit|HTTP/1.1'
HTTP/1.1 200 OK
X-Ratelimit-Remaining: 999
X-Ratelimit-Reset: 3600
```

