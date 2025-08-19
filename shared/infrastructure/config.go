package infrastructure

import (
	"os"
	"strconv"
	"time"
)

// GetJWTExpiration returns JWT expiration in seconds from env or fallback
func GetJWTExpiration() int {
	if value := os.Getenv("JWT_EXPIRATION_SECONDS"); value != "" {
		if seconds, err := strconv.Atoi(value); err == nil {
			return seconds
		}
	}
	return 24 * 60 * 60 // 24 hours in seconds
}

// GetDefaultRateLimit returns rate limit from env or fallback
func GetDefaultRateLimit() int {
	if value := os.Getenv("RATE_LIMIT"); value != "" {
		if limit, err := strconv.Atoi(value); err == nil {
			return limit
		}
	}
	return 100 // requests per minute
}

// GetMaxDBConnections returns max DB connections from env or fallback
func GetMaxDBConnections() int {
	if value := os.Getenv("MAX_DB_CONNECTIONS"); value != "" {
		if conns, err := strconv.Atoi(value); err == nil {
			return conns
		}
	}
	return 100
}

// GetMaxIdleConnections returns max idle connections from env or fallback
func GetMaxIdleConnections() int {
	if value := os.Getenv("MAX_IDLE_CONNECTIONS"); value != "" {
		if conns, err := strconv.Atoi(value); err == nil {
			return conns
		}
	}
	return 10
}

// GetDefaultCacheTTL returns cache TTL from env or fallback
func GetDefaultCacheTTL() time.Duration {
	if value := os.Getenv("CACHE_TTL_SECONDS"); value != "" {
		if seconds, err := strconv.Atoi(value); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}
	return 5 * 60 * time.Second // 5 minutes
}
