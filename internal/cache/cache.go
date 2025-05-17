package cache

import (
	"context"
	"time"
)

type Repository interface {
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	SubscribeExpired(ctx context.Context, handler func(key string)) error
}
