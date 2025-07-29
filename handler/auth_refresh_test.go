package handler_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"app/database"
	"app/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func setupRefreshApp() *fiber.App {
	database.ConnectDBWithDSN(":memory:")
	os.Setenv("SECRET", "testsecret")
	os.Setenv("REFRESH_SECRET", "refreshsecret")

	app := fiber.New()
	app.Get("/auth/refresh", handler.RefreshToken)
	return app
}

func createRefreshCookie(userID uint, username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	signed, _ := token.SignedString([]byte("refreshsecret"))
	return signed
}

func TestRefreshToken_Success(t *testing.T) {
	app := setupRefreshApp()
	req := httptest.NewRequest("GET", "/auth/refresh", nil)
	cookie := createRefreshCookie(1, "testuser")
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: cookie})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestRefreshToken_Invalid(t *testing.T) {
	app := setupRefreshApp()
	req := httptest.NewRequest("GET", "/auth/refresh", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "invalidtoken"})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}
