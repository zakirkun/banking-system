package customer

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
	customers := api.Group("/customers")

	customers.Post("/", h.CreateCustomer)
	customers.Get("/:id", h.GetCustomer)
	customers.Put("/:id", h.UpdateCustomer)
	customers.Delete("/:id", h.DeleteCustomer)
	customers.Get("/", h.ListCustomers)
}

func (h *Handler) CreateCustomer(c *fiber.Ctx) error {
	var req struct {
		UserID      string `json:"user_id"`
		FullName    string `json:"full_name"`
		DateOfBirth string `json:"date_of_birth"`
		Address     string `json:"address"`
		IDNumber    string `json:"id_number"`
		IDType      string `json:"id_type"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	resp, err := h.svc.CreateCustomer(req.UserID, req.FullName, req.DateOfBirth, req.Address, req.IDNumber, req.IDType)
	if err != nil {
		if err == ErrUserExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create customer",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *Handler) GetCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Customer ID is required",
		})
	}

	resp, err := h.svc.GetCustomer(id)
	if err != nil {
		if err == ErrCustomerNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Customer not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get customer",
		})
	}

	return c.JSON(resp)
}

func (h *Handler) UpdateCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Customer ID is required",
		})
	}

	var req struct {
		FullName string `json:"full_name"`
		Address  string `json:"address"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	resp, err := h.svc.UpdateCustomer(id, req.FullName, req.Address)
	if err != nil {
		if err == ErrCustomerNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Customer not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update customer",
		})
	}

	return c.JSON(resp)
}

func (h *Handler) DeleteCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Customer ID is required",
		})
	}

	if err := h.svc.DeleteCustomer(id); err != nil {
		if err == ErrCustomerNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Customer not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete customer",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) ListCustomers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")

	resp, err := h.svc.ListCustomers(page, limit, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list customers",
		})
	}

	return c.JSON(resp)
}
