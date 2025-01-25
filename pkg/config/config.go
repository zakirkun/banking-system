package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServiceName string

	HTTP struct {
		Port string
	}

	GRPC struct {
		Port string
	}

	Services struct {
		Auth     string
		Bank     string
		Card     string
		Customer string
	}

	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}

	RabbitMQ struct {
		Host     string
		Port     string
		User     string
		Password string
	}
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
