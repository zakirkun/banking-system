package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(authService Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		// Remove "Bearer " prefix if present
		token := strings.TrimPrefix(authHeader, "Bearer ")

		resp, err := authService.ValidateToken(token)
		if err != nil || !resp.IsValid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Add user ID to context for downstream handlers
		c.Locals("user_id", resp.UserID)
		return c.Next()
	}
}
