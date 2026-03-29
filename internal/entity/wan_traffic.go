package entity

import "time"

type WanTraffic struct {
	ID                 uint   `gorm:"primaryKey"`
	WanID              string `gorm:"index;column:wan_id"`
	Hostname           string
	AgentHost          string
	InterfaceName      string
	RxMbps             float64
	TxMbps             float64
	CapacityMbps       float64
	UtilizationPercent float64
	CreatedAt          time.Time
}
