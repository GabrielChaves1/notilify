package request

import "time"

type CreateNotificationDTO struct {
	Channel     string     `json:"channel"`
	Recipient   string     `json:"recipient"`
	Message     string     `json:"message"`
	Priority    string     `json:"priority"`
	ScheduledAt *time.Time `json:"scheduled_at"`
	Data        string     `json:"data"`
}
