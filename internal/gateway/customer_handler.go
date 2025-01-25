package gateway

import (
	"context"

	"github.com/gofiber/fiber/v2"
	pb "github.com/zakirkun/banking-microservices/api/proto"
)

func (s *Server) handleGetCustomerProfile(c *fiber.Ctx) error {
	customerID := c.Locals("customerID").(string)

	resp, err := s.customerClient.GetCustomer(context.Background(), &pb.GetCustomerRequest{
		CustomerId: customerID,
	})
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

func (s *Server) handleUpdateCustomerProfile(c *fiber.Ctx) error {
	customerID := c.Locals("customerID").(string)

	var req struct {
		FullName string `json:"full_name"`
		Address  string `json:"address"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	resp, err := s.customerClient.UpdateCustomer(context.Background(), &pb.UpdateCustomerRequest{
		CustomerId: customerID,
		FullName:   req.FullName,
		Address:    req.Address,
	})
	if err != nil {
		return err
	}

	return c.JSON(resp)
}
