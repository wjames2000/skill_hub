package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hpds/skill-hub/pkg/errno"
	"github.com/hpds/skill-hub/pkg/redis"
	"github.com/hpds/skill-hub/pkg/response"
)

type RateLimiter struct {
	client  *redis.Client
	limit   int
	window  time.Duration
	enabled bool
}

func NewRateLimiter(client *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		client:  client,
		limit:   limit,
		window:  window,
		enabled: client != nil,
	}
}

var slidingWindowScript = `
local key = KEYS[1]
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])

redis.call("ZREMRANGEBYSCORE", key, 0, now - window)
local count = redis.call("ZCARD", key)

if count < limit then
	redis.call("ZADD", key, now, now .. ":" .. math.random())
	redis.call("EXPIRE", key, math.ceil(window / 1000))
	return 1
end
return 0
`

func (rl *RateLimiter) Limit(rateKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.enabled {
			c.Next()
			return
		}

		key := fmt.Sprintf("ratelimit:%s:%s", rateKey, c.ClientIP())

		now := time.Now().UnixMilli()
		windowMs := rl.window.Milliseconds()

		ret, err := rl.client.Eval(context.Background(), slidingWindowScript, []string{key}, now, windowMs, rl.limit)
		if err != nil {
			c.Next()
			return
		}

		ok, _ := ret.(int64)
		if ok != 1 {
			response.Error(c, errno.TooManyRequests)
			c.Abort()
			return
		}

		c.Next()
	}
}

func IPRateLimiter(client *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	rl := NewRateLimiter(client, limit, window)
	return rl.Limit("ip")
}
