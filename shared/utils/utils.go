package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// GenerateID generates a random string ID
func GenerateID(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("ID_%d", time.Now().UnixNano())
	}
	return strings.ToUpper(hex.EncodeToString(bytes))
}

// ParseDuration parses duration string with units
func ParseDuration(duration string) (time.Duration, error) {
	if duration == "" {
		return time.Second, nil
	}
	return time.ParseDuration(duration)
}

// FormatDuration formats duration to string
func FormatDuration(d time.Duration) string {
	return d.String()
}

// StringToInt converts string to int with default value
func StringToInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return defaultValue
}

// StringToFloat64 converts string to float64 with default value
func StringToFloat64(s string, defaultValue float64) float64 {
	if s == "" {
		return defaultValue
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return defaultValue
}

// Contains checks if slice contains element
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// FormatTimestamp formats timestamp to RFC3339
func FormatTimestamp(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ParseTimestamp parses RFC3339 timestamp
func ParseTimestamp(timestamp string) (time.Time, error) {
	return time.Parse(time.RFC3339, timestamp)
}

// GetEnvOrDefault gets environment variable or returns default
func GetEnvOrDefault(key, defaultValue string) string {
	if value := getEnv(key); value != "" {
		return value
	}
	return defaultValue
}

// Mock function for environment variable (replace with os.Getenv in real implementation)
func getEnv(key string) string {
	// This would be os.Getenv(key) in real implementation
	_ = key // Suppress unused parameter warning
	return ""
}
