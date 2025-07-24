package main

import (
	"log"

	"app/database"
	"app/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "App Name",
	})
	app.Options("/api/*", func(c *fiber.Ctx) error {
    c.Set("Access-Control-Allow-Origin", "http://localhost:3000")
    c.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    c.Set("Access-Control-Allow-Credentials", "true")
    return c.SendStatus(fiber.StatusOK)
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000", // Your Next.js frontend URL
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true, // Allow cookies (like JWT tokens)
	}))

	database.ConnectDB()

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
