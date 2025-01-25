package auth

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID          string `gorm:"primaryKey;type:uuid"`
	Username    string `gorm:"uniqueIndex;not null"`
	Password    string `gorm:"not null"`
	Email       string `gorm:"uniqueIndex;not null"`
	PhoneNumber string `gorm:"uniqueIndex;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
