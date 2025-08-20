package constants

const (
	// Default pagination
	DefaultPageSize = 10
	MaxPageSize     = 100

	// Sensor types
	SensorTypeTemperature = "temperature"
	SensorTypeHumidity    = "humidity"
	SensorTypePressure    = "pressure"
	SensorTypeLight       = "light"
	SensorTypeMotion      = "motion"
)

// Error messages
const (
	ErrInvalidRequest    = "invalid request"
	ErrUnauthorized      = "unauthorized"
	ErrForbidden         = "forbidden"
	ErrNotFound          = "not found"
	ErrInternalServer    = "internal server error"
	ErrValidationFailed  = "validation failed"
	ErrDuplicateEntry    = "duplicate entry"
	ErrRateLimitExceeded = "rate limit exceeded"
)

// HTTP Status messages
const (
	StatusSuccess = "success"
	StatusError   = "error"
	StatusFailed  = "failed"
)
