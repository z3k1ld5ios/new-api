package middleware

import (
	"os"
	"strings"
)

// DefaultIPWhitelistConfig returns an empty whitelist config (allow all).
func DefaultIPWhitelistConfig() *IPWhitelistConfig {
	return &IPWhitelistConfig{}
}

// IPWhitelistConfigFromEnv builds an IPWhitelistConfig from environment variables.
//
// Environment variables:
//   IP_WHITELIST_IPS   — comma-separated list of exact IPs
//   IP_WHITELIST_CIDRS — comma-separated list of CIDR ranges
//
// Note: both variables also support semicolons as an alternative separator,
// which is handy when commas conflict with certain shell or Docker env configs.
func IPWhitelistConfigFromEnv() *IPWhitelistConfig {
	cfg := DefaultIPWhitelistConfig()

	if raw := os.Getenv("IP_WHITELIST_IPS"); raw != "" {
		// Normalise semicolon separators to commas before splitting.
		raw = strings.ReplaceAll(raw, ";", ",")
		for _, ip := range strings.Split(raw, ",") {
			ip = strings.TrimSpace(ip)
			if ip != "" {
				cfg.AllowedIPs = append(cfg.AllowedIPs, ip)
			}
		}
	}

	if raw := os.Getenv("IP_WHITELIST_CIDRS"); raw != "" {
		// Normalise semicolon separators to commas before splitting.
		raw = strings.ReplaceAll(raw, ";", ",")
		for _, cidr := range strings.Split(raw, ",") {
			cidr = strings.TrimSpace(cidr)
			if cidr != "" {
				cfg.AllowedCIDRs = append(cfg.AllowedCIDRs, cidr)
			}
		}
	}

	return cfg
}
