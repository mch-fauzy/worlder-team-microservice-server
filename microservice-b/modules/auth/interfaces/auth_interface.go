package interfaces

import (
	"context"

	"github.com/worlder-team/microservice-server/microservice-b/modules/auth/dtos"
)

// AuthServiceInterface defines the interface for authentication service
type AuthServiceInterface interface {
	Login(ctx context.Context, request *dtos.LoginRequest) (*dtos.LoginResponse, error)
}
