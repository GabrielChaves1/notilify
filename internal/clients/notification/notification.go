package notifications

import "context"

type Client interface {
	SendSMS(ctx context.Context, phoneNumber, message string) error
}
