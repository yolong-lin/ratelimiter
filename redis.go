package ratelimiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
}

func (r *Redis) Do(key string, expiraion time.Duration) (int, float64, error) {
	ctx := context.Background()
	var val int
	var ttl float64

	_, err := r.client.Pipelined(ctx, func(p redis.Pipeliner) error {
		ok, err := r.client.SetNX(ctx, key, 1, expiraion).Result()
		if err != nil {
			return err
		}

		if ok {
			val = 1
		} else {
			val = int(r.client.Incr(ctx, key).Val())
		}

		ttl = r.client.TTL(ctx, key).Val().Seconds()
		return nil
	})

	if err != nil {
		return 0, 0, err
	}
	return val, ttl, nil
}

func NewRedisStore(addr, username, password string, db int) *Redis {
	r := new(Redis)
	r.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
	})
	return r
}
