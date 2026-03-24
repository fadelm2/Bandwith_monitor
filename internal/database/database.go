package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"wan-system/internal/models"
)

var DB *gorm.DB

func Init() {

	dsn := "root:password@tcp(127.0.0.1:3306)/network?parseTime=true"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DB = db

	DB.AutoMigrate(
		&models.WanTraffic{},
		&models.WanCapacity{},
	)
}
