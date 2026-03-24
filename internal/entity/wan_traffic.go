package entity

import "time"

type WanTraffic struct {
	ID                 uint   `gorm:"primaryKey"`
	WanID              string `gorm:"index"`
	Hostname           string
	InterfaceName      string
	RxMbps             float64
	TxMbps             float64
	CapacityMbps       float64
	UtilizationPercent float64
	CreatedAt          time.Time
}
