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

// TODO: configure the context
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db:  db,
		ctx: context.Background(),
	}
}

func (service *UserService) CreateUser(req *RegisterRequest) error {
	user, err := ToUser(req)
	if err != nil {
		return err
	}
	return service.db.WithContext(service.ctx).Create(user).Error
}

func (service *UserService) AuthenticateUser(req *LoginRequest) (string, error) {
	// Find user and check password
	var user User
	username := html.EscapeString(req.Username)
	password := html.EscapeString(req.Password)
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

func (service *UserService) GetUserByID(userID uint) (*GetUserResponse, error) {
	var user User
	if err := service.db.WithContext(service.ctx).First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	resp := ToGetUserResponse(&user)
	return resp, nil
}
