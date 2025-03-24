package localmemory

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/personal/task-management/pkg/cache"
)

// Should singleton
// NewCache initializes a new Cache instance with cleanup interval
func NewCache(cleanupInterval time.Duration) (cache.Cache, error) {
	if cleanupInterval <= 0 {
		return nil, cache.ErrInvalidParams
	}

	c := &localMemory{
		store:    sync.Map{},
		mu:       sync.Mutex{},
		ticker:   time.NewTicker(cleanupInterval),
		stopChan: make(chan struct{}),
	}

	c.wg.Add(1)
	go c.startCleanupRoutine()

	return c, nil
}

type cacheItem struct {
	value      any
	expireTime *time.Time
}

func (item *cacheItem) isExpired() bool {
	return item.expireTime != nil && time.Now().After(*item.expireTime)
}

type localMemory struct {
	store    sync.Map
	mu       sync.Mutex
	ticker   *time.Ticker
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func (c *localMemory) Set(ctx context.Context, key, value any) error {
	if key == nil || value == nil {
		return cache.ErrInvalidParams
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		defaultExp := time.Now().Add(5 * time.Minute)
		c.mu.Lock()
		defer c.mu.Unlock()
		c.store.Store(key, cacheItem{value: value, expireTime: &defaultExp})
		return nil
	}
}

func (c *localMemory) SetWithExpire(ctx context.Context, key, value any, expire time.Duration) error {
	if key == nil || value == nil || expire <= 0 {
		return cache.ErrInvalidParams
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		expiration := time.Now().Add(expire)
		c.mu.Lock()
		defer c.mu.Unlock()
		c.store.Store(key, cacheItem{value: value, expireTime: &expiration})
		return nil
	}
}

func (c *localMemory) Get(ctx context.Context, key any) (any, error) {
	if key == nil {
		return nil, cache.ErrInvalidParams
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		c.mu.Lock()
		defer c.mu.Unlock()
		item, ok := c.store.Load(key)
		if !ok {
			return nil, cache.ErrKeyNotFound
		}
		if item.(cacheItem).expireTime != nil && time.Now().After(*item.(cacheItem).expireTime) {
			c.store.Delete(key)
			return nil, cache.ErrKeyExpired
		}
		return item.(cacheItem).value, nil
	}
}

func (c *localMemory) Update(ctx context.Context, key, value any) error {
	if key == nil || value == nil {
		return cache.ErrInvalidParams
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		c.mu.Lock()
		defer c.mu.Unlock()
		item, ok := c.store.Load(key)
		if !ok {
			return cache.ErrKeyNotFound
		}

		c.store.Store(key, cacheItem{value: value, expireTime: item.(cacheItem).expireTime})
		return nil
	}
}

func (c *localMemory) Delete(ctx context.Context, key any) error {
	if key == nil {
		return cache.ErrInvalidParams
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		c.mu.Lock()
		defer c.mu.Unlock()
		c.store.Delete(key)
		return nil
	}
}

func (c *localMemory) DeleteByPrefix(ctx context.Context, prefix string) error {
	if prefix == "" {
		return cache.ErrInvalidParams
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		c.mu.Lock()
		defer c.mu.Unlock()

		// Collect keys to delete
		var keysToDelete []any
		c.store.Range(func(key, _ any) bool {
			// Check if key is a string and has the prefix
			if keyStr, ok := key.(string); ok && len(keyStr) >= len(prefix) && strings.HasPrefix(keyStr, prefix) {
				keysToDelete = append(keysToDelete, key)
			}
			return true
		})

		// Delete collected keys
		for _, key := range keysToDelete {
			c.store.Delete(key)
		}

		return nil
	}
}

func (c *localMemory) Close() error {
	close(c.stopChan)
	c.ticker.Stop()
	c.wg.Wait()
	return nil
}

func (c *localMemory) startCleanupRoutine() {
	defer c.wg.Done()

	for {
		select {
		case <-c.ticker.C:
			c.cleanupExpired()
		case <-c.stopChan:
			return
		}
	}
}

func (c *localMemory) cleanupExpired() {
	c.store.Range(func(key, value any) bool {
		if item, ok := value.(*cacheItem); ok && item.isExpired() {
			c.store.Delete(key)
		}
		return true
	})
}

// For singleton usage (optional)
var (
	instance cache.Cache
	once     sync.Once
)

// GetInstance returns a singleton cache instance with default cleanup interval
func GetInstance() (cache.Cache, error) {
	var err error
	once.Do(func() {
		instance, err = NewCache(time.Minute)
	})
	return instance, err
}
