package main

import (
	"GabrielChaves1/notilify/internal/domain/types"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	rabbitmqConnectionString string
	redisConnectionString    string
	environment              types.Environment
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

	rabbitmqConnectionString := os.Getenv("RABBITMQ_CONNECTION_STRING")
	if rabbitmqConnectionString == "" {
		return nil, fmt.Errorf("env var RABBITMQ_CONNECTION_STRING not defined")
	}

	redisConnectionString := os.Getenv("REDIS_CONNECTION_STRING")
	if redisConnectionString == "" {
		return nil, fmt.Errorf("env var REDIS_CONNECTION_STRING not defined")
	}

	return &Config{
		rabbitmqConnectionString: rabbitmqConnectionString,
		redisConnectionString:    redisConnectionString,
		environment:              environment,
	}, nil
}

func (c *Config) GetRabbitMQConnectionString() string {
	return c.rabbitmqConnectionString
}

func (c *Config) GetRedisConnectionString() string {
	return c.redisConnectionString
}
