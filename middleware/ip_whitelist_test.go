package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func newWhitelistRouter(cfg *IPWhitelistConfig) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(IPWhitelist(cfg))
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	return r
}

func TestIPWhitelist_EmptyAllowsAll(t *testing.T) {
	r := newWhitelistRouter(DefaultIPWhitelistConfig())
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req.RemoteAddr = "1.2.3.4:5678"
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestIPWhitelist_AllowedIP(t *testing.T) {
	cfg := &IPWhitelistConfig{AllowedIPs: []string{"1.2.3.4"}}
	r := newWhitelistRouter(cfg)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req.RemoteAddr = "1.2.3.4:9999"
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestIPWhitelist_BlockedIP(t *testing.T) {
	cfg := &IPWhitelistConfig{AllowedIPs: []string{"1.2.3.4"}}
	r := newWhitelistRouter(cfg)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req.RemoteAddr = "9.9.9.9:9999"
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestIPWhitelist_CIDRAllowed(t *testing.T) {
	cfg := &IPWhitelistConfig{AllowedCIDRs: []string{"192.168.1.0/24"}}
	r := newWhitelistRouter(cfg)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("X-Real-IP", "192.168.1.55")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestIPWhitelistConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("IP_WHITELIST_IPS")
	os.Unsetenv("IP_WHITELIST_CIDRS")
	cfg := IPWhitelistConfigFromEnv()
	if len(cfg.AllowedIPs) != 0 || len(cfg.AllowedCIDRs) != 0 {
		t.Fatal("expected empty config from unset env vars")
	}
}

func TestIPWhitelistConfigFromEnv_Custom(t *testing.T) {
	os.Setenv("IP_WHITELIST_IPS", "10.0.0.1, 10.0.0.2")
	os.Setenv("IP_WHITELIST_CIDRS", "172.16.0.0/12")
	defer os.Unsetenv("IP_WHITELIST_IPS")
	defer os.Unsetenv("IP_WHITELIST_CIDRS")
	cfg := IPWhitelistConfigFromEnv()
	if len(cfg.AllowedIPs) != 2 {
		t.Fatalf("expected 2 IPs, got %d", len(cfg.AllowedIPs))
	}
	if len(cfg.AllowedCIDRs) != 1 {
		t.Fatalf("expected 1 CIDR, got %d", len(cfg.AllowedCIDRs))
	}
}
