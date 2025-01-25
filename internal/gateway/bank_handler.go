package gateway

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	pb "github.com/zakirkun/banking-microservices/api/proto"
)

func (s *Server) handleListBanks(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")
	country := c.Query("country", "")
	status := c.Query("status", "")

	resp, err := s.bankClient.ListBanks(context.Background(), &pb.ListBanksRequest{
		Page:    int32(page),
		Limit:   int32(limit),
		Search:  search,
		Country: country,
		Status:  status,
	})
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

func (s *Server) handleGetBank(c *fiber.Ctx) error {
	bankID := c.Params("id")
	if bankID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Bank ID is required")
	}

	resp, err := s.bankClient.GetBank(context.Background(), &pb.GetBankRequest{
		BankId: bankID,
	})
	if err != nil {
		return err
	}

	return c.JSON(resp)
}
