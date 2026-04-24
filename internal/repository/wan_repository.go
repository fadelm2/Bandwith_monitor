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

// FindById overrides the generic Repository.FindById to use wan_id as the primary key
func (r *WanCapacityRepository) FindById(db *gorm.DB, entity *entity.WanCapacity, id any) error {
	return db.Where("wan_id = ?", id).Take(entity).Error
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
	if request.SinceMinutes > 0 {
		query = query.Where("created_at > NOW() - INTERVAL ? MINUTE", request.SinceMinutes)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((request.Page - 1) * request.Size).Limit(request.Size).Order("created_at desc").Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *WanTrafficRepository) GetHighUtilizationAlerts(db *gorm.DB) ([]model.WanAlertResponse, error) {
	var alerts []model.WanAlertResponse

	// Calculate average RX/TX over the last 5 minutes (or latest minute for responsiveness)
	// and join with CURRENT capacity
	err := db.Table("wan_traffics t").
		Select("t.wan_id, c.description, AVG(t.rx_mbps) as avg_rx_mbps, AVG(t.tx_mbps) as avg_tx_mbps, c.capacity_mbps as current_capacity, c.threshold_percent as current_threshold").
		Joins("JOIN wan_capacities c ON t.wan_id = c.wan_id").
		Where("t.created_at > (SELECT MAX(created_at) FROM wan_traffics) - INTERVAL 5 MINUTE").
		Group("t.wan_id, c.description, c.capacity_mbps, c.threshold_percent").
		Having("AVG(t.rx_mbps) / c.capacity_mbps * 100 > c.threshold_percent").
		Scan(&alerts).Error

	if err != nil {
		return nil, err
	}

	for i := range alerts {
		if alerts[i].CurrentCapacity > 0 {
			alerts[i].AvgUtilization = (alerts[i].AvgRxMbps / alerts[i].CurrentCapacity) * 100
		}
	}

	return alerts, nil
}
