package config

import (
	"wan-system/internal/delivery/http"
	"wan-system/internal/delivery/http/middleware"
	"wan-system/internal/delivery/http/route"
	"wan-system/internal/delivery/nats"
	"wan-system/internal/repository"
	"wan-system/internal/usecase"
	"wan-system/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB        *gorm.DB
	App       *fiber.App
	Log       *logrus.Logger
	Config    *viper.Viper
	Validate  *validator.Validate
	SecretKey string
}

func Bootstrap(config *BootstrapConfig) {
	// CORS
	config.App.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,OPTIONS,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Repositories
	wanCapacityRepository := repository.NewWanCapacityRepository(config.Log)
	wanTrafficRepository := repository.NewWanTrafficRepository(config.Log)
	userRepository := repository.NewUserRepository(config.Log)
	telegrafRepository := repository.NewTelegrafRepository()

	// Token util
	tokenUtil := util.NewTokenUtil(config.SecretKey)

	// Use cases
	wanUseCase := usecase.NewWanUseCase(config.DB, config.Log, wanCapacityRepository, wanTrafficRepository)
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, tokenUtil)
	telegrafUseCase := usecase.NewTelegrafUseCase(config.DB, config.Log, telegrafRepository)

	// Controllers
	wanController := http.NewWanController(config.Log, wanUseCase)
	userController := http.NewUserController(config.Log, userUseCase)
	telegrafController := http.NewTelegrafController(config.Log, telegrafUseCase)

	// Middleware
	authMiddleware := middleware.NewAuth(tokenUtil, config.Log)

	// Consumer
	wanConsumer := nats.NewWanConsumer(config.Log, config.Config, wanUseCase)
	go wanConsumer.Start()

	routeConfig := route.RouteConfig{
		App:            config.App,
		WanController:      wanController,
		UserController:     userController,
		TelegrafController: telegrafController,
		AuthMiddleware:     authMiddleware,
	}
	routeConfig.Setup()
}
