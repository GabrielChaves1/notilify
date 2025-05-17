package context

import (
	"context"
	"fmt"
	"time"
)

const (
	CurrentTimeContextKey = "currentTime"
	AgentIDContextKey     = "agentID"
)

func NewContextFieldNotFoundError(fieldName any) error {
	err := fmt.Errorf("couldn't find context field with name %q", fieldName)
	return err
}

func NewContextWithCurrentTime(ctx context.Context, currentTime time.Time) context.Context {
	return context.WithValue(ctx, CurrentTimeContextKey, currentTime)
}

func ExtractCurrentTimeFromContext(ctx context.Context) (time.Time, error) {
	currentTime, ok := ctx.Value(CurrentTimeContextKey).(time.Time)

	if !ok {
		return time.Time{}, NewContextFieldNotFoundError(CurrentTimeContextKey)
	}

	return currentTime, nil
}
