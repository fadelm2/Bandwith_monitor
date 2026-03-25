package middleware

import (
	"wan-system/internal/model"
	"wan-system/internal/util"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// NewAuth returns a Fiber middleware that validates JWT from Authorization header or cookie
func NewAuth(tokenUtil *util.TokenUtil, log *logrus.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Try Authorization header first, fall back to cookie
		tokenStr := ctx.Get("Authorization", "")
		if tokenStr == "" {
			tokenStr = ctx.Cookies("token")
		}

		if tokenStr == "" {
			log.Warn("Missing token in request")
			return fiber.ErrUnauthorized
		}

		auth, err := tokenUtil.ParseToken(ctx.UserContext(), tokenStr)
		if err != nil {
			log.Warnf("Invalid token: %+v", err)
			return fiber.ErrUnauthorized
		}

		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

// GetUser retrieves the authenticated user from fiber context locals
func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
