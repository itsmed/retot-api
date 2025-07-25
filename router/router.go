package router

import (
	"app/handler"
	"app/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Hello)

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)

	// User
	user := api.Group("/user")
	user.Get("/id/:id", handler.GetUser)
	user.Post("/", handler.CreateUser)
	user.Get("/all", handler.GetAllUsers)
	user.Patch("/id/:id", middleware.Protected(), handler.UpdateUser)
	user.Delete("/id/:id", middleware.Protected(), handler.DeleteUser)

	// Post
	post := api.Group("/post")
	post.Get("/user/:id", handler.GetPostFromUser)
	post.Post("/user/:id/new", middleware.Protected(), handler.CreatePost)
}
