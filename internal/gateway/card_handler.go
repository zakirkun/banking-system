package gateway

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	pb "github.com/zakirkun/banking-microservices/api/proto"
)

func (s *Server) handleIssueCard(c *fiber.Ctx) error {
	var req struct {
		AccountID   string `json:"account_id"`
		CardType    string `json:"card_type"`
		CardNetwork string `json:"card_network"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	customerID := c.Locals("customerID").(string)

	resp, err := s.cardClient.IssueCard(context.Background(), &pb.IssueCardRequest{
		CustomerId:  customerID,
		AccountId:   req.AccountID,
		CardType:    req.CardType,
		CardNetwork: req.CardNetwork,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (s *Server) handleGetCard(c *fiber.Ctx) error {
	cardID := c.Params("id")
	if cardID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Card ID is required")
	}

	resp, err := s.cardClient.GetCard(context.Background(), &pb.GetCardRequest{
		CardId: cardID,
	})
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

func (s *Server) handleBlockCard(c *fiber.Ctx) error {
	cardID := c.Params("id")
	if cardID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Card ID is required")
	}

	var req struct {
		Reason string `json:"reason"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	resp, err := s.cardClient.BlockCard(context.Background(), &pb.BlockCardRequest{
		CardId: cardID,
		Reason: req.Reason,
	})
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

func (s *Server) handleUnblockCard(c *fiber.Ctx) error {
	cardID := c.Params("id")
	if cardID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Card ID is required")
	}

	resp, err := s.cardClient.UnblockCard(context.Background(), &pb.UnblockCardRequest{
		CardId: cardID,
	})
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

func (s *Server) handleListCards(c *fiber.Ctx) error {
	customerID := c.Locals("customerID").(string)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	status := c.Query("status", "")

	resp, err := s.cardClient.ListCards(context.Background(), &pb.ListCardsRequest{
		CustomerId: customerID,
		Page:       int32(page),
		Limit:      int32(limit),
		Status:     status,
	})
	if err != nil {
		return err
	}

	return c.JSON(resp)
}
