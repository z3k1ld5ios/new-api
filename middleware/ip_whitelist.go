package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// IPWhitelistConfig holds configuration for IP whitelist middleware
type IPWhitelistConfig struct {
	AllowedCIDRs []string
	AllowedIPs   []string
	nets         []*net.IPNet
	ips          map[string]struct{}
}

// compile parses CIDR and IP strings into usable net structures
func (c *IPWhitelistConfig) compile() error {
	c.ips = make(map[string]struct{}, len(c.AllowedIPs))
	for _, ip := range c.AllowedIPs {
		c.ips[strings.TrimSpace(ip)] = struct{}{}
	}
	for _, cidr := range c.AllowedCIDRs {
		_, network, err := net.ParseCIDR(strings.TrimSpace(cidr))
		if err != nil {
			return err
		}
		c.nets = append(c.nets, network)
	}
	return nil
}

// isAllowed checks whether the given IP address is permitted
func (c *IPWhitelistConfig) isAllowed(ipStr string) bool {
	if _, ok := c.ips[ipStr]; ok {
		return true
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	for _, network := range c.nets {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

// getClientIP extracts the real client IP from the request.
// NOTE: X-Forwarded-For can be spoofed by clients; only trust it if your
// reverse proxy (nginx, etc.) is the one setting it.
func getClientIP(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return ip
}

// IPWhitelist returns a Gin middleware that restricts access to allowed IPs/CIDRs.
// If config has no entries, all IPs are allowed.
func IPWhitelist(cfg *IPWhitelistConfig) gin.HandlerFunc {
	if err := cfg.compile(); err != nil {
		panic("ip_whitelist: invalid config: " + err.Error())
	}
	return func(c *gin.Context) {
		if len(cfg.ips) == 0 && len(cfg.nets) == 0 {
			c.Next()
			return
		}
		clientIP := getClientIP(c)
		if !cfg.isAllowed(clientIP) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":  "IP address not allowed",
				"remote": clientIP,
			})
			return
		}
		c.Next()
	}
}
