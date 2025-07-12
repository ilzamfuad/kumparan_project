package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitMiddlewareWithRedis(t *testing.T) {
	// Mock Redis client
	mockRedis, err := miniredis.Run()
	assert.NoError(t, err)
	defer mockRedis.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mockRedis.Addr(), // like "localhost:6379"
	})

	// Flush Redis before starting the test
	rdb.FlushAll()

	// Create a Gin router with the middleware
	r := gin.Default()
	r.Use(RateLimitMiddlewareWithRedis(rdb))
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	t.Run("Allow requests under the rate limit", func(t *testing.T) {
		for i := 0; i < rateLimit; i++ {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = "127.0.0.1:12345"
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "OK", w.Body.String())
		}
	})

	t.Run("Block requests exceeding the rate limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.Equal(t, "Too many requests", w.Body.String())
	})

}
