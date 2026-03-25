package main

import (
	"wan-system/internal/config"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// @title WAN Monitoring System API
// @version 1.0.0
// @description API for monitoring WAN bandwidth and managing user authentication.
// @host localhost:9090
// @BasePath /
// @securityDefinitions.apikey bearerAuth
// @in header
// @name Authorization
// @description Type "Bearer <your_token>" to authenticate

func main() {
	log := config.NewLogger()
	viper := config.NewViper()
	db := config.NewDatabase(viper, log)
	app := fiber.New()
	validate := validator.New()
	secretKey := config.SecretKey(viper)

	config.Bootstrap(&config.BootstrapConfig{
		DB:        db,
		App:       app,
		Log:       log,
		Config:    viper,
		Validate:  validate,
		SecretKey: secretKey,
	})

	err := app.Listen(":9090")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
