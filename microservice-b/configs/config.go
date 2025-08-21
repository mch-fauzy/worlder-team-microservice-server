package configs

import (
	"time"

	"github.com/worlder-team/microservice-server/shared/utils"
)

// Config holds all configuration for microservice-b
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	GRPC      GRPCConfig
	JWT       JWTConfig
	RateLimit RateLimitConfig
	Cache     CacheConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         string
	LogLevel     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

// GRPCConfig holds gRPC server configuration
type GRPCConfig struct {
	Port string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string
	Issuer     string
	Expiration time.Duration
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	TTL time.Duration
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         utils.GetEnvOrDefault("MICROSERVICE_B_PORT", "8080"),
			LogLevel:     utils.GetEnvOrDefault("LOG_LEVEL", "info"),
			ReadTimeout:  utils.ParseDurationOrZero(utils.GetEnvOrDefault("READ_TIMEOUT", "10s")),
			WriteTimeout: utils.ParseDurationOrZero(utils.GetEnvOrDefault("WRITE_TIMEOUT", "10s")),
		},
		Database: DatabaseConfig{
			Host:     utils.GetEnvOrDefault("DB_HOST", "localhost"),
			Port:     utils.GetEnvOrDefault("DB_PORT", "3306"),
			Name:     utils.GetEnvOrDefault("DB_NAME", "sensor_data"),
			Username: utils.GetEnvOrDefault("DB_USER", "root"),
			Password: utils.GetEnvOrDefault("DB_PASSWORD", "password"),
		},
		Redis: RedisConfig{
			Host:     utils.GetEnvOrDefault("REDIS_HOST", "localhost"),
			Port:     utils.GetEnvOrDefault("REDIS_PORT", "6379"),
			Password: utils.GetEnvOrDefault("REDIS_PASSWORD", ""),
		},
		GRPC: GRPCConfig{
			Port: utils.GetEnvOrDefault("GRPC_PORT", "50051"),
		},
		JWT: JWTConfig{
			Secret:     utils.GetEnvOrDefault("JWT_SECRET", "your-super-secret-jwt-key-here"),
			Issuer:     utils.GetEnvOrDefault("JWT_ISSUER", "microservice-b"),
			Expiration: time.Duration(utils.ParseInt(utils.GetEnvOrDefault("JWT_EXPIRATION_SECONDS", "86400"))) * time.Second,
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: utils.ParseInt(utils.GetEnvOrDefault("RATE_LIMIT", "100")),
		},
		Cache: CacheConfig{
			TTL: utils.ParseDurationOrZero(utils.GetEnvOrDefault("CACHE_TTL_SECONDS", "300s")),
		},
	}
}

// GetDSN returns database connection string
func (c *Config) GetDSN() string {
	return c.Database.Username + ":" + c.Database.Password + "@tcp(" + c.Database.Host + ":" + c.Database.Port + ")/" + c.Database.Name + "?charset=utf8mb4&parseTime=True&loc=Local"
}

// GetRedisAddr returns Redis connection address
func (c *Config) GetRedisAddr() string {
	return c.Redis.Host + ":" + c.Redis.Port
}
