package ratelimiter

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// All availabe options for the middleware
type Config struct {
	// Maximum number of requests during the time window
	TimeWindow   time.Duration
	RequestQuota int

	// TODO:Abstract
	// Redis options
	RedisIP       string
	RedisPort     string
	RedisUsername string
	RedisPassword string
	RedisDB       int

	// Store Key Prefix
	KeyPrefix string
}

type rateLimiter struct {
	// Maximum number of requests during the time window
	timeWindow   time.Duration
	requestQuota int

	// Redis Client
	redisClient *redis.Client

	// Store Key Prefix
	keyPrefix string
}

func (r *rateLimiter) conduct(c *gin.Context) {
	ctx := context.Background()
	key := r.genrateKey(c.ClientIP())

	var remaining int
	var ttl float64
	// Using pipelined to speed up query
	_, err := r.redisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		// create key, return true if key created
		ok, err := r.redisClient.SetNX(ctx, key, 1, r.timeWindow).Result()
		if err != nil {
			return err
		}

		if ok {
			remaining = r.requestQuota - 1
		} else {
			val := r.redisClient.Incr(ctx, key).Val()
			remaining = r.requestQuota - int(val)
		}

		ttl = r.redisClient.TTL(ctx, key).Val().Seconds()

		return nil
	})

	if err != nil {
		panic(err)
	}

	c.Header("X-RateLimit-Reset", strconv.FormatInt(int64(ttl), 10))
	if remaining < 0 {
		c.Header("X-RateLimit-Remaining", "0")
		c.AbortWithStatus(http.StatusTooManyRequests)
	} else {
		c.Header("X-RateLimit-Remaining", fmt.Sprint(remaining))
	}

}

func (r *rateLimiter) genrateKey(ip string) string {
	return r.keyPrefix + ":" + ip
}

func newRateLimiter(config Config) *rateLimiter {
	rl := &rateLimiter{
		timeWindow:   config.TimeWindow,
		requestQuota: config.RequestQuota,
		redisClient: redis.NewClient(&redis.Options{
			Addr:     config.RedisIP + ":" + config.RedisPort,
			Username: config.RedisUsername,
			Password: config.RedisPassword,
			DB:       config.RedisDB,
		}),
	}

	if config.KeyPrefix != "" {
		rl.keyPrefix = config.KeyPrefix
	} else {
		rl.keyPrefix = "ratelimiter"
	}

	return rl
}

func New(config Config) gin.HandlerFunc {
	ratelimiter := newRateLimiter(config)
	return func(c *gin.Context) {
		ratelimiter.conduct(c)
	}
}
