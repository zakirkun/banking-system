package bank

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(bank *Bank) error
	GetByID(id string) (*Bank, error)
	GetByCode(code string) (*Bank, error)
	Update(bank *Bank) error
	Delete(id string) error
	List(page, limit int, search, country, status string) ([]Bank, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(bank *Bank) error {
	return r.db.Create(bank).Error
}

func (r *repository) GetByID(id string) (*Bank, error) {
	var bank Bank
	if err := r.db.First(&bank, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &bank, nil
}

func (r *repository) GetByCode(code string) (*Bank, error) {
	var bank Bank
	if err := r.db.First(&bank, "bank_code = ?", code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &bank, nil
}

func (r *repository) Update(bank *Bank) error {
	return r.db.Save(bank).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&Bank{}, "id = ?", id).Error
}

func (r *repository) List(page, limit int, search, country, status string) ([]Bank, int64, error) {
	var banks []Bank
	var total int64

	query := r.db.Model(&Bank{})

	if search != "" {
		query = query.Where("bank_name ILIKE ? OR bank_code LIKE ? OR swift_code LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if country != "" {
		query = query.Where("country = ?", country)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&banks).Error; err != nil {
		return nil, 0, err
	}

	return banks, total, nil
}
