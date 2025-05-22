package store

import (
	"context"

	"github.com/ZhuoyangM/ConfigLeak/internal/models"
	"gorm.io/gorm"
)

type UserService struct {
	db  *gorm.DB
	ctx context.Context
}

func (service *UserService) CreateUser(user *models.User) error {
	return service.db.WithContext(service.ctx).Create(user).Error
}

func (service *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	return nil, nil
}
