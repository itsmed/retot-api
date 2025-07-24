package middleware

// UserMiddleware provides user-related middleware functions
import (
	"github.com/gofiber/fiber/v2"
)

// GetUserID retrieves the user ID from the request context
func GetUserID(c *fiber.Ctx) (uint, error) {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	return userID, nil
}
