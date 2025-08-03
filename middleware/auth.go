package middleware

import (
	"github.com/gofiber/fiber/v2"
	u "github.com/strbagus/fiber-auth/utils"
)

func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !u.IsJWTUsable(c.Cookies("access_token")) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token is not valid",
			})
		}
		return c.Next()
	}
}
