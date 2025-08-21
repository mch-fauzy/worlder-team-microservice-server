package interfaces

import (
	"time"

	"github.com/worlder-team/microservice-server/microservice-b/modules/auth/entities"
)

// TokenResponse represents JWT token with expiration information
type TokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// JWTServiceInterface defines the interface for JWT service
type JWTServiceInterface interface {
	GenerateToken(user *entities.User) (*TokenResponse, error)
	ValidateToken(tokenString string) (*entities.User, error)
}
