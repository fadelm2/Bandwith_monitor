package api

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
	"wan-system/internal/database"
	"wan-system/internal/models"
)

func StartInternal() {

	app := fiber.New()

	app.Get("/internal/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Put("/internal/capacity/bulk", func(c *fiber.Ctx) error {

		var body []models.WanCapacity

		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid body",
			})
		}

		database.DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&body)

		return c.JSON(fiber.Map{
			"message": "bulk update success",
			"count":   len(body),
		})
	})

	app.Listen(":9090")
}
