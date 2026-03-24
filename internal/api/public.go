package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"wan-system/internal/database"
	"wan-system/internal/models"
)

var jwtSecret = []byte("SUPER_SECRET_KEY")

func StartPublic() {

	app := fiber.New()

	app.Post("/login", func(c *fiber.Ctx) error {

		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

		t, _ := token.SignedString(jwtSecret)

		return c.JSON(fiber.Map{"token": t})
	})

	api := app.Group("/api",
		jwtware.New(jwtware.Config{
			SigningKey: jwtSecret,
		}),
	)

	api.Get("/latest/:wan", func(c *fiber.Ctx) error {

		wan := c.Params("wan")

		var data models.WanTraffic
		database.DB.Where("wan_id = ?", wan).
			Order("created_at DESC").
			First(&data)

		return c.JSON(data)
	})

	app.Listen(":8080")
}
