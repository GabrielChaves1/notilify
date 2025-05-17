package notification

import (
	"GabrielChaves1/notilify/internal/application/dto/request"
	errors "GabrielChaves1/notilify/internal/application/error"
	"GabrielChaves1/notilify/internal/application/helper"
	"GabrielChaves1/notilify/internal/domain/types"
	"time"
)

type Validator struct{}

func NewNotificationValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(dto request.CreateNotificationDTO, currentTime time.Time) (*Notification, error) {
	validationErrs := &errors.ValidationErrors{}

	channel := NotificationChannel(dto.Channel)

	if dto.Channel == "" {
		validationErrs.Add("channel", channel, errors.Required, nil)
	}

	validChannels := map[string]bool{
		"email": true,
		"sms":   true,
	}

	if !validChannels[dto.Channel] {
		channels := ""

		for channel := range validChannels {
			channels += channel + ", "
		}

		channels = channels[:len(channels)-2]

		validationErrs.Add("channel", channel, errors.InvalidFormat, nil)
	}

	switch channel {
	case "email":
		if _, err := types.NewEmail(dto.Recipient); err != nil {
			validationErrs.Add("recipient", dto.Recipient, errors.InvalidFormat, "email must be in the format user@example.com")
		}
	case "sms":
		if _, err := types.NewPhoneNumber(dto.Recipient); err != nil {
			validationErrs.Add("recipient", dto.Recipient, errors.InvalidFormat, "phone number must be in the format +551100000000")
		}
	}

	if dto.Recipient == "" {
		validationErrs.Add("recipient", dto.Recipient, errors.Required, "recipient must not be empty")
	}

	if dto.Message == "" {
		validationErrs.Add("message", dto.Message, errors.Required, nil)
	}

	if dto.Priority == "" {
		validationErrs.Add("priority", dto.Priority, errors.Required, "priority must not be empty")
	}

	duration, err := helper.DurationFromTimestamp(dto.ScheduledAt)
	if err != nil {
		validationErrs.Add("scheduled_at", dto.ScheduledAt, errors.InvalidValue, "scheduled_at must be fe")
	}

	if duration > 24*time.Hour {
		validationErrs.Add("scheduled_at", dto.ScheduledAt, errors.InvalidValue, "scheduled_at must be a date less than 24 hours away")
	}

	if validationErrs.HasErrors() {
		return nil, validationErrs
	}

	createdNotification := NewNotification(
		channel,
		dto.Recipient,
		dto.Message,
		dto.Data,
		currentTime,
		dto.ScheduledAt,
	)

	return createdNotification, nil
}
