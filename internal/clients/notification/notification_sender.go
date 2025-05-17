package notifications

import (
	"GabrielChaves1/notilify/internal/domain/notification"
	"context"
)

type SMSSender struct {
	client Client
}

func NewSMSSender(client Client) *SMSSender {
	return &SMSSender{client: client}
}

func (s *SMSSender) Send(ctx context.Context, notification *notification.Notification) error {
	return s.client.SendSMS(ctx, notification.Recipient, notification.Message)
}
