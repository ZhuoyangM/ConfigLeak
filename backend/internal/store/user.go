package store

import (
	"context"
	"fmt"
	"html"

	"github.com/ZhuoyangM/ConfigLeak/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db  *gorm.DB
	ctx context.Context
}

func (service *UserService) CreateUser(user *User) error {
	// Hash the password before storing it
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPwd)

	// Escape special characters in username and email
	user.Username = html.EscapeString(user.Username)
	user.Email = html.EscapeString(user.Email)

	return service.db.WithContext(service.ctx).Create(user).Error
}

func (service *UserService) AuthenticateUser(username, password string) (string, error) {
	// Find user and check password
	var user User
	if err := service.db.WithContext(service.ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid password: %w", err)
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return token, nil
}
