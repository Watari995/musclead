package cacheinfra

import (
	"context"
	"errors"
	"time"

	shareddomain "github.com/Watari995/musclead/internal/shared/domain"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", shareddomain.ErrCacheMiss
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}
