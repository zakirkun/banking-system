package main

import (
	"log"

	"github.com/zakirkun/banking-microservices/internal/gateway"
	"github.com/zakirkun/banking-microservices/pkg/config"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize gateway server
	server, err := gateway.NewServer(&gateway.Config{
		AuthService:     cfg.Services.Auth,
		BankService:     cfg.Services.Bank,
		CardService:     cfg.Services.Card,
		CustomerService: cfg.Services.Customer,
	})
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start server
	log.Printf("Starting API Gateway on port %s...", cfg.HTTP.Port)
	if err := server.Start(":" + cfg.HTTP.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
