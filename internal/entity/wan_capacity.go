package entity

import "time"

type WanCapacity struct {
	WanID            string `gorm:"primaryKey;column:wan_id"`
	CapacityMbps     float64
	ThresholdPercent float64
	CreatedAt        time.Time
}
