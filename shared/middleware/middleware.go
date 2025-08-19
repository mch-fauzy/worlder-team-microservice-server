package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/worlder-team/microservice-server/shared/constants"
)

// ErrorResponse represents error response structure
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// RateLimiter middleware for rate limiting
func RateLimiter(redisClient *redis.Client, limit int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			key := fmt.Sprintf("rate_limit:%s", ip)

			// Get current count
			val, err := redisClient.Get(c.Request().Context(), key).Result()
			if err != nil && err != redis.Nil {
				// If Redis is down, allow the request
				return next(c)
			}

			count := 0
			if val != "" {
				count, _ = strconv.Atoi(val)
			}

			if count >= limit {
				return c.JSON(http.StatusTooManyRequests, ErrorResponse{
					Status:  constants.StatusError,
					Message: constants.ErrRateLimitExceeded,
				})
			}

			// Increment counter
			redisClient.Incr(c.Request().Context(), key)
			redisClient.Expire(c.Request().Context(), key, time.Minute)

			return next(c)
		}
	}
}

// CORS middleware with default settings
func CORS() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			origin := c.Request().Header.Get("Origin")

			// Allow common development origins
			allowedOrigins := []string{
				"http://localhost:3000",
				"http://localhost:8080",
				"http://localhost:8081",
			}

			allowed := false
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					allowed = true
					break
				}
			}

			if allowed || origin == "" {
				c.Response().Header().Set("Access-Control-Allow-Origin", origin)
			}

			c.Response().Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Requested-With")
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")

			if c.Request().Method == "OPTIONS" {
				return c.NoContent(http.StatusOK)
			}

			return next(c)
		}
	}
}

// RequestID middleware adds unique request ID
func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Request().Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = generateRequestID()
			}

			c.Response().Header().Set("X-Request-ID", requestID)
			c.Set("request_id", requestID)

			return next(c)
		}
	}
}

// Authentication middleware
func JWTAuth(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, ErrorResponse{
					Status:  constants.StatusError,
					Message: constants.ErrUnauthorized,
				})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return c.JSON(http.StatusUnauthorized, ErrorResponse{
					Status:  constants.StatusError,
					Message: constants.ErrUnauthorized,
				})
			}

			// TODO: Implement JWT validation
			// For now, we'll skip validation for demo purposes
			c.Set("user_id", "demo_user")

			return next(c)
		}
	}
}

// Security headers middleware
func SecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			c.Response().Header().Set("X-Frame-Options", "DENY")
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
			c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000")

			return next(c)
		}
	}
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}
