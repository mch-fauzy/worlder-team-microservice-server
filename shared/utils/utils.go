package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
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
// Returns 1 second as default if duration is empty string
func ParseDuration(duration string) (time.Duration, error) {
	if duration == "" {
		return time.Second, nil
	}
	return time.ParseDuration(duration)
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

// GetEnvOrDefault gets environment variable or returns default
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ParseInt parses string to int with fallback to 0
func ParseInt(s string) int {
	if value, err := strconv.Atoi(s); err == nil {
		return value
	}
	return 0
}

// ParseDurationOrZero parses duration string with fallback to 0
func ParseDurationOrZero(s string) time.Duration {
	if duration, err := time.ParseDuration(s); err == nil {
		return duration
	}
	return 0
}
