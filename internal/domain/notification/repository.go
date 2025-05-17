package notification

import "context"

type Repository interface {
	Save(ctx context.Context, notification *Notification) error
	GetByID(ctx context.Context, notificationID ID) (*Notification, error)
}

