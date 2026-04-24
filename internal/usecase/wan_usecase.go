package usecase

import (
	"context"
	"time"
	"wan-system/internal/entity"
	"wan-system/internal/model"
	"wan-system/internal/model/converter"
	"wan-system/internal/repository"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WanUseCase struct {
	DB                    *gorm.DB
	Log                   *logrus.Logger
	WanCapacityRepository *repository.WanCapacityRepository
	WanTrafficRepository  *repository.WanTrafficRepository
}

func NewWanUseCase(db *gorm.DB, logger *logrus.Logger, capRepo *repository.WanCapacityRepository, trafficRepo *repository.WanTrafficRepository) *WanUseCase {
	return &WanUseCase{
		DB:                    db,
		Log:                   logger,
		WanCapacityRepository: capRepo,
		WanTrafficRepository:  trafficRepo,
	}
}

func (c *WanUseCase) CreateCapacity(ctx context.Context, request *model.WanCapacityRequest) (*model.WanCapacityResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	entity := &entity.WanCapacity{
		WanID:            request.WanID,
		CapacityMbps:     request.CapacityMbps,
		ThresholdPercent: request.ThresholdPercent,
		Description:      request.Description,
		CreatedAt:        time.Now(),
	}

	if err := c.WanCapacityRepository.Create(tx, entity); err != nil {
		c.Log.Warnf("Failed to create capacity: %v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return converter.WanCapacityToResponse(entity), nil
}

func (c *WanUseCase) UpdateCapacity(ctx context.Context, request *model.WanCapacityRequest) (*model.WanCapacityResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var entity entity.WanCapacity
	if err := c.WanCapacityRepository.FindById(tx, &entity, request.WanID); err != nil {
		c.Log.Warnf("Failed to find capacity for update: %v", err)
		return nil, err
	}

	entity.CapacityMbps = request.CapacityMbps
	entity.ThresholdPercent = request.ThresholdPercent
	entity.Description = request.Description

	if err := c.WanCapacityRepository.Update(tx, &entity); err != nil {
		c.Log.Warnf("Failed to update capacity: %v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return converter.WanCapacityToResponse(&entity), nil
}

func (c *WanUseCase) DeleteCapacity(ctx context.Context, wanID string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var entity entity.WanCapacity
	if err := c.WanCapacityRepository.FindById(tx, &entity, wanID); err != nil {
		c.Log.Warnf("Failed to find capacity for delete: %v", err)
		return err
	}

	if err := c.WanCapacityRepository.Delete(tx, &entity); err != nil {
		return err
	}

	return tx.Commit().Error
}

func (c *WanUseCase) GetCapacity(ctx context.Context, wanID string) (*model.WanCapacityResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var entity entity.WanCapacity
	if err := c.WanCapacityRepository.FindById(tx, &entity, wanID); err != nil {
		c.Log.Warnf("Failed to find capacity: %v", err)
		return nil, err
	}

	return converter.WanCapacityToResponse(&entity), nil
}

func (c *WanUseCase) ListCapacity(ctx context.Context) ([]*model.WanCapacityResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	entities, err := c.WanCapacityRepository.List(tx)
	if err != nil {
		return nil, err
	}

	var responses []*model.WanCapacityResponse
	for _, e := range entities {
		responses = append(responses, converter.WanCapacityToResponse(&e))
	}

	return responses, nil
}

func (c *WanUseCase) BulkUpdateCapacity(ctx context.Context, request []model.WanCapacityRequest) ([]*model.WanCapacityResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var entities []entity.WanCapacity
	for _, req := range request {
		entities = append(entities, entity.WanCapacity{
			WanID:            req.WanID,
			CapacityMbps:     req.CapacityMbps,
			ThresholdPercent: req.ThresholdPercent,
			CreatedAt:        time.Now(),
		})
	}

	if err := tx.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&entities).Error; err != nil {
		c.Log.Warnf("Failed to bulk update capacity: %v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	var responses []*model.WanCapacityResponse
	for _, e := range entities {
		responses = append(responses, converter.WanCapacityToResponse(&e))
	}

	return responses, nil
}

func (c *WanUseCase) ProcessTraffic(ctx context.Context, traffic *entity.WanTraffic) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Get capacity to calculate utilization
	var cap entity.WanCapacity
	err := c.WanCapacityRepository.FindById(tx, &cap, traffic.WanID)
	if err == nil {
		traffic.CapacityMbps = cap.CapacityMbps
		if cap.CapacityMbps > 0 {
			traffic.UtilizationPercent = (traffic.RxMbps / cap.CapacityMbps) * 100
		}
	}

	if err := c.WanTrafficRepository.Create(tx, traffic); err != nil {
		c.Log.Warnf("Failed to record traffic: %v", err)
		return err
	}

	return tx.Commit().Error
}

func (c *WanUseCase) SearchTraffic(ctx context.Context, request *model.SearchTrafficRequest) ([]*model.WanTrafficResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	entities, total, err := c.WanTrafficRepository.Search(tx, request)
	if err != nil {
		c.Log.Warnf("Failed to search traffic: %v", err)
		return nil, 0, err
	}

	var responses []*model.WanTrafficResponse
	for _, e := range entities {
		responses = append(responses, converter.WanTrafficToResponse(&e))
	}

	return responses, total, nil
}
