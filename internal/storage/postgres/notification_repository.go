package postgres

import (
	"GabrielChaves1/notilify/internal/domain/notification"
	repository "GabrielChaves1/notilify/internal/storage"
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type NotificationRepository struct {
	db *sqlx.DB
}

func NewNotificationRepository(db *sqlx.DB) notification.Repository {
	return &NotificationRepository{
		db: db,
	}
}

func (r *NotificationRepository) Save(ctx context.Context, notification *notification.Notification) error {
	query := `
		INSERT INTO notifications (id, channel, recipient, message, data, status, created_at, updated_at, scheduled_at)
	VALUES (:id, :channel, :recipient, :message, :data, :status, :created_at, :updated_at, :scheduled_at)
		ON CONFLICT (id) DO UPDATE SET
			channel = :channel,
			recipient = :recipient,
			message = :message,
			data = :data,
			status = :status,
			created_at = :created_at,
			updated_at = :updated_at,
			scheduled_at = :scheduled_at,
			deleted_at = :deleted_at
	`

	_, err := r.db.NamedExecContext(ctx, query, notification)
	if err != nil {
		return err
	}

	return nil
}

func (r *NotificationRepository) GetByID(ctx context.Context, notificationID notification.ID) (*notification.Notification, error) {
	query := `
		SELECT * FROM notifications WHERE id = $1
	`

	var notificationValue notification.Notification

	err := r.db.GetContext(ctx, &notificationValue, query, notificationID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &notificationValue, err
}
