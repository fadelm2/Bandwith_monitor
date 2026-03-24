package models

import "time"

type WanCapacity struct {
	WanID            string    `gorm:"primaryKey"`
	CapacityMbps     float64
	ThresholdPercent float64
	CreatedAt        time.Time
}
