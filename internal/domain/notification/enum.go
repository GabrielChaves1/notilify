package notification

type NotificationChannel string

const (
	EmailChannel NotificationChannel = "email"
	SMSChannel   NotificationChannel = "sms"
)

func (c NotificationChannel) IsValid() bool {
	switch c {
	case EmailChannel, SMSChannel:
		return true
	}
	return false
}

type NotificationStatus string

const (
	Pending   NotificationStatus = "pending"
	Delivered NotificationStatus = "delivered"
)
