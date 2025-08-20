package configs

import (
	"time"

	"github.com/worlder-team/microservice-server/shared/utils"
)

// Config holds all configuration for microservice-a
type Config struct {
	Server    ServerConfig
	GRPC      GRPCConfig
	Generator GeneratorConfig
	RateLimit RateLimitConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         string
	ExternalPort string // Swagger Port (setting in docker-compose.yml)
	LogLevel     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// GRPCConfig holds gRPC client configuration
type GRPCConfig struct {
	ServerHost string
	ServerPort string
}

// GeneratorConfig holds sensor generator configuration
type GeneratorConfig struct {
	SensorType string
	Frequency  string
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: utils.GetEnvOrDefault("MICROSERVICE_A_PORT", "8081"),
			// Swagger Port (setting in docker-compose.yml)
			ExternalPort: utils.GetEnvOrDefault("EXTERNAL_PORT", "8081"),
			LogLevel:     utils.GetEnvOrDefault("LOG_LEVEL", "info"),
			ReadTimeout:  utils.ParseDurationOrZero(utils.GetEnvOrDefault("READ_TIMEOUT", "10s")),
			WriteTimeout: utils.ParseDurationOrZero(utils.GetEnvOrDefault("WRITE_TIMEOUT", "10s")),
		},
		GRPC: GRPCConfig{
			// GRPC_HOST points to microservice-b (storage service) where generators send data, check service name in docker-compose.yml
			ServerHost: utils.GetEnvOrDefault("GRPC_HOST", "microservice-b"),
			ServerPort: utils.GetEnvOrDefault("GRPC_PORT", "50051"),
		},
		Generator: GeneratorConfig{
			// SENSOR_TYPE is set by docker-compose.yml for each service instance (not in .env file)
			SensorType: utils.GetEnvOrDefault("SENSOR_TYPE", "temperature"),
			Frequency:  utils.GetEnvOrDefault("GENERATION_FREQUENCY", "300s"),
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: utils.ParseInt(utils.GetEnvOrDefault("RATE_LIMIT", "100")),
		},
	}
}

// GetGRPCAddress returns gRPC server address
func (c *Config) GetGRPCAddress() string {
	return c.GRPC.ServerHost + ":" + c.GRPC.ServerPort
}
