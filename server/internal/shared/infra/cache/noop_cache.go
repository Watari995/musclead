package cacheinfra

import (
	"context"
	"time"

	shareddomain "github.com/Watari995/musclead/internal/shared/domain"
)

type NoOpCache struct{}

func NewNoOpCache() *NoOpCache {
	return &NoOpCache{}
}

func (c *NoOpCache) Get(_ context.Context, _ string) (string, error) {
	return "", shareddomain.ErrCacheMiss
}

func (c *NoOpCache) Set(_ context.Context, _ string, _ string, _ time.Duration) error {
	return nil
}

func (c *NoOpCache) DeleteByPattern(ctx context.Context, pattern string) error {
	return nil
}
