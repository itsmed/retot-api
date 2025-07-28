// middleware/auth_test.go
package middleware_test

import (
	"app/middleware"
	"os"
	"testing"
	"net/http/httptest"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func setupProtectedApp() *fiber.App {
	// Set test secret in env
	os.Setenv("SECRET", "testsecret")

	app := fiber.New()
	app.Use(middleware.Protected())
	app.Get("/secure", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	return app
}

func TestProtected_MissingToken(t *testing.T) {
	app := setupProtectedApp()
	req := httptest.NewRequest("GET", "/secure", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestProtected_InvalidToken(t *testing.T) {
	app := setupProtectedApp()
	req := httptest.NewRequest("GET", "/secure", nil)
	req.Header.Set("Authorization", "Bearer BADTOKEN")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestProtected_ValidToken(t *testing.T) {
	app := setupProtectedApp()

	// Create valid token
	claims := jwt.MapClaims{
		"user_id": 1,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte("testsecret"))

	req := httptest.NewRequest("GET", "/secure", nil)
	req.Header.Set("Authorization", "Bearer "+signed)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
