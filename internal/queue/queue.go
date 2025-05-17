package queue

import (
	"context"
)

type Repository interface {
	Publish(ctx context.Context, queueName, message string) error
	Consume(ctx context.Context, queueName string, handler func(msg string) error) error
}
