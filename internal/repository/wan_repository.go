package repository

import (
	"wan-system/internal/entity"
	"wan-system/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WanCapacityRepository struct {
	Repository[entity.WanCapacity]
	Log *logrus.Logger
}

func NewWanCapacityRepository(log *logrus.Logger) *WanCapacityRepository {
	return &WanCapacityRepository{
		Log: log,
	}
}

func (r *WanCapacityRepository) List(db *gorm.DB) ([]entity.WanCapacity, error) {
	var entities []entity.WanCapacity
	if err := db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

type WanTrafficRepository struct {
	Repository[entity.WanTraffic]
	Log *logrus.Logger
}

func NewWanTrafficRepository(log *logrus.Logger) *WanTrafficRepository {
	return &WanTrafficRepository{
		Log: log,
	}
}

func (r *WanTrafficRepository) Search(db *gorm.DB, request *model.SearchTrafficRequest) ([]entity.WanTraffic, int64, error) {
	var entities []entity.WanTraffic

	if request.Page <= 0 {
		request.Page = 1
	}
	if request.Size <= 0 {
		request.Size = 10
	}

	query := db.Model(&entity.WanTraffic{})
	if request.WanID != "" {
		query = query.Where("wan_id LIKE ?", "%"+request.WanID+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((request.Page - 1) * request.Size).Limit(request.Size).Order("utilization_percent desc, created_at desc").Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}
