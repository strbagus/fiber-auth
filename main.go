package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/strbagus/fiber-auth/database"
	"github.com/strbagus/fiber-auth/routes"
)

func main() {
	database.DBConnect()
	database.RedisConnect()
	app := fiber.New(fiber.Config{})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, UPDATE, DELETE",
		AllowCredentials: true,
	}))

	routes.RegisterRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Service up and running!")
	})

	log.Fatal(app.Listen(":3000"))
}
