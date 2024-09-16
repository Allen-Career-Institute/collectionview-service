package cache

import (
	"context"
	"time"
)

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	GetByPattern(ctx context.Context, pattern string, count int) ([]string, error)
	SetMultiple(ctx context.Context, dataMap map[string]interface{}, ttl time.Duration) (int, error)
	Delete(ctx context.Context, key string) (int64, error)
	DeleteByPattern(ctx context.Context, pattern string) error
}
