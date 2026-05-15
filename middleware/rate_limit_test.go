package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RateLimit())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	return r
}

func TestRateLimiter_Allow(t *testing.T) {
	rl := newRateLimiter(time.Minute, 3)

	assert.True(t, rl.allow("192.168.1.1"))
	assert.True(t, rl.allow("192.168.1.1"))
	assert.True(t, rl.allow("192.168.1.1"))
	assert.False(t, rl.allow("192.168.1.1"))
}

func TestRateLimiter_DifferentKeys(t *testing.T) {
	rl := newRateLimiter(time.Minute, 2)

	assert.True(t, rl.allow("10.0.0.1"))
	assert.True(t, rl.allow("10.0.0.1"))
	assert.False(t, rl.allow("10.0.0.1"))

	assert.True(t, rl.allow("10.0.0.2"))
	assert.True(t, rl.allow("10.0.0.2"))
	assert.False(t, rl.allow("10.0.0.2"))
}

func TestRateLimiter_WindowExpiry(t *testing.T) {
	rl := newRateLimiter(50*time.Millisecond, 2)

	assert.True(t, rl.allow("192.168.1.5"))
	assert.True(t, rl.allow("192.168.1.5"))
	assert.False(t, rl.allow("192.168.1.5"))

	time.Sleep(60 * time.Millisecond)

	assert.True(t, rl.allow("192.168.1.5"))
}

func TestRateLimitMiddleware_HTTP(t *testing.T) {
	globalLimiter = newRateLimiter(time.Minute, 2)
	r := setupRouter()

	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
		req.RemoteAddr = "127.0.0.1:1234"
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}
