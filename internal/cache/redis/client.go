package redis

import (
	"GabrielChaves1/notilify/internal/cache"
	"context"
	"strings"
	"time"

	rediscache "github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *rediscache.Client
}

func NewRedisClient(addr string) cache.Repository {
	rdb := rediscache.NewClient(&rediscache.Options{
		Addr: strings.TrimPrefix(addr, "redis://"),
	})

	return &RedisCache{client: rdb}
}

func (r *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := r.client.Get(ctx, key).Bytes()
	if err == rediscache.Nil {
		return nil, nil
	}

	return val, err
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) SubscribeExpired(ctx context.Context, handler func(key string)) error {
	pubsub := r.client.PSubscribe(ctx, "__keyevent@0__:expired")

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}

			//---- log
			continue
		}

		key := msg.Payload
		if strings.HasPrefix(key, "notification:") {
			handler(key)
		}
	}
}
