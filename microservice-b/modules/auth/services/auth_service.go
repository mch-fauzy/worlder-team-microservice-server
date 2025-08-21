package services

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/worlder-team/microservice-server/microservice-b/modules/auth/dtos"
	"github.com/worlder-team/microservice-server/microservice-b/modules/auth/entities"
	"github.com/worlder-team/microservice-server/microservice-b/modules/auth/interfaces"
)

type authService struct {
	db         *gorm.DB
	jwtService interfaces.JWTServiceInterface
}

func NewAuthService(db *gorm.DB, jwtService interfaces.JWTServiceInterface) interfaces.AuthServiceInterface {
	return &authService{
		db:         db,
		jwtService: jwtService,
	}
}

// Login authenticates user with email and password, returns JWT token
func (s *authService) Login(ctx context.Context, request *dtos.LoginRequest) (*dtos.LoginResponse, error) {
	// Find user by email
	var user entities.User
	err := s.db.Where("email = ?", request.Email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Validate password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	tokenResponse, err := s.jwtService.GenerateToken(&user)
	if err != nil {
		return nil, err
	}

	// Create user response (without password)
	userResponse := dtos.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	response := &dtos.LoginResponse{
		Token:     tokenResponse.Token,
		ExpiresAt: tokenResponse.ExpiresAt,
		User:      userResponse,
	}

	return response, nil
}

// HashPassword creates a bcrypt hash of the password
func (s *authService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CreateUser creates a new user with hashed password
func (s *authService) CreateUser(username, email, password, role string) (*entities.User, error) {
	// Check if user already exists
	var existingUser entities.User
	err := s.db.Where("username = ? OR email = ?", username, email).First(&existingUser).Error
	if err == nil {
		return nil, errors.New("user already exists")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Hash password
	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := entities.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Role:     role,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
