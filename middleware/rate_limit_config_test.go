package middleware

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultRateLimitConfig(t *testing.T) {
	cfg := DefaultRateLimitConfig()
	assert.Equal(t, 60, cfg.WindowSeconds)
	assert.Equal(t, 60, cfg.MaxRequests)
}

func TestRateLimitConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("RATE_LIMIT_WINDOW_SECONDS")
	os.Unsetenv("RATE_LIMIT_MAX_REQUESTS")

	cfg := RateLimitConfigFromEnv()
	assert.Equal(t, 60, cfg.WindowSeconds)
	assert.Equal(t, 60, cfg.MaxRequests)
}

func TestRateLimitConfigFromEnv_Custom(t *testing.T) {
	os.Setenv("RATE_LIMIT_WINDOW_SECONDS", "30")
	os.Setenv("RATE_LIMIT_MAX_REQUESTS", "100")
	defer os.Unsetenv("RATE_LIMIT_WINDOW_SECONDS")
	defer os.Unsetenv("RATE_LIMIT_MAX_REQUESTS")

	cfg := RateLimitConfigFromEnv()
	assert.Equal(t, 30, cfg.WindowSeconds)
	assert.Equal(t, 100, cfg.MaxRequests)
}

func TestRateLimitConfigFromEnv_InvalidValues(t *testing.T) {
	os.Setenv("RATE_LIMIT_WINDOW_SECONDS", "invalid")
	os.Setenv("RATE_LIMIT_MAX_REQUESTS", "-5")
	defer os.Unsetenv("RATE_LIMIT_WINDOW_SECONDS")
	defer os.Unsetenv("RATE_LIMIT_MAX_REQUESTS")

	cfg := RateLimitConfigFromEnv()
	// Should fall back to defaults on invalid values
	assert.Equal(t, 60, cfg.WindowSeconds)
	assert.Equal(t, 60, cfg.MaxRequests)
}

func TestRateLimitWithConfig(t *testing.T) {
	cfg := RateLimitConfig{WindowSeconds: 60, MaxRequests: 1}
	mw := RateLimitWithConfig(cfg)
	assert.NotNil(t, mw)
}
