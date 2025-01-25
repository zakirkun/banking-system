package main

import (
	"log"
	"net"

	"github.com/gofiber/fiber/v2"
	pb "github.com/zakirkun/banking-microservices/api/proto"
	"github.com/zakirkun/banking-microservices/internal/card"
	"github.com/zakirkun/banking-microservices/pkg/config"
	"github.com/zakirkun/banking-microservices/pkg/database"
	"github.com/zakirkun/banking-microservices/pkg/rabbitmq"
	"google.golang.org/grpc"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.NewPostgresDB(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate database
	if err := db.AutoMigrate(&card.Card{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize RabbitMQ
	rmq, err := rabbitmq.NewRabbitMQ(rabbitmq.Config{
		Host:     cfg.RabbitMQ.Host,
		Port:     cfg.RabbitMQ.Port,
		User:     cfg.RabbitMQ.User,
		Password: cfg.RabbitMQ.Password,
	})
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rmq.Close()

	// Initialize repository
	repo := card.NewRepository(db)

	// Initialize service
	svc := card.NewService(repo, rmq)

	// Create gRPC listener
	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	grpcSvc := card.NewGRPCServer(svc)
	pb.RegisterCardServiceServer(grpcServer, grpcSvc)

	// Initialize HTTP server
	app := fiber.New()

	// Register HTTP handlers
	handler := card.NewHandler(svc)
	handler.RegisterRoutes(app)

	// Start servers
	go func() {
		if err := app.Listen(":" + cfg.HTTP.Port); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	log.Printf("Starting %s service...", cfg.ServiceName)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
