package entity

// User represents a user entity in the system
type User struct {
	ID        string `gorm:"column:id;primaryKey"`
	Password  string `gorm:"column:password"`
	Name      string `gorm:"column:username"`
	Email     string `gorm:"column:email"`
	Token     string `gorm:"column:token"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
}

func (u *User) TableName() string {
	return "users"
}
