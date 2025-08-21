package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/worlder-team/microservice-server/microservice-b/modules/auth/entities"
	"github.com/worlder-team/microservice-server/microservice-b/modules/auth/interfaces"
)

type jwtService struct {
	secretKey  string
	issuer     string
	expiration time.Duration
}

type jwtClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTService(secretKey, issuer string, expiration time.Duration) interfaces.JWTServiceInterface {
	return &jwtService{
		secretKey:  secretKey,
		issuer:     issuer,
		expiration: expiration,
	}
}

func (s *jwtService) GenerateToken(user *entities.User) (*interfaces.TokenResponse, error) {
	expirationTime := time.Now().Add(s.expiration)

	claims := &jwtClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
			Subject:   user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return nil, err
	}

	return &interfaces.TokenResponse{
		Token:     tokenString,
		ExpiresAt: expirationTime,
	}, nil
}

func (s *jwtService) ValidateToken(tokenString string) (*entities.User, error) {
	claims := &jwtClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	user := &entities.User{
		ID:   claims.UserID,
		Role: claims.Role,
	}

	return user, nil
}
