package ratelimiter

import (
	"time"
)

type storage interface {
	// Give a key and a expiration time, return value and ttl
	Do(string, time.Duration) (int, float64, error)
}
