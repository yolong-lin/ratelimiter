package ratelimiter

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// All availabe options for the middleware
type Config struct {
	// Maximum number of requests during the time window
	TimeWindow   time.Duration
	RequestQuota int

	// Store Key Prefix
	KeyPrefix string
}

type rateLimiter struct {
	// Maximum number of requests during the time window
	timeWindow   time.Duration
	requestQuota int

	// Store Key Prefix
	keyPrefix string
}

func (r *rateLimiter) setHeader(c *gin.Context, val int, ttl float64) {
	c.Header("X-RateLimit-Reset", strconv.FormatInt(int64(ttl), 10))
	if val > r.requestQuota {
		c.Header("X-RateLimit-Remaining", "0")
		c.AbortWithStatus(http.StatusTooManyRequests)
	} else {
		c.Header("X-RateLimit-Remaining", fmt.Sprint(r.requestQuota-val))
	}
}

func (r *rateLimiter) genrateKey(ip string) string {
	return r.keyPrefix + ":" + ip
}

func newRateLimiter(config Config) *rateLimiter {
	rl := &rateLimiter{
		timeWindow:   config.TimeWindow,
		requestQuota: config.RequestQuota,
	}

	if config.KeyPrefix != "" {
		rl.keyPrefix = config.KeyPrefix
	} else {
		rl.keyPrefix = "ratelimiter"
	}

	return rl
}

func New(config Config, store storage) gin.HandlerFunc {
	rateLimiter := newRateLimiter(config)
	return func(c *gin.Context) {
		key := rateLimiter.genrateKey(c.ClientIP())

		val, ttl, err := store.Do(key, rateLimiter.timeWindow)

		if err != nil {
			panic(err)
		}

		rateLimiter.setHeader(c, val, ttl)
	}
}
