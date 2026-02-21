package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// DefaultTTL is the default cache TTL (5 minutes)
	DefaultTTL = 5 * time.Minute
)

// Cache represents a cache interface
type Cache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

// RedisCache implements Cache using Redis
type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisCache creates a new Redis cache
func NewRedisCache(client *redis.Client, ttl time.Duration) *RedisCache {
	if ttl == 0 {
		ttl = DefaultTTL
	}
	return &RedisCache{
		client: client,
		ttl:    ttl,
	}
}

// Get retrieves a value from cache
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// Set stores a value in cache
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if ttl == 0 {
		ttl = c.ttl
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

// Delete removes a key from cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Exists checks if a key exists
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	return n > 0, err
}

// NoOpCache is a no-op cache for when Redis is not available
type NoOpCache struct{}

// NewNoOpCache creates a new no-op cache
func NewNoOpCache() *NoOpCache {
	return &NoOpCache{}
}

// Get always returns an error
func (c *NoOpCache) Get(ctx context.Context, key string, dest interface{}) error {
	return ErrCacheMiss
}

// Set does nothing
func (c *NoOpCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}

// Delete does nothing
func (c *NoOpCache) Delete(ctx context.Context, key string) error {
	return nil
}

// Exists always returns false
func (c *NoOpCache) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

// ErrCacheMiss is returned when cache key is not found
var ErrCacheMiss = &CacheError{Message: "cache miss"}

// CacheError represents a cache error
type CacheError struct {
	Message string
}

func (e *CacheError) Error() string {
	return e.Message
}

// Key generators for consistent key naming
func ONUListKey(oltID string, boardID, ponID int) string {
	return "onu_list:" + oltID + ":" + intToString(boardID) + ":" + intToString(ponID)
}

func ONUDetailKey(oltID string, boardID, ponID, onuID int) string {
	return "onu_detail:" + oltID + ":" + intToString(boardID) + ":" + intToString(ponID) + ":" + intToString(onuID)
}

func EmptySlotsKey(oltID string, boardID, ponID int) string {
	return "empty_slots:" + oltID + ":" + intToString(boardID) + ":" + intToString(ponID)
}

func intToString(i int) string {
	return string(rune('0' + i))
}
