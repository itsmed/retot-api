package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// Logout handles user logout
// Logout clears the JWT cookie
func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set expired date
		HTTPOnly: true,
		Secure:   false, // true in production
	})
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Logged out successfully",
	})
}
