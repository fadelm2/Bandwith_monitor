package main

import (
	"wan-system/internal/config"

	"github.com/gofiber/fiber/v2"
)

func main() {
	log := config.NewLogger()
	viper := config.NewViper()
	db := config.NewDatabase(viper, log)
	app := fiber.New()

	config.Bootstrap(&config.BootstrapConfig{
		DB:     db,
		App:    app,
		Log:    log,
		Config: viper,
	})

	err := app.Listen(":9090")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
