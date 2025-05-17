package usecase

import (
	appcontext "GabrielChaves1/notilify/internal/application/context"
	"GabrielChaves1/notilify/internal/application/dto/request"
	"GabrielChaves1/notilify/internal/application/dto/response"
	apperrors "GabrielChaves1/notilify/internal/application/error"
	"GabrielChaves1/notilify/internal/application/helper"
	"GabrielChaves1/notilify/internal/cache"
	"GabrielChaves1/notilify/internal/domain/notification"
	"GabrielChaves1/notilify/internal/queue"
	"context"
	"errors"
	"fmt"
)

type CreateNotification struct {
	repo      notification.Repository
	queue     queue.Repository
	cache     cache.Repository
	validator *notification.Validator
}

func NewCreateNotification(
	repo notification.Repository,
	queue queue.Repository,
	cache cache.Repository,
	validator *notification.Validator,
) *CreateNotification {
	return &CreateNotification{
		repo:      repo,
		queue:     queue,
		cache:     cache,
		validator: validator,
	}
}

func (uc *CreateNotification) Execute(ctx context.Context, cmd request.CreateNotificationDTO) (*response.NotificationDTO, error) {
	currentTime, err := appcontext.ExtractCurrentTimeFromContext(ctx)
	if err != nil {
		return nil, apperrors.NewInvalidContextError(ctx, appcontext.CurrentTimeContextKey, err.Error())
	}

	newNotification, err := uc.validator.Validate(cmd, currentTime)
	if err != nil {
		var validationErrs *apperrors.ValidationErrors
		if errors.As(err, &validationErrs) {
			return nil, apperrors.NewValidationError(ctx, validationErrs)
		}
		return nil, err
	}

	err = uc.repo.Save(ctx, newNotification)
	if err != nil {
		return nil, apperrors.NewInternalServerError(ctx, err.Error())
	}

	notificationDTO := &response.NotificationDTO{
		ID:          newNotification.ID.String(),
		Channel:     string(newNotification.Channel),
		Recipient:   newNotification.Recipient,
		Message:     newNotification.Message,
		Data:        newNotification.Data,
		Status:      string(notification.Pending),
		ScheduledAt: newNotification.ScheduledAt,
		CreatedAt:   newNotification.CreatedAt,
		UpdatedAt:   newNotification.UpdatedAt,
	}

	if newNotification.ScheduledAt == nil {
		err = uc.queue.Publish(ctx, queue.NotificationQueue, newNotification.ID.String())
		if err != nil {
			return nil, apperrors.NewInternalServerError(ctx, err.Error())
		}

		return notificationDTO, nil
	}

	dur, err := helper.DurationFromTimestamp(newNotification.ScheduledAt)
	if err != nil {
		return nil, err
	}

	uc.cache.Set(ctx, fmt.Sprintf("notification:%s", newNotification.ID.String()), []byte(""), dur)

	return notificationDTO, nil
}
