package test

import (
	"wan-system/internal/config"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var app *fiber.App
var db *gorm.DB
var log *logrus.Logger
var validate *validator.Validate
var viperConfig *viper.Viper
var secretKey string

func init() {
	log = config.NewLogger()
	viperConfig = config.NewViper()
	db = config.NewDatabase(viperConfig, log)
	validate = validator.New()
	app = fiber.New()
	secretKey = config.SecretKey(viperConfig)

	config.Bootstrap(&config.BootstrapConfig{
		DB:        db,
		App:       app,
		Log:       log,
		Config:    viperConfig,
		Validate:  validate,
		SecretKey: secretKey,
	})
}
