package ratelimiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var script = `
local val, ttl

local ok = redis.call('SET', KEYS[1], '1', 'EX', ARGV[1], 'NX')
if ok then	
	val = 1
else
	val = redis.call('INCR', KEYS[1])
end

ttl = redis.call('TTL', KEYS[1])
return {val, ttl}
`

type Redis struct {
	client *redis.Client
}

func (r *Redis) Do(key string, expiraion time.Duration) (int, int, error) {
	ctx := context.Background()
	expiraion /= time.Second
	res, err := r.client.Eval(ctx, script, []string{key}, int(expiraion)).Result()

	if err != nil {
		return 0, 0, err
	}

	val := res.([]interface{})

	return int(val[0].(int64)), int(val[1].(int64)), nil
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
