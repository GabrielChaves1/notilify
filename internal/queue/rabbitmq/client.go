package rabbitmq

import (
	"GabrielChaves1/notilify/internal/queue"
	"context"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	client   *amqp091.Connection
	ch       *amqp091.Channel
	exchange string
}

func NewRabbitMQClient(url string) (queue.Repository, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	ch.ExchangeDeclare(queue.DLQNotificationExchange, "direct", true, false, false, false, nil)
	ch.QueueDeclare(queue.DlqNotificationQueue, true, false, false, false, nil)
	ch.QueueBind(queue.DlqNotificationQueue, "dlq_notifications_queue_rt", queue.DLQNotificationExchange, false, nil)

	return &RabbitMQ{
		client: conn,
		ch:     ch,
	}, nil
}

func (q *RabbitMQ) Publish(ctx context.Context, queueName, message string) error {
	_, err := q.ch.QueueDeclare(queueName, true, false, false, false, amqp091.Table{
		"x-dead-letter-exchange":    queue.DLQNotificationExchange,
		"x-dead-letter-routing-key": "dlq_notifications_queue_rt",
		"x-queue-type":              "quorum",
		"x-delivery-limit":          3,
	})
	if err != nil {
		return err
	}

	return q.ch.PublishWithContext(ctx, "", queueName, false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	})
}

func (q *RabbitMQ) Consume(ctx context.Context, queueName string, handler func(msg string) error) error {
	msgs, err := q.ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			if err := handler(string(msg.Body)); err != nil {
				msg.Nack(false, true)
			}

			msg.Nack(false, false)
		}
	}()

	return err
}
