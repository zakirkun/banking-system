package database

import (
	"github.com/zakirkun/banking-microservices/internal/auth"
	"github.com/zakirkun/banking-microservices/internal/bank"
	"github.com/zakirkun/banking-microservices/internal/card"
	"github.com/zakirkun/banking-microservices/internal/customer"
	"gorm.io/gorm"
)

// AutoMigrateAll melakukan migrasi untuk semua model database
func AutoMigrateAll(db *gorm.DB) error {
	return db.AutoMigrate(
		&auth.User{},
		&bank.Bank{},
		&card.Card{},
		&customer.Customer{},
	)
}

// AutoMigrateAuth melakukan migrasi untuk auth service
func AutoMigrateAuth(db *gorm.DB) error {
	return db.AutoMigrate(&auth.User{})
}

// AutoMigrateBank melakukan migrasi untuk bank service
func AutoMigrateBank(db *gorm.DB) error {
	return db.AutoMigrate(&bank.Bank{})
}

// AutoMigrateCard melakukan migrasi untuk card service
func AutoMigrateCard(db *gorm.DB) error {
	return db.AutoMigrate(&card.Card{})
}

// AutoMigrateCustomer melakukan migrasi untuk customer service
func AutoMigrateCustomer(db *gorm.DB) error {
	return db.AutoMigrate(&customer.Customer{})
}
