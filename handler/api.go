package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Hello handle api status
func Hello(c *fiber.Ctx) error {
	fmt.Println("Hello I'm ok!")
	return c.JSON(fiber.Map{"status": "success", "message": "Hello i'm ok!", "data": nil})
}
