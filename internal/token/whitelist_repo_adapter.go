package token

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const keyPrefix = "token-wl:"

type RedisWhitelistRepository struct {
	client  *redis.Client
	expTime time.Duration
}

func NewRedisWhitelistRepository(client *redis.Client, tokenExpirationTime time.Duration) *RedisWhitelistRepository {
	return &RedisWhitelistRepository{
		client:  client,
		expTime: tokenExpirationTime,
	}
}

func (r *RedisWhitelistRepository) Add(ctx context.Context, id uuid.UUID) error {
	status := r.client.SetEx(ctx, keyPrefix+id.String(), "", r.expTime)

	return status.Err()
}

func (r *RedisWhitelistRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	status := r.client.Exists(ctx, keyPrefix+id.String())
	val, err := status.Result()
	if err != nil {
		return false, err
	}

	return val != 0, nil
}

func (r *RedisWhitelistRepository) Remove(ctx context.Context, id uuid.UUID) error {
	status := r.client.Del(ctx, keyPrefix+id.String())

	return status.Err()
}
