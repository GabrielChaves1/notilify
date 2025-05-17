package main

import (
	"GabrielChaves1/notilify/internal/application/logger"
	"GabrielChaves1/notilify/internal/application/manager"
	"GabrielChaves1/notilify/internal/application/usecase"
	"GabrielChaves1/notilify/internal/clients/mailer"
	awsses "GabrielChaves1/notilify/internal/clients/mailer/ses"
	notifications "GabrielChaves1/notilify/internal/clients/notification"
	awssns "GabrielChaves1/notilify/internal/clients/notification/sns"
	"GabrielChaves1/notilify/internal/domain/notification"
	"GabrielChaves1/notilify/internal/queue"
	"GabrielChaves1/notilify/internal/queue/rabbitmq"
	"GabrielChaves1/notilify/internal/storage/postgres"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/jmoiron/sqlx"
)

func initializeDependencies(cfg *Config) (*sqlx.DB, queue.Repository, *ses.Client, *sns.Client, error) {
	postgresClient, err := postgres.NewPostgresClient(cfg.GetPostgresConnectionString(), postgres.DefaultPostgresOptions())
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("couldn't initialize postgres: %w", err)
	}

	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("sa-east-1"))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("couldn't load aws config")
	}

	sesClient := ses.NewFromConfig(awsCfg)
	snsClient := sns.NewFromConfig(awsCfg)

	amqpClient, err := rabbitmq.NewRabbitMQClient(cfg.GetRabbitMQConnectionString())
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("couldn't initialize rabbitmq: %w", err)
	}

	return postgresClient, amqpClient, sesClient, snsClient, nil
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

	postgresClient, amqpClient, sesClient, snsClient, err := initializeDependencies(config)
	if err != nil {
		panic(err)
	}

	ses := awsses.NewSESClient(sesClient, config.fromSourceEmail)
	emailSender := mailer.NewMailSender(ses)

	sns := awssns.NewSNSClient(snsClient)
	notificationSender := notifications.NewSMSSender(sns)

	channelManager := manager.NewChannelManager(emailSender, notificationSender)

	notificationRepo := postgres.NewNotificationRepository(postgresClient)
	handleSendNotificationTaskUseCase := usecase.NewHandleSendNotificationTask(notificationRepo, channelManager, logger)

	err = amqpClient.Consume(context.Background(), queue.NotificationQueue, func(msg string) error {
		ctx := context.Background()
		notificationID, err := notification.NewIDFromString(msg)
		if err != nil {
			return err
		}

		if err := handleSendNotificationTaskUseCase.Execute(ctx, notificationID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("[*] Waiting for messages. To exit press CTRL+C")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	logger.Print("Shutdown...")
}
