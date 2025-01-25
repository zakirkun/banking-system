package customer

import (
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrUserExists       = errors.New("user already has a customer profile")
)

type Service interface {
	CreateCustomer(userID, fullName, dateOfBirth, address, idNumber, idType string) (*CustomerResponse, error)
	GetCustomer(id string) (*CustomerResponse, error)
	UpdateCustomer(id, fullName, address string) (*CustomerResponse, error)
	DeleteCustomer(id string) error
	ListCustomers(page, limit int, search string) (*ListCustomersResponse, error)
}

type service struct {
	repo Repository
	rmq  *amqp.Connection
}

type CustomerResponse struct {
	ID          string
	UserID      string
	FullName    string
	DateOfBirth string
	Address     string
	IDNumber    string
	IDType      string
	Status      string
	CreatedAt   string
	UpdatedAt   string
}

type ListCustomersResponse struct {
	Customers []CustomerResponse
	Total     int64
}

func NewService(repo Repository, rmq *amqp.Connection) Service {
	return &service{repo: repo, rmq: rmq}
}

func (s *service) CreateCustomer(userID, fullName, dateOfBirth, address, idNumber, idType string) (*CustomerResponse, error) {
	// Check if user already has a customer profile
	existing, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserExists
	}

	customer := &Customer{
		UserID:      userID,
		FullName:    fullName,
		DateOfBirth: dateOfBirth,
		Address:     address,
		IDNumber:    idNumber,
		IDType:      idType,
		Status:      "active",
	}

	if err := s.repo.Create(customer); err != nil {
		return nil, err
	}

	return s.customerToResponse(customer), nil
}

func (s *service) GetCustomer(id string) (*CustomerResponse, error) {
	customer, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, ErrCustomerNotFound
	}

	return s.customerToResponse(customer), nil
}

func (s *service) UpdateCustomer(id, fullName, address string) (*CustomerResponse, error) {
	customer, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, ErrCustomerNotFound
	}

	if fullName != "" {
		customer.FullName = fullName
	}
	if address != "" {
		customer.Address = address
	}

	if err := s.repo.Update(customer); err != nil {
		return nil, err
	}

	return s.customerToResponse(customer), nil
}

func (s *service) DeleteCustomer(id string) error {
	customer, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if customer == nil {
		return ErrCustomerNotFound
	}

	return s.repo.Delete(id)
}

func (s *service) ListCustomers(page, limit int, search string) (*ListCustomersResponse, error) {
	customers, total, err := s.repo.List(page, limit, search)
	if err != nil {
		return nil, err
	}

	var responses []CustomerResponse
	for _, customer := range customers {
		responses = append(responses, *s.customerToResponse(&customer))
	}

	return &ListCustomersResponse{
		Customers: responses,
		Total:     total,
	}, nil
}

func (s *service) customerToResponse(c *Customer) *CustomerResponse {
	return &CustomerResponse{
		ID:          c.ID,
		UserID:      c.UserID,
		FullName:    c.FullName,
		DateOfBirth: c.DateOfBirth,
		Address:     c.Address,
		IDNumber:    c.IDNumber,
		IDType:      c.IDType,
		Status:      c.Status,
		CreatedAt:   c.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   c.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
