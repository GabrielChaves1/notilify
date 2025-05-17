package manager

import (
	"GabrielChaves1/notilify/internal/domain/notification"
)

type ChannelManager struct {
	emailSender notification.ChannelSender
	smsSender   notification.ChannelSender
}

func NewChannelManager(emailSender, smsSender notification.ChannelSender) *ChannelManager {
	return &ChannelManager{
		emailSender: emailSender,
		smsSender:   smsSender,
	}
}

func (m *ChannelManager) GetChannel(channelType notification.NotificationChannel) notification.ChannelSender {
	switch channelType {
	case notification.EmailChannel:
		return m.emailSender
	case notification.SMSChannel:
		return m.smsSender
	default:
		return m.emailSender
	}
}
