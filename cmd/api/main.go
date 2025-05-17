package main

import (
	"GabrielChaves1/notilify/internal/api/handlers"
	"GabrielChaves1/notilify/internal/api/router"
	"GabrielChaves1/notilify/internal/application/logger"
	"GabrielChaves1/notilify/internal/application/usecase"
	"GabrielChaves1/notilify/internal/cache"
	"GabrielChaves1/notilify/internal/cache/redis"
	"GabrielChaves1/notilify/internal/domain/notification"
	"GabrielChaves1/notilify/internal/queue"
	"GabrielChaves1/notilify/internal/queue/rabbitmq"
	"GabrielChaves1/notilify/internal/storage/postgres"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/doganarif/govisual"
	"github.com/jmoiron/sqlx"
)

func initializeDependencies(config *Config) (*sqlx.DB, queue.Repository, cache.Repository, error) {
	postgresClient, err := postgres.NewPostgresClient(
		config.GetPostgresConnectionString(),
		postgres.DefaultPostgresOptions(),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("couldn't initialize postgres: %w", err)
	}

	amqpClient, err := rabbitmq.NewRabbitMQClient(
		config.rabbitmqConnectionString,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("couldn't initialize rabbitmq: %w", err)
	}

	redisClient := redis.NewRedisClient(config.GetRedisConnectionString())

	return postgresClient, amqpClient, redisClient, nil
}

func main() {
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}

	postgresClient, amqpClient, redisClient, err := initializeDependencies(config)
	if err != nil {
		panic(err)
	}

	postgresNotificationRepository := postgres.NewNotificationRepository(postgresClient)
	notificationValidator := notification.NewNotificationValidator()

	createNotificationUseCase := usecase.NewCreateNotification(postgresNotificationRepository, amqpClient, redisClient, notificationValidator)
	notificationHandlers := handlers.NewNotificationHandlers(createNotificationUseCase)

	loggerConfig := logger.Config{
		Environment: config.environment,
	}

	logger := logger.NewLogger(loggerConfig)

	routerConfig := router.APIRouterConfig{
		Environment: config.environment,
	}

	router := router.SetupAPIRouter(
		routerConfig,
		logger,
		notificationHandlers,
	)

	handler := govisual.Wrap(router, govisual.WithRequestBodyLogging(true), govisual.WithResponseBodyLogging(true))

	server := &http.Server{
		Addr:           ":8080",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logger.Printf("[*] Waiting for messages. To exit press CTRL+C")

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	logger.Println("Starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Println("Failed to shutdown HTTP server")
	} else {
		logger.Println("HTTP Server shutdown")
	}
}
