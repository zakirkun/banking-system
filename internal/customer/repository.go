package customer

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(customer *Customer) error
	GetByID(id string) (*Customer, error)
	GetByUserID(userID string) (*Customer, error)
	Update(customer *Customer) error
	Delete(id string) error
	List(page, limit int, search string) ([]Customer, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(customer *Customer) error {
	return r.db.Create(customer).Error
}

func (r *repository) GetByID(id string) (*Customer, error) {
	var customer Customer
	if err := r.db.First(&customer, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

func (r *repository) GetByUserID(userID string) (*Customer, error) {
	var customer Customer
	if err := r.db.First(&customer, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

func (r *repository) Update(customer *Customer) error {
	return r.db.Save(customer).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&Customer{}, "id = ?", id).Error
}

func (r *repository) List(page, limit int, search string) ([]Customer, int64, error) {
	var customers []Customer
	var total int64

	query := r.db.Model(&Customer{})

	if search != "" {
		query = query.Where("full_name ILIKE ? OR id_number LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}
