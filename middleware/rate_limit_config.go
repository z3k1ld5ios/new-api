package middleware

import (
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimitConfig holds configuration for the rate limiter.
type RateLimitConfig struct {
	WindowSeconds int
	MaxRequests   int
}

// DefaultRateLimitConfig returns sensible defaults.
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		WindowSeconds: 60,
		MaxRequests:   60,
	}
}

// RateLimitConfigFromEnv reads config from environment variables:
//   RATE_LIMIT_WINDOW_SECONDS
//   RATE_LIMIT_MAX_REQUESTS
func RateLimitConfigFromEnv() RateLimitConfig {
	cfg := DefaultRateLimitConfig()

	if v := os.Getenv("RATE_LIMIT_WINDOW_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.WindowSeconds = n
		}
	}

	if v := os.Getenv("RATE_LIMIT_MAX_REQUESTS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.MaxRequests = n
		}
	}

	return cfg
}

// RateLimitWithConfig creates a rate-limit middleware from the given config.
func RateLimitWithConfig(cfg RateLimitConfig) gin.HandlerFunc {
	limiter := newRateLimiter(
		time.Duration(cfg.WindowSeconds)*time.Second,
		cfg.MaxRequests,
	)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.allow(ip) {
			c.AbortWithStatusJSON(429, gin.H{
				"error": gin.H{
					"message": "Rate limit exceeded. Please try again later.",
					"type":    "rate_limit_error",
				},
			})
			return
		}
		c.Next()
	}
}
