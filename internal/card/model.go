package card

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Card struct {
	ID          string `gorm:"primaryKey;type:uuid"`
	CustomerID  string `gorm:"not null;index"`
	AccountID   string `gorm:"not null;index"`
	CardNumber  string `gorm:"uniqueIndex;not null"`
	CardType    string `gorm:"not null"` // DEBIT, CREDIT
	CardNetwork string `gorm:"not null"` // VISA, MASTERCARD
	ExpiryDate  string `gorm:"not null"`
	CVV         string `gorm:"not null"`
	Status      string `gorm:"not null;default:'active'"` // active, blocked, expired
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (c *Card) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}
