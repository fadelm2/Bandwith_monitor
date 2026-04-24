package repository

import (
	"wan-system/internal/entity"

	"gorm.io/gorm"
)

type TelegrafRepository struct {
	Repository[entity.TelegrafAgent]
}

func NewTelegrafRepository() *TelegrafRepository {
	return &TelegrafRepository{}
}

func (r *TelegrafRepository) ListActive(db *gorm.DB) ([]entity.TelegrafAgent, error) {
	var agents []entity.TelegrafAgent
	err := db.Where("is_active = ?", true).Find(&agents).Error
	return agents, err
}
