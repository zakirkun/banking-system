package card

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ErrCardNotFound = errors.New("card not found")
	ErrInvalidCard  = errors.New("invalid card")
)

type Service interface {
	IssueCard(customerID, accountID, cardType, cardNetwork string) (*CardResponse, error)
	GetCard(id string) (*CardResponse, error)
	BlockCard(id, reason string) (*CardResponse, error)
	UnblockCard(id string) (*CardResponse, error)
	ListCards(customerID string, page, limit int, status string) (*ListCardsResponse, error)
}

type service struct {
	repo Repository
	rmq  *amqp.Connection
}

type CardResponse struct {
	ID          string
	CustomerID  string
	AccountID   string
	CardNumber  string
	CardType    string
	CardNetwork string
	ExpiryDate  string
	CVV         string
	Status      string
	CreatedAt   string
	UpdatedAt   string
}

type ListCardsResponse struct {
	Cards []CardResponse
	Total int64
}

func NewService(repo Repository, rmq *amqp.Connection) Service {
	return &service{repo: repo, rmq: rmq}
}

func (s *service) IssueCard(customerID, accountID, cardType, cardNetwork string) (*CardResponse, error) {
	// Generate card details
	cardNumber := generateCardNumber()
	expiryDate := generateExpiryDate()
	cvv := generateCVV()

	card := &Card{
		CustomerID:  customerID,
		AccountID:   accountID,
		CardNumber:  cardNumber,
		CardType:    cardType,
		CardNetwork: cardNetwork,
		ExpiryDate:  expiryDate,
		CVV:         cvv,
		Status:      "active",
	}

	if err := s.repo.Create(card); err != nil {
		return nil, err
	}

	// TODO: Publish card issued event to RabbitMQ

	return s.cardToResponse(card), nil
}

func (s *service) GetCard(id string) (*CardResponse, error) {
	card, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, ErrCardNotFound
	}

	return s.cardToResponse(card), nil
}

func (s *service) BlockCard(id, reason string) (*CardResponse, error) {
	card, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, ErrCardNotFound
	}

	card.Status = "blocked"
	if err := s.repo.Update(card); err != nil {
		return nil, err
	}

	// TODO: Publish card blocked event to RabbitMQ

	return s.cardToResponse(card), nil
}

func (s *service) UnblockCard(id string) (*CardResponse, error) {
	card, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, ErrCardNotFound
	}

	card.Status = "active"
	if err := s.repo.Update(card); err != nil {
		return nil, err
	}

	// TODO: Publish card unblocked event to RabbitMQ

	return s.cardToResponse(card), nil
}

func (s *service) ListCards(customerID string, page, limit int, status string) (*ListCardsResponse, error) {
	cards, total, err := s.repo.List(customerID, page, limit, status)
	if err != nil {
		return nil, err
	}

	var responses []CardResponse
	for _, card := range cards {
		responses = append(responses, *s.cardToResponse(&card))
	}

	return &ListCardsResponse{
		Cards: responses,
		Total: total,
	}, nil
}

func (s *service) cardToResponse(c *Card) *CardResponse {
	return &CardResponse{
		ID:          c.ID,
		CustomerID:  c.CustomerID,
		AccountID:   c.AccountID,
		CardNumber:  maskCardNumber(c.CardNumber),
		CardType:    c.CardType,
		CardNetwork: c.CardNetwork,
		ExpiryDate:  c.ExpiryDate,
		CVV:         maskCVV(c.CVV),
		Status:      c.Status,
		CreatedAt:   c.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   c.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func generateCardNumber() string {
	return fmt.Sprintf("%016d", rand.Int63n(10000000000000000))
}

func generateExpiryDate() string {
	expiry := time.Now().AddDate(5, 0, 0)
	return expiry.Format("01/06")
}

func generateCVV() string {
	return fmt.Sprintf("%03d", rand.Intn(1000))
}

func maskCardNumber(number string) string {
	if len(number) != 16 {
		return number
	}
	return fmt.Sprintf("****-****-****-%s", number[12:])
}

func maskCVV(cvv string) string {
	return "***"
}
