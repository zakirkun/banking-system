package bank

import (
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ErrBankNotFound = errors.New("bank not found")
	ErrBankExists   = errors.New("bank code already exists")
)

type Service interface {
	CreateBank(bankName, bankCode, swiftCode, country, currency string) (*BankResponse, error)
	GetBank(id string) (*BankResponse, error)
	UpdateBank(id, bankName, swiftCode, status string) (*BankResponse, error)
	DeleteBank(id string) error
	ListBanks(page, limit int, search, country, status string) (*ListBanksResponse, error)
}

type service struct {
	repo Repository
	rmq  *amqp.Connection
}

type BankResponse struct {
	ID        string
	BankName  string
	BankCode  string
	SwiftCode string
	Country   string
	Currency  string
	Status    string
	CreatedAt string
	UpdatedAt string
}

type ListBanksResponse struct {
	Banks []BankResponse
	Total int64
}

func NewService(repo Repository, rmq *amqp.Connection) Service {
	return &service{repo: repo, rmq: rmq}
}

func (s *service) CreateBank(bankName, bankCode, swiftCode, country, currency string) (*BankResponse, error) {
	existing, err := s.repo.GetByCode(bankCode)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrBankExists
	}

	bank := &Bank{
		BankName:  bankName,
		BankCode:  bankCode,
		SwiftCode: swiftCode,
		Country:   country,
		Currency:  currency,
		Status:    "active",
	}

	if err := s.repo.Create(bank); err != nil {
		return nil, err
	}

	return s.bankToResponse(bank), nil
}

func (s *service) GetBank(id string) (*BankResponse, error) {
	bank, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if bank == nil {
		return nil, ErrBankNotFound
	}

	return s.bankToResponse(bank), nil
}

func (s *service) UpdateBank(id, bankName, swiftCode, status string) (*BankResponse, error) {
	bank, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if bank == nil {
		return nil, ErrBankNotFound
	}

	if bankName != "" {
		bank.BankName = bankName
	}
	if swiftCode != "" {
		bank.SwiftCode = swiftCode
	}
	if status != "" {
		bank.Status = status
	}

	if err := s.repo.Update(bank); err != nil {
		return nil, err
	}

	return s.bankToResponse(bank), nil
}

func (s *service) DeleteBank(id string) error {
	bank, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if bank == nil {
		return ErrBankNotFound
	}

	return s.repo.Delete(id)
}

func (s *service) ListBanks(page, limit int, search, country, status string) (*ListBanksResponse, error) {
	banks, total, err := s.repo.List(page, limit, search, country, status)
	if err != nil {
		return nil, err
	}

	var responses []BankResponse
	for _, bank := range banks {
		responses = append(responses, *s.bankToResponse(&bank))
	}

	return &ListBanksResponse{
		Banks: responses,
		Total: total,
	}, nil
}

func (s *service) bankToResponse(b *Bank) *BankResponse {
	return &BankResponse{
		ID:        b.ID,
		BankName:  b.BankName,
		BankCode:  b.BankCode,
		SwiftCode: b.SwiftCode,
		Country:   b.Country,
		Currency:  b.Currency,
		Status:    b.Status,
		CreatedAt: b.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: b.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
