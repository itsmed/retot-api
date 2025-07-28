// handler/user_test.go
package handler_test

import (
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"app/database"
	"app/handler"
	"app/model"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	database.ConnectDBWithDSN(":memory:")  // uses in-memory SQLite
	database.DB.AutoMigrate(&model.User{}) // create user table
	database.DB.Create(&model.User{
		ID:       1,
		Username: fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Email:    fmt.Sprintf("testuser_%d@example.com", time.Now().UnixNano()),
		Password: "hashedpass",
	})

	app := fiber.New()
	app.Get("/user/:id", handler.GetUser)
	return app
}

func TestGetUser_Success(t *testing.T) {
	app := setupTestApp()

	// Seed user
	user := model.User{
		Username: fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Email:    fmt.Sprintf("testuser_%d@example.com", time.Now().UnixNano()),
		Password: "secretpass",
	}
	database.DB.Create(&user)

	// Send request
	req := httptest.NewRequest("GET", "/user/"+strconv.Itoa(int(user.ID)), nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetUser_NotFound(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/user/9999", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}
