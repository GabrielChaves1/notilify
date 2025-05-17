package notification

import (
	"context"
)

type ChannelSender interface {
	Send(ctx context.Context, notification *Notification) error
}
