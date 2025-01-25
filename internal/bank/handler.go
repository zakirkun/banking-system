package bank

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
	banks := api.Group("/banks")

	banks.Post("/", h.CreateBank)
	banks.Get("/:id", h.GetBank)
	banks.Put("/:id", h.UpdateBank)
	banks.Delete("/:id", h.DeleteBank)
	banks.Get("/", h.ListBanks)
}

func (h *Handler) CreateBank(c *fiber.Ctx) error {
	var req struct {
		BankName  string `json:"bank_name"`
		BankCode  string `json:"bank_code"`
		SwiftCode string `json:"swift_code"`
		Country   string `json:"country"`
		Currency  string `json:"currency"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	resp, err := h.svc.CreateBank(req.BankName, req.BankCode, req.SwiftCode, req.Country, req.Currency)
	if err != nil {
		if err == ErrBankExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create bank",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *Handler) GetBank(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bank ID is required",
		})
	}

	resp, err := h.svc.GetBank(id)
	if err != nil {
		if err == ErrBankNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Bank not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get bank",
		})
	}

	return c.JSON(resp)
}

func (h *Handler) UpdateBank(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bank ID is required",
		})
	}

	var req struct {
		BankName  string `json:"bank_name"`
		SwiftCode string `json:"swift_code"`
		Status    string `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	resp, err := h.svc.UpdateBank(id, req.BankName, req.SwiftCode, req.Status)
	if err != nil {
		if err == ErrBankNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Bank not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update bank",
		})
	}

	return c.JSON(resp)
}

func (h *Handler) DeleteBank(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bank ID is required",
		})
	}

	if err := h.svc.DeleteBank(id); err != nil {
		if err == ErrBankNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Bank not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete bank",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) ListBanks(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")
	country := c.Query("country", "")
	status := c.Query("status", "")

	resp, err := h.svc.ListBanks(page, limit, search, country, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list banks",
		})
	}

	return c.JSON(resp)
}
