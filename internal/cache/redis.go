package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// DefaultTTL adalah durasi penyimpanan cache default (5 menit).
	// Jika data sudah lebih dari 5 menit, program akan bertanya lagi ke OLT.
	DefaultTTL = 5 * time.Minute
)

// Cache mendefinisikan standar (interface) untuk sistem penyimpanan sementara.
type Cache interface {
	Get(ctx context.Context, key string, dest interface{}) error // Mengambil data
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error // Menyimpan data
	Delete(ctx context.Context, key string) error // Menghapus data
	Exists(ctx context.Context, key string) (bool, error) // Cek apakah data ada
}

// RedisCache mengimplementasikan Cache menggunakan Redis
type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisCache membuat cache Redis baru
func NewRedisCache(client *redis.Client, ttl time.Duration) *RedisCache {
	if ttl == 0 {
		ttl = DefaultTTL
	}
	return &RedisCache{
		client: client,
		ttl:    ttl,
	}
}

// Get mengambil nilai dari cache
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// Set menyimpan nilai dalam cache
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

// Delete menghapus kunci dari cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Exists memeriksa apakah kunci ada
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	return n > 0, err
}

// NoOpCache adalah cache no-op (tidak melakukan apa-apa) untuk saat Redis tidak tersedia
type NoOpCache struct{}

// NewNoOpCache membuat cache no-op baru
func NewNoOpCache() *NoOpCache {
	return &NoOpCache{}
}

// Get selalu mengembalikan error
func (c *NoOpCache) Get(ctx context.Context, key string, dest interface{}) error {
	return ErrCacheMiss
}

// Set tidak melakukan apa-apa
func (c *NoOpCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}

// Delete tidak melakukan apa-apa
func (c *NoOpCache) Delete(ctx context.Context, key string) error {
	return nil
}

// Exists selalu mengembalikan false
func (c *NoOpCache) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

// ErrCacheMiss dikembalikan saat kunci cache tidak ditemukan
var ErrCacheMiss = &CacheError{Message: "cache miss"}

// CacheError merepresentasikan error cache
type CacheError struct {
	Message string
}

func (e *CacheError) Error() string {
	return e.Message
}

// Pembuat kunci (key generator) untuk penamaan kunci yang konsisten
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
