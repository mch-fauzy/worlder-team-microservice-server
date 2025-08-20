package services

import (
	"github.com/worlder-team/microservice-server/microservice-b/modules/auth/entities"
)

type JWTServiceInterface interface {
	GenerateToken(user *entities.User) (string, error)
	ValidateToken(tokenString string) (*entities.User, error)
}
