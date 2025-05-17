package helper

import (
	"errors"
	"time"
)

func DurationFromTimestamp(scheduled *time.Time) (time.Duration, error) {
	now := time.Now()
	if scheduled.Before(now) {
		return 0, errors.New("scheduled time is in the past")
	}

	return scheduled.Sub(now), nil
}
