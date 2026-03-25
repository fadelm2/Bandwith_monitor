package route

import (
	"wan-system/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	WanController  *http.WanController
	UserController *http.UserController
	AuthMiddleware fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.App.Get("/internal/health", c.WanController.Health)

	// Protected Internal Routes (Dashboard)
	internal := c.App.Group("/internal", c.AuthMiddleware)

	// Capacity CRUD
	internal.Post("/capacity", c.WanController.CreateCapacity)
	internal.Put("/capacity/bulk", c.WanController.BulkUpdateCapacity)
	internal.Get("/capacity", c.WanController.ListCapacity)
	internal.Get("/capacity/:wanId", c.WanController.GetCapacity)
	internal.Put("/capacity/:wanId", c.WanController.UpdateCapacity)
	internal.Delete("/capacity/:wanId", c.WanController.DeleteCapacity)

	// Traffic
	internal.Get("/traffic", c.WanController.SearchTraffic)

	// Auth — public
	c.App.Post("/api/auth/register", c.UserController.Register)
	c.App.Post("/api/auth/login", c.UserController.Login)

	// Auth — protected (requires valid JWT)
	c.App.Get("/api/auth/current", c.AuthMiddleware, c.UserController.Current)
	c.App.Post("/api/auth/logout", c.AuthMiddleware, c.UserController.Logout)
}
