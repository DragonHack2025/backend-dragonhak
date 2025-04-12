package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
}

func NewRateLimiter(redisAddr string) *RateLimiter {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return &RateLimiter{client: client}
}

func (rl *RateLimiter) Limit(maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use IP address as the key
		key := c.ClientIP()

		// Get current count
		count, err := rl.client.Get(c, key).Int()
		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		// If count exceeds max requests, return error
		if count >= maxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests",
				"retry_after": window.Seconds(),
			})
			c.Abort()
			return
		}

		// Increment count
		pipe := rl.client.Pipeline()
		pipe.Incr(c, key)
		pipe.Expire(c, key, window)
		_, err = pipe.Exec(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		// Add remaining requests to header
		remaining := maxRequests - count - 1
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(window).Unix(), 10))

		c.Next()
	}
}
