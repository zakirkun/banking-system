package gateway

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	pb "github.com/zakirkun/banking-microservices/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	app            *fiber.App
	authClient     pb.AuthServiceClient
	bankClient     pb.BankServiceClient
	cardClient     pb.CardServiceClient
	customerClient pb.CustomerServiceClient
}

func NewServer(cfg *Config) (*Server, error) {
	// Initialize gRPC connections
	authConn, err := grpc.Dial(cfg.AuthService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %v", err)
	}

	bankConn, err := grpc.Dial(cfg.BankService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to bank service: %v", err)
	}

	cardConn, err := grpc.Dial(cfg.CardService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to card service: %v", err)
	}

	customerConn, err := grpc.Dial(cfg.CustomerService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to customer service: %v", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	// Add middlewares
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	server := &Server{
		app:            app,
		authClient:     pb.NewAuthServiceClient(authConn),
		bankClient:     pb.NewBankServiceClient(bankConn),
		cardClient:     pb.NewCardServiceClient(cardConn),
		customerClient: pb.NewCustomerServiceClient(customerConn),
	}

	// Register routes
	server.registerRoutes()

	return server, nil
}

func (s *Server) Start(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) registerRoutes() {
	// API version group
	api := s.app.Group("/api/v1")

	// Public routes
	auth := api.Group("/auth")
	auth.Post("/register", s.handleRegister)
	auth.Post("/login", s.handleLogin)

	// Protected routes
	protected := api.Group("/", s.authMiddleware)

	// Customer routes
	customers := protected.Group("/customers")
	customers.Get("/me", s.handleGetCustomerProfile)
	customers.Put("/me", s.handleUpdateCustomerProfile)

	// Bank routes
	banks := protected.Group("/banks")
	banks.Get("/", s.handleListBanks)
	banks.Get("/:id", s.handleGetBank)

	// Card routes
	cards := protected.Group("/cards")
	cards.Post("/", s.handleIssueCard)
	cards.Get("/", s.handleListCards)
	cards.Get("/:id", s.handleGetCard)
	cards.Post("/:id/block", s.handleBlockCard)
	cards.Post("/:id/unblock", s.handleUnblockCard)

}
