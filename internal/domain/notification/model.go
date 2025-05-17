package notification

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
)

type ID uuid.UUID

type Notification struct {
	ID          ID                  `db:"id"`
	Channel     NotificationChannel `db:"channel"`
	Recipient   string              `db:"recipient"`
	Message     string              `db:"message"`
	Data        string              `db:"data"`
	Status      NotificationStatus  `db:"status"`
	CreatedAt   time.Time           `db:"created_at"`
	UpdatedAt   time.Time           `db:"updated_at"`
	ScheduledAt *time.Time          `db:"scheduled_at"`
	DeletedAt   *time.Time          `db:"deleted_at"`
}

func (id *ID) Scan(src interface{}) error {
	var u uuid.UUID
	if err := u.Scan(src); err != nil {
		return err
	}

	*id = ID(u)
	return nil
}

func (id ID) Value() (driver.Value, error) {
	return uuid.UUID(id).String(), nil
}

func (id ID) String() string {
	return uuid.UUID(id).String()
}
