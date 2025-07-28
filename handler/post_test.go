package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"app/database"
	"app/handler"
	"app/middleware"
	"app/model"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// func authHeader() string {
// 	claims := jwt.MapClaims{
// 		"user_id": 1,
// 		"exp":     jwt.NewNumericDate(jwt.TimeFunc().Add(time.Hour)),
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	signed, _ := token.SignedString([]byte("testsecret"))
// 	return "Bearer " + signed
// }

func setupProtectedPostApp() *fiber.App {
	database.ConnectDBWithDSN(":memory:")
	database.DB.AutoMigrate(&model.User{}, &model.Post{})

	database.DB.Create(&model.User{ID: 1, Username: "testuser", Email: "test@example.com", Password: "secret"})

	os.Setenv("SECRET", "testsecret")

	app := fiber.New()
	post := app.Group("/api/posts", middleware.Protected())
	post.Post("/", handler.CreatePost)
	app.Get("/api/users/:id/posts", handler.GetPostFromUser)

	return app
}

func TestCreatePost(t *testing.T) {
	app := setupProtectedPostApp()

	post := map[string]interface{}{
		"title":   "Test Post",
		"content": "This is a test post",
		"user_id": 1,
	}
	body, _ := json.Marshal(post)

	req := httptest.NewRequest("POST", "/api/posts/", bytes.NewReader(body))
	req.Header.Set("Authorization", authHeader())
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetPostFromUser(t *testing.T) {
	app := setupProtectedPostApp()
	database.DB.Create(&model.Post{Title: "Post1", Body: "Content1", UserID: 1})

	req := httptest.NewRequest("GET", "/api/users/1/posts", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetPostFromUser_NotFound(t *testing.T) {
	app := setupProtectedPostApp()

	req := httptest.NewRequest("GET", "/api/users/999/posts", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}
