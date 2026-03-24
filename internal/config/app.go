package config

import (
	"wan-system/internal/delivery/http"
	"wan-system/internal/delivery/http/route"
	"wan-system/internal/delivery/nats"
	"wan-system/internal/repository"
	"wan-system/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB     *gorm.DB
	App    *fiber.App
	Log    *logrus.Logger
	Config *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// CORS
	config.App.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,OPTIONS,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// setup repositories
	wanCapacityRepository := repository.NewWanCapacityRepository(config.Log)
	wanTrafficRepository := repository.NewWanTrafficRepository(config.Log)

	// setup use cases
	wanUseCase := usecase.NewWanUseCase(config.DB, config.Log, wanCapacityRepository, wanTrafficRepository)

	// setup controller
	wanController := http.NewWanController(config.Log, wanUseCase)

	// setup consumer
	wanConsumer := nats.NewWanConsumer(config.Log, config.Config, wanUseCase)
	go wanConsumer.Start()

	routeConfig := route.RouteConfig{
		App:           config.App,
		WanController: wanController,
	}
	routeConfig.Setup()
}
