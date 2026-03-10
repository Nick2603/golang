package middlewares

import (
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Nick2603/golang/final_project/internal/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-jwt-secret"
const testUserID = "507f1f77bcf86cd799439011"

func setupAuthApp(t *testing.T) *fiber.App {
	t.Helper()

	app := fiber.New()
	app.Use(Auth(testSecret))
	app.Get("/test", func(c fiber.Ctx) error {
		userID, _ := c.Locals("userID").(string)
		return c.SendString(userID)
	})
	return app
}

func TestAuth(t *testing.T) {
	validToken, err := utils.GenerateToken(testUserID, testSecret, time.Hour)
	require.NoError(t, err)

	expiredToken, err := utils.GenerateToken(testUserID, testSecret, -time.Second)
	require.NoError(t, err)

	wrongSecretToken, err := utils.GenerateToken(testUserID, "wrong-secret", time.Hour)
	require.NoError(t, err)

	tests := []struct {
		name       string
		authHeader string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "passes valid bearer token and sets userID in locals",
			authHeader: "Bearer " + validToken,
			wantStatus: fiber.StatusOK,
			wantBody:   testUserID,
		},
		{
			name:       "returns 401 when authorization header is missing",
			authHeader: "",
			wantStatus: fiber.StatusUnauthorized,
		},
		{
			name:       "returns 401 when token is expired",
			authHeader: "Bearer " + expiredToken,
			wantStatus: fiber.StatusUnauthorized,
		},
		{
			name:       "returns 401 when token has wrong secret",
			authHeader: "Bearer " + wrongSecretToken,
			wantStatus: fiber.StatusUnauthorized,
		},
		{
			name:       "returns 401 when authorization format is invalid",
			authHeader: "Token " + validToken,
			wantStatus: fiber.StatusUnauthorized,
		},
		{
			name:       "returns 401 when token string is malformed",
			authHeader: "Bearer not.a.valid.token",
			wantStatus: fiber.StatusUnauthorized,
		},
		{
			name:       "is case-insensitive for bearer prefix",
			authHeader: "bearer " + validToken,
			wantStatus: fiber.StatusOK,
			wantBody:   testUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupAuthApp(t)

			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantBody != "" {
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				assert.Equal(t, tt.wantBody, string(body))
			}
		})
	}
}
