package interfaces

import (
	"context"

	"github.com/worlder-team/microservice-server/microservice-b/modules/auth/dtos"
	"github.com/worlder-team/microservice-server/microservice-b/modules/auth/entities"
)

// AuthServiceInterface defines the interface for authentication service
type AuthServiceInterface interface {
	Login(ctx context.Context, request *dtos.LoginRequest) (*dtos.LoginResponse, error)
	ValidateToken(ctx context.Context, token string) (*entities.User, error)
	GenerateToken(user *entities.User) (string, error)
}
