package gateway

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	pb "github.com/zakirkun/banking-microservices/api/proto"
)

func (s *Server) authMiddleware(c *fiber.Ctx) error {
	// Get token from header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing authorization header")
	}

	// Check token format
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization header format")
	}

	token := parts[1]

	// Validate token
	resp, err := s.authClient.ValidateToken(context.Background(), &pb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	// Set customer ID in context
	c.Locals("customerID", resp.UserId)
	return c.Next()
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
