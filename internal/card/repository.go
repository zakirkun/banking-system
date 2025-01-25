package card

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(card *Card) error
	GetByID(id string) (*Card, error)
	GetByCardNumber(cardNumber string) (*Card, error)
	Update(card *Card) error
	List(customerID string, page, limit int, status string) ([]Card, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(card *Card) error {
	return r.db.Create(card).Error
}

func (r *repository) GetByID(id string) (*Card, error) {
	var card Card
	if err := r.db.First(&card, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &card, nil
}

func (r *repository) GetByCardNumber(cardNumber string) (*Card, error) {
	var card Card
	if err := r.db.First(&card, "card_number = ?", cardNumber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &card, nil
}

func (r *repository) Update(card *Card) error {
	return r.db.Save(card).Error
}

func (r *repository) List(customerID string, page, limit int, status string) ([]Card, int64, error) {
	var cards []Card
	var total int64

	query := r.db.Model(&Card{}).Where("customer_id = ?", customerID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&cards).Error; err != nil {
		return nil, 0, err
	}

	return cards, total, nil
}
