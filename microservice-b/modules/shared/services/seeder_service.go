package services

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/worlder-team/microservice-server/microservice-b/modules/auth/entities"
)

type SeederService struct {
	db *gorm.DB
}

func NewSeederService(db *gorm.DB) *SeederService {
	return &SeederService{
		db: db,
	}
}

// SeedDefaultUsers creates default users if they don't exist
func (s *SeederService) SeedDefaultUsers() error {
	defaultUsers := []struct {
		Username string
		Email    string
		Password string
		Role     string
	}{
		{
			Username: "admin",
			Email:    "admin@example.com",
			Password: "password",
			Role:     "admin",
		},
		{
			Username: "testuser",
			Email:    "user@example.com",
			Password: "testuser123",
			Role:     "user",
		},
	}

	for _, userData := range defaultUsers {
		// Check if user already exists
		var existingUser entities.User
		err := s.db.Where("username = ? OR email = ?", userData.Username, userData.Email).First(&existingUser).Error

		if err == gorm.ErrRecordNotFound {
			// User doesn't exist, create it
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("Error hashing password for user %s: %v", userData.Username, err)
				continue
			}

			user := entities.User{
				Username: userData.Username,
				Email:    userData.Email,
				Password: string(hashedPassword),
				Role:     userData.Role,
			}

			if err := s.db.Create(&user).Error; err != nil {
				log.Printf("Error creating user %s: %v", userData.Username, err)
				continue
			}

			log.Printf("Created default user: %s (%s)", userData.Username, userData.Role)
		} else if err != nil {
			log.Printf("Error checking for existing user %s: %v", userData.Username, err)
			continue
		} else {
			log.Printf("User %s already exists, skipping", userData.Username)
		}
	}

	return nil
}

// SeedAll runs all seeding operations
func (s *SeederService) SeedAll() error {
	log.Println("Starting database seeding...")

	if err := s.SeedDefaultUsers(); err != nil {
		return err
	}

	log.Println("Database seeding completed successfully")
	return nil
}
