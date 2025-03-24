package cache

import (
	"context"
	"errors"
	"time"
)

var (
	ErrKeyNotFound   = errors.New("key not found")
	ErrKeyExpired    = errors.New("key has expired")
	ErrInvalidParams = errors.New("invalid parameters")
)

// Cache defines the interface for cache operations
type Cache interface {
	Set(ctx context.Context, key, value any) error
	SetWithExpire(ctx context.Context, key, value any, expireTime time.Duration) error
	Get(ctx context.Context, key any) (any, error)
	Update(ctx context.Context, key, value any) error
	Delete(ctx context.Context, key any) error
	DeleteByPrefix(ctx context.Context, prefix string) error
	Close() error
}
