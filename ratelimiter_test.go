package ratelimiter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	rr        *httptest.ResponseRecorder
	got       string
	keyPrefix string = "KeyPrefix"
)

func serve(route http.Handler) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", "/ping", nil)
	rr = httptest.NewRecorder()
	route.ServeHTTP(rr, req)
	return rr
}

func TestMiddlewareWithRedis(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	gin.SetMode(gin.TestMode)
	route := gin.New()

	config := Config{
		TimeWindow:   time.Hour,
		RequestQuota: 5,
		KeyPrefix:    keyPrefix,
	}
	store := NewRedisStore(s.Addr(), "default", "", 0)

	route.Use(New(config, store))
	route.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	rr = serve(route)
	assert.Equal(t, 200, rr.Result().StatusCode)
	assert.Equal(t, "4", rr.Header().Get("X-RateLimit-Remaining"))
	assert.NotEmpty(t, rr.Header().Get("X-RateLimit-Reset"))
	got, _ = s.Get(keyPrefix + ":")
	assert.Equal(t, "1", got)

	rr = serve(route)
	assert.Equal(t, 200, rr.Result().StatusCode)
	assert.Equal(t, "3", rr.Header().Get("X-RateLimit-Remaining"))
	assert.NotEmpty(t, rr.Header().Get("X-RateLimit-Reset"))
	got, _ = s.Get(keyPrefix + ":")
	assert.Equal(t, "2", got)

	rr = serve(route)
	rr = serve(route)

	rr = serve(route)
	assert.Equal(t, 200, rr.Result().StatusCode)
	assert.Equal(t, "0", rr.Header().Get("X-RateLimit-Remaining"))
	assert.NotEmpty(t, rr.Header().Get("X-RateLimit-Reset"))
	got, _ = s.Get(keyPrefix + ":")
	assert.Equal(t, "5", got)

	rr = serve(route)
	assert.Equal(t, 429, rr.Result().StatusCode)
	assert.Equal(t, "0", rr.Header().Get("X-RateLimit-Remaining"))
	assert.NotEmpty(t, rr.Header().Get("X-RateLimit-Reset"))
	got, _ = s.Get(keyPrefix + ":")
	assert.Equal(t, "6", got)

	s.FastForward(time.Hour)

	rr = serve(route)
	assert.Equal(t, 200, rr.Result().StatusCode)
	assert.Equal(t, "4", rr.Header().Get("X-RateLimit-Remaining"))
	assert.NotEmpty(t, rr.Header().Get("X-RateLimit-Reset"))
	got, _ = s.Get(keyPrefix + ":")
	assert.Equal(t, "1", got)
}
