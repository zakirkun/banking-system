package card

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	cards := api.Group("/cards")

	cards.Post("/", h.IssueCard)
	cards.Get("/:id", h.GetCard)
	cards.Post("/:id/block", h.BlockCard)
	cards.Post("/:id/unblock", h.UnblockCard)
	cards.Get("/", h.ListCards)
}

func (h *Handler) IssueCard(c *fiber.Ctx) error {
	var req struct {
		CustomerID  string `json:"customer_id"`
		AccountID   string `json:"account_id"`
		CardType    string `json:"card_type"`
		CardNetwork string `json:"card_network"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	resp, err := h.svc.IssueCard(req.CustomerID, req.AccountID, req.CardType, req.CardNetwork)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to issue card",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *Handler) GetCard(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Card ID is required",
		})
	}

	resp, err := h.svc.GetCard(id)
	if err != nil {
		if err == ErrCardNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Card not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get card",
		})
	}

	return c.JSON(resp)
}

func (h *Handler) BlockCard(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Card ID is required",
		})
	}

	var req struct {
		Reason string `json:"reason"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	resp, err := h.svc.BlockCard(id, req.Reason)
	if err != nil {
		if err == ErrCardNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Card not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to block card",
		})
	}

	return c.JSON(resp)
}

func (h *Handler) UnblockCard(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Card ID is required",
		})
	}

	resp, err := h.svc.UnblockCard(id)
	if err != nil {
		if err == ErrCardNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Card not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to unblock card",
		})
	}

	return c.JSON(resp)
}

func (h *Handler) ListCards(c *fiber.Ctx) error {
	customerID := c.Query("customer_id")
	if customerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Customer ID is required",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	status := c.Query("status", "")

	resp, err := h.svc.ListCards(customerID, page, limit, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list cards",
		})
	}

	return c.JSON(resp)
}
