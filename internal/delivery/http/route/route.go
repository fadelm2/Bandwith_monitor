package route

import (
	"wan-system/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App           *fiber.App
	WanController *http.WanController
}

func (c *RouteConfig) Setup() {
	c.App.Get("/internal/health", c.WanController.Health)

	// Capacity CRUD
	c.App.Post("/internal/capacity", c.WanController.CreateCapacity)
	c.App.Put("/internal/capacity/bulk", c.WanController.BulkUpdateCapacity)
	c.App.Get("/internal/capacity", c.WanController.ListCapacity)
	c.App.Get("/internal/capacity/:wanId", c.WanController.GetCapacity)
	c.App.Put("/internal/capacity/:wanId", c.WanController.UpdateCapacity)
	c.App.Delete("/internal/capacity/:wanId", c.WanController.DeleteCapacity)

	// Traffic
	c.App.Get("/internal/traffic", c.WanController.SearchTraffic)
}
