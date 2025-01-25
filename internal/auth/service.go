package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	jwtSecret             = []byte("your-secret-key") // In production, use environment variable
)

type Service interface {
	Register(username, password, email, phoneNumber string) (*AuthResponse, error)
	Login(username, password string) (*AuthResponse, error)
	ValidateToken(token string) (*ValidateTokenResponse, error)
}

type service struct {
	repo Repository
	rmq  *amqp.Connection
}

type AuthResponse struct {
	Token     string
	UserID    string
	ExpiresAt string
}

type ValidateTokenResponse struct {
	IsValid bool
	UserID  string
}

func NewService(repo Repository, rmq *amqp.Connection) Service {
	return &service{repo: repo, rmq: rmq}
}

func (s *service) Register(username, password, email, phoneNumber string) (*AuthResponse, error) {
	// Check if user exists
	existingUser, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// Create new user
	user := &User{
		Username:    username,
		Password:    password,
		Email:       email,
		PhoneNumber: phoneNumber,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	// Generate token
	return s.generateToken(user.ID)
}

func (s *service) Login(username, password string) (*AuthResponse, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := user.ComparePassword(password); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.generateToken(user.ID)
}

func (s *service) ValidateToken(tokenString string) (*ValidateTokenResponse, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return &ValidateTokenResponse{IsValid: false}, nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(string)
		return &ValidateTokenResponse{
			IsValid: true,
			UserID:  userID,
		}, nil
	}

	return &ValidateTokenResponse{IsValid: false}, nil
}

func (s *service) generateToken(userID string) (*AuthResponse, error) {
	expiresAt := time.Now().Add(24 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt.Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:     tokenString,
		UserID:    userID,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
}
