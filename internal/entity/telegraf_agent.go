package entity

import "time"

type TelegrafAgent struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement"`
	IPAddress   string    `gorm:"column:ip_address"`
	Port        int       `gorm:"column:port"`
	Protocol    string    `gorm:"column:protocol"`
	Description string    `gorm:"column:description"`
	IsActive    bool      `gorm:"column:is_active"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (a *TelegrafAgent) TableName() string {
	return "telegraf_agents"
}
