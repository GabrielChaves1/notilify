package notification

import (
	"time"

	"github.com/google/uuid"
)

func NewNotification(
	channel NotificationChannel,
	recipient string,
	message string,
	data string,
	currentTime time.Time,
	scheduledAt *time.Time,
) *Notification {
	return &Notification{
		ID:          ID(uuid.New()),
		Channel:     channel,
		Recipient:   recipient,
		Message:     message,
		Data:        data,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
		ScheduledAt: scheduledAt,
		DeletedAt:   nil,
	}
}

func NewIDFromString(id string) (ID, error) {
	notificationUUID, err := uuid.Parse(id)
	if err != nil {
		return ID{}, err
	}

	return ID(notificationUUID), err
}
