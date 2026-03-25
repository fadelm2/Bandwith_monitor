package repository

import (
	"wan-system/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// UserRepository handles database operations for User entity
type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{Log: log}
}

// FindByToken fetches a user by their JWT token column
func (r *UserRepository) FindByToken(db *gorm.DB, user *entity.User, token string) error {
	return db.Where("token = ?", token).First(user).Error
}
