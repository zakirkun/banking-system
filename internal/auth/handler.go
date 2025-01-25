package auth

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	svc Service
}

type RegisterRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/validate", h.ValidateToken)
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	resp, err := h.svc.Register(req.Username, req.Password, req.Email, req.PhoneNumber)
	if err != nil {
		if err == ErrUserExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "User already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(AuthResponse{
		Token:     resp.Token,
		UserID:    resp.UserID,
		ExpiresAt: resp.ExpiresAt,
	})
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	resp, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		if err == ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to login",
		})
	}

	return c.JSON(AuthResponse{
		Token:     resp.Token,
		UserID:    resp.UserID,
		ExpiresAt: resp.ExpiresAt,
	})
}

func (h *Handler) ValidateToken(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token is required",
		})
	}

	resp, err := h.svc.ValidateToken(token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to validate token",
		})
	}

	return c.JSON(fiber.Map{
		"is_valid": resp.IsValid,
		"user_id":  resp.UserID,
	})
}
