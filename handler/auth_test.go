package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"app/database"
	"app/handler"
	"app/model"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type RegisterPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginPayload struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

func setupAuthApp() *fiber.App {
	database.ConnectDBWithDSN(":memory:")
	database.DB.AutoMigrate(&model.User{})
	os.Setenv("SECRET", "testsecret")

	app := fiber.New()
	app.Post("/register", handler.Register)
	app.Post("/login", handler.Login)
	return app
}

func TestRegister_Success(t *testing.T) {
	app := setupAuthApp()
	payload := RegisterPayload{"testuser", "test@example.com", "securepass"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestRegister_Conflict(t *testing.T) {
	app := setupAuthApp()
	// Create user
	database.DB.Create(&model.User{Username: "dupuser", Email: "dup@example.com", Password: "x"})
	payload := RegisterPayload{"dupuser", "dup@example.com", "securepass"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 409, resp.StatusCode)
}

func TestLogin_Success(t *testing.T) {
	app := setupAuthApp()
	// Register first
	hash, _ := handler.HashPassword("securepass")
	database.DB.Create(&model.User{Username: "tester", Email: "tester@example.com", Password: hash})

	payload := LoginPayload{"tester", "securepass"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestLogin_InvalidPassword(t *testing.T) {
	app := setupAuthApp()
	hash, _ := handler.HashPassword("securepass")
	database.DB.Create(&model.User{Username: "wrongpass", Email: "wp@example.com", Password: hash})

	payload := LoginPayload{"wrongpass", "wrongpass123"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}
