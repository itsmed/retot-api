package handler_test

import (
	"net/http/httptest"
	"os"
	"testing"

	"app/database"
	"app/handler"
	"app/model"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupApiApp() *fiber.App {
	database.ConnectDBWithDSN(":memory:")
	database.DB.AutoMigrate(&model.User{})
	os.Setenv("SECRET", "testsecret")

	app := fiber.New()
	app.Get("/", handler.Hello)
	return app
}

func TestHello(t *testing.T) {
	app := setupApiApp()

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
