package customer

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	ID          string `gorm:"primaryKey;type:uuid"`
	UserID      string `gorm:"uniqueIndex;not null"`
	FullName    string `gorm:"not null"`
	DateOfBirth string `gorm:"not null"`
	Address     string `gorm:"not null"`
	IDNumber    string `gorm:"uniqueIndex;not null"`
	IDType      string `gorm:"not null"`
	Status      string `gorm:"not null;default:'active'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}
