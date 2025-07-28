package handler_test

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"app/database"
	"app/handler"
	"app/middleware"
	"app/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func authHeader() string {
	claims := jwt.MapClaims{
		"user_id": 1, // must match key used in middleware.GetUserID
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte("testsecret"))

	return "Bearer " + signed
}

func truncateTables() {
	database.DB.Exec("DELETE FROM reviews")
	database.DB.Exec("DELETE FROM items")
	database.DB.Exec("DELETE FROM categories")
	database.DB.Exec("DELETE FROM users")
}

func setupTestDB() {
	database.ConnectDBWithDSN(":memory:")
	database.DB.AutoMigrate(&model.User{}, &model.Category{}, &model.Review{}, &model.Item{})
	truncateTables()

	// Create test user
	database.DB.Create(&model.User{
		ID:       1,
		Username: "testuser1",
		Email:    "test1@example.com",
		Password: "hashedpass",
	})

	// Create test category
	category := model.Category{Name: "Default"}
	database.DB.Create(&category)

	// Create test item tied to user
	database.DB.Create(&model.Item{
		Name:        "Sample Item",
		Description: "This is a sample item",
		Price:       100.0,
		UserID:      1,
		Reviews:     []model.Review{{UserID: 1, Rating: 5, Comment: "Great item!"}},
	})
}

func setupProtectedItemApp() *fiber.App {
	setupTestDB()
	os.Setenv("SECRET", "testsecret") // used by middleware

	app := fiber.New()

	// Apply your real middleware to the /api/items group
	item := app.Group("/api/items", middleware.Protected())
	item.Post("/", handler.CreateItem)
	item.Patch("/:id", handler.UpdateItem)
	item.Delete("/:id", handler.DeleteItem)

	return app
}

func TestCreateItem(t *testing.T) {
	app := setupProtectedItemApp()
	category := model.Category{Name: "Electronics"}
	database.DB.Create(&category)

	body := `{"name":"Laptop","description":"Gaming laptop","price":1299.99,"category_id":` + strconv.Itoa(int(category.ID)) + `}`

	req := httptest.NewRequest("POST", "/api/items/", bytes.NewReader([]byte(body)))
	req.Header.Set("Authorization", authHeader())
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestUpdateItem(t *testing.T) {
	app := setupProtectedItemApp()

	category := model.Category{Name: "Books"}
	database.DB.Create(&category)

	item := model.Item{
		Name:        "Go Book",
		Description: "Learn Go",
		Price:       20,
		CategoryID:  category.ID,
		UserID:      1, // assuming user ID 1 is the test user
	}
	database.DB.Create(&item)

	body := `{"name":"Updated Go Book","description":"Master Go","price":25}`
	req := httptest.NewRequest("PATCH", "/api/items/"+strconv.Itoa(int(item.ID)), bytes.NewReader([]byte(body)))
	req.Header.Set("Authorization", authHeader())
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestDeleteItem(t *testing.T) {
	app := setupProtectedItemApp()

	category := model.Category{Name: "Games"}
	database.DB.Create(&category)

	item := model.Item{
		Name:        "PS5",
		Description: "Console",
		Price:       500,
		UserID:      1, // assuming user ID 1 is the test user
		CategoryID:  category.ID,
	}
	database.DB.Create(&item)

	req := httptest.NewRequest("DELETE", "/api/items/"+strconv.Itoa(int(item.ID)), nil)
	req.Header.Set("Authorization", authHeader())

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestUpdateItem_Unauthorized(t *testing.T) {
	app := setupProtectedItemApp()

	// Create item owned by user ID 2
	user := model.User{ID: 2, Username: "other", Email: "other@example.com", Password: "x"}
	database.DB.Create(&user)

	cat := model.Category{Name: "Test"}
	database.DB.Create(&cat)

	item := model.Item{
		Name:        "Not Yours",
		Description: "Secret",
		Price:       10,
		UserID:      2,
	}
	database.DB.Create(&item)

	body := `{"name":"Hack","description":"Hack","price":999}`
	req := httptest.NewRequest("PATCH", "/api/items/"+strconv.Itoa(int(item.ID)), bytes.NewReader([]byte(body)))
	req.Header.Set("Authorization", authHeader()) // logged in as user ID 1
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)
}

func TestDeleteItem_Unauthorized(t *testing.T) {
	app := setupProtectedItemApp()

	// Create item owned by user ID 2
	user := model.User{ID: 2, Username: "other", Email: "other@example.com", Password: "x"}
	database.DB.Create(&user)

	cat := model.Category{Name: "Test"}
	database.DB.Create(&cat)

	item := model.Item{
		Name:        "Stolen",
		Description: "Not yours",
		Price:       99,
		UserID:      2,
	}
	database.DB.Create(&item)

	req := httptest.NewRequest("DELETE", "/api/items/"+strconv.Itoa(int(item.ID)), nil)
	req.Header.Set("Authorization", authHeader()) // user ID 1

	resp, err := app.Test(req)
	fmt.Println("Response status:", resp.StatusCode)
	fmt.Println("Error:", err)
	assert.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)
}
