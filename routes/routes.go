package routes

import (
	"github.com/gofiber/fiber/v2"
	h "github.com/strbagus/fiber-auth/handlers"
	m "github.com/strbagus/fiber-auth/middleware"
)

func RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")
	v1.Get("/users", m.RequireAuth(), h.ListUsers)
	v1.Post("/login", h.SignIn)
	v1.Post("/check", h.CheckToken)
	v1.Post("/refresh", h.RefreshToken)
	auth := v1.Group("/auth", m.RequireAuth())
	auth.Post("/invalidate", h.InvalidateToken)

}
