package main

import (
	"GabrielChaves1/notilify/internal/domain/types"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	postgresConnectionString string
	rabbitmqConnectionString string
	fromSourceEmail          string
	environment              types.Environment
	maxRetries               int
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("couldn't load environment variables: %w", err)
	}

	env := os.Getenv("ENVIRONMENT")
	var environment types.Environment
	switch env {
	case "production":
		environment = types.Production
	case "staging":
		environment = types.Staging
	default:
		environment = types.Development
	}

	postgresConnectionString := os.Getenv("POSTGRES_CONNECTION_STRING")
	if postgresConnectionString == "" {
		return nil, fmt.Errorf("env var POSTGRES_CONNECTION_STRING not defined")
	}

	rabbitmqConnectionString := os.Getenv("RABBITMQ_CONNECTION_STRING")
	if rabbitmqConnectionString == "" {
		return nil, fmt.Errorf("env var RABBITMQ_CONNECTION_STRING not defined")
	}

	fromSourceEmail := os.Getenv("FROM_SOURCE_EMAIL")
	if fromSourceEmail == "" {
		return nil, fmt.Errorf("env var FROM_SOURCE_EMAIL not defined")
	}

	return &Config{
		postgresConnectionString: postgresConnectionString,
		rabbitmqConnectionString: rabbitmqConnectionString,
		fromSourceEmail:          fromSourceEmail,
		environment:              environment,
		maxRetries:               5,
	}, nil
}

func (c Config) GetPostgresConnectionString() string {
	return c.postgresConnectionString
}

func (c Config) GetRabbitMQConnectionString() string {
	return c.rabbitmqConnectionString
}
