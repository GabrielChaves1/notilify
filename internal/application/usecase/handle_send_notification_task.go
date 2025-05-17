package usecase

import (
	apperrors "GabrielChaves1/notilify/internal/application/error"
	"GabrielChaves1/notilify/internal/application/manager"
	"GabrielChaves1/notilify/internal/domain/notification"
	"context"

	"github.com/sirupsen/logrus"
)

type HandleSendNotificationTask struct {
	notificationRepo notification.Repository
	channelManager   *manager.ChannelManager
	logger           *logrus.Entry
}

func NewHandleSendNotificationTask(notificationRepo notification.Repository, channelManager *manager.ChannelManager, logger *logrus.Entry) *HandleSendNotificationTask {
	return &HandleSendNotificationTask{
		notificationRepo: notificationRepo,
		channelManager:   channelManager,
		logger:           logger,
	}
}

func (uc *HandleSendNotificationTask) Execute(ctx context.Context, notificationStr notification.ID) error {
	uc.logger.Info("Fetching notification in database")

	rawNotification, err := uc.notificationRepo.GetByID(ctx, notificationStr)
	if err != nil {
		return apperrors.NewInternalServerError(ctx, err.Error())
	}

	uc.logger.Info("Sending message to communication channel")

	ch := uc.channelManager.GetChannel(rawNotification.Channel)
	if err := ch.Send(ctx, rawNotification); err != nil {
		return apperrors.NewInternalServerError(ctx, err.Error())
	}

	rawNotification.Status = notification.Delivered

	if err := uc.notificationRepo.Save(ctx, rawNotification); err != nil {
		return apperrors.NewInternalServerError(ctx, err.Error())
	}

	return nil
}
