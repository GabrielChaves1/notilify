package mailer

import (
	"GabrielChaves1/notilify/internal/domain/notification"
	"context"
)

type MailSender struct {
	client Client
}

func NewMailSender(client Client) *MailSender {
	return &MailSender{client: client}
}

func (s *MailSender) Send(ctx context.Context, notification *notification.Notification) error {
	return s.client.SendMail(ctx, notification.Recipient, "Notificação", notification.Message)
}
