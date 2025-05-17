package main

import (
	"GabrielChaves1/notilify/internal/application/logger"
	"GabrielChaves1/notilify/internal/cache"
	"GabrielChaves1/notilify/internal/cache/redis"
	"GabrielChaves1/notilify/internal/queue"
	"GabrielChaves1/notilify/internal/queue/rabbitmq"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func initializeDependencies(cfg *Config) (queue.Repository, cache.Repository, error) {
	amqpClient, err := rabbitmq.NewRabbitMQClient(cfg.GetRabbitMQConnectionString())
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't initialize rabbitmq: %w", err)
	}

	redisClient := redis.NewRedisClient(cfg.GetRedisConnectionString())

	return amqpClient, redisClient, nil
}

func main() {
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}

	loggerConfig := logger.Config{
		Environment: config.environment,
	}

	logger := logger.NewLogger(loggerConfig)

	amqpClient, redisClient, err := initializeDependencies(config)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	logger.Printf("[*] Waiting for messages. To exit press CTRL+C")

	err = redisClient.SubscribeExpired(ctx, func(key string) {
		id := strings.TrimPrefix(key, "notification:")
		logger.Infof("Publishing notification %s", id)
		amqpClient.Publish(ctx, "notifications_queue", id)
	})
	if err != nil {
		logger.Fatal(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	logger.Print("Shutdown...")
}
