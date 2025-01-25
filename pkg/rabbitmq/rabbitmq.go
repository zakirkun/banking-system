package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
}

func NewRabbitMQ(cfg Config) (*amqp.Connection, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.User, cfg.Password, cfg.Host, cfg.Port)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return conn, nil
}
