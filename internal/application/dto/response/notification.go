package response

import "time"

type NotificationDTO struct {
	ID          string     `json:"id"`
	Channel     string     `json:"channel"`
	Recipient   string     `json:"recipient"`
	Message     string     `json:"message"`
	Priority    string     `json:"priority"`
	Data        string     `json:"data"`
	Status      string     `json:"status"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
