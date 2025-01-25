package gateway

import (
	"context"

	"github.com/gofiber/fiber/v2"
	pb "github.com/zakirkun/banking-microservices/api/proto"
)

func (s *Server) handleRegister(c *fiber.Ctx) error {
	var req struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FullName  string `json:"full_name"`
		Phone     string `json:"phone"`
		Address   string `json:"address"`
		IDNumber  string `json:"id_number"`
		IDType    string `json:"id_type"`
		BirthDate string `json:"birth_date"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	resp, err := s.authClient.Register(context.Background(), &pb.RegisterRequest{
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.Phone,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (s *Server) handleLogin(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	resp, err := s.authClient.Login(context.Background(), &pb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

// func (s *Server) handleRefreshToken(c *fiber.Ctx) error {
// 	var req struct {
// 		RefreshToken string `json:"refresh_token"`
// 	}

// 	if err := c.BodyParser(&req); err != nil {
// 		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
// 	}

// 	resp, err := s.authClient.RefreshToken(context.Background(), &pb.RefreshTokenRequest{
// 		RefreshToken: req.RefreshToken,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return c.JSON(resp)
// }
