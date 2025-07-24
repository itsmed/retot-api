package main

import (
	"app/database"
	"app/router"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "ReTot",
	})

	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3001, http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	app.Options("/api/*", func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "http://localhost:3001, http://localhost:3000")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization")
		c.Set("Access-Control-Allow-Credentials", "true")
		return c.SendStatus(fiber.StatusNoContent)
	})

	database.ConnectDB()

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
