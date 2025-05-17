package mailer

import "context"

type Client interface {
	SendMail(ctx context.Context, to, subject, text string) error
}
