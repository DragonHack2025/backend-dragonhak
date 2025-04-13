package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
	dummy  bool
}

func NewRateLimiter(addr string) *RateLimiter {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RateLimiter{
		client: client,
		dummy:  false,
	}
}

func NewDummyRateLimiter() *RateLimiter {
	return &RateLimiter{
		dummy: true,
	}
}

func (rl *RateLimiter) Limit(maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rl.dummy {
			c.Next()
			return
		}

		key := fmt.Sprintf("rate_limit:%s", c.ClientIP())
		ctx := context.Background()

		count, err := rl.client.Get(ctx, key).Int()
		if err == redis.Nil {
			rl.client.Set(ctx, key, 1, window)
			c.Next()
			return
		} else if err != nil {
			c.Next() // On error, let the request through
			return
		}

		if count >= maxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests",
				"retry_after": window.Seconds(),
			})
			c.Abort()
			return
		}

		pipe := rl.client.Pipeline()
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, window)
		_, err = pipe.Exec(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		remaining := maxRequests - count - 1
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(window).Unix(), 10))

		c.Next()
	}
}
