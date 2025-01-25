package bank

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bank struct {
	ID        string `gorm:"primaryKey;type:uuid"`
	BankName  string `gorm:"not null"`
	BankCode  string `gorm:"uniqueIndex;not null"`
	SwiftCode string `gorm:"uniqueIndex;not null"`
	Country   string `gorm:"not null"`
	Currency  string `gorm:"not null"`
	Status    string `gorm:"not null;default:'active'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *Bank) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}
