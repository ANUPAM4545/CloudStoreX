package cache

import (
	"context"
	"errors"
	"sync"
)

type AICache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
}

type inMemoryCache struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewCache() AICache {
	return &inMemoryCache{
		store: make(map[string]string),
	}
}

func (c *inMemoryCache) Get(ctx context.Context, key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, exists := c.store[key]
	if !exists {
		return "", errors.New("cache miss")
	}
	return val, nil
}

func (c *inMemoryCache) Set(ctx context.Context, key string, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
	return nil
}
