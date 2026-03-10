package middlewares

import (
	"log/slog"
	"strings"

	"github.com/Nick2603/golang/final_project/internal/utils"

	"github.com/gofiber/fiber/v3"
)

func Auth(jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			slog.Warn("Auth middleware: missing Authorization header", "path", c.Path())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "missing authorization header",
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			slog.Warn("Auth middleware: malformed Authorization header", "path", c.Path())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "invalid authorization format, expected: Bearer <token>",
			})
		}

		claims, err := utils.ParseToken(parts[1], jwtSecret)
		if err != nil {
			slog.Warn("Auth middleware: invalid token", "path", c.Path(), "error", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "invalid or expired token",
			})
		}

		c.Locals("userID", claims.UserID)
		slog.Info("Auth middleware: request authorized", "userID", claims.UserID, "path", c.Path())
		return c.Next()
	}
}
