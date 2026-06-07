package shareddomain

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type StorageClient interface {
	GeneratePutURL(ctx context.Context, key, contentType string, ttl time.Duration) (valueobject.URL, error)
	GenerateGetURL(ctx context.Context, key string, ttl time.Duration) (valueobject.URL, error)
	DeleteObject(ctx context.Context, key string) error
}

type URLBuilder interface {
	BuildPublicURL(path string) string
}
