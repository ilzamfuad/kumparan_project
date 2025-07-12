package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

const (
	rateLimit     = 10
	rateLimitTime = time.Minute
)

func RateLimitMiddlewareWithRedis(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		key := "rate_limit:" + clientIP

		count, err := redisClient.Incr(key).Result()
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error")
			c.Abort()
			return
		}

		if count == 1 {
			redisClient.Expire(key, rateLimitTime)
		}

		if count > rateLimit {
			c.String(http.StatusTooManyRequests, "Too many requests")
			c.Abort()
			return
		}

		c.Next()
	}
}
