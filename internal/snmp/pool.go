package snmp

import (
	"context"
	"sync"
	"time"

	"github.com/gosnmp/gosnmp"
)

// Pool manages SNMP connections with concurrency limits
type Pool struct {
	maxConcurrent int
	sem           chan struct{}
	mu            sync.Mutex
}

// NewPool creates a new SNMP connection pool
func NewPool(maxConcurrent int) *Pool {
	return &Pool{
		maxConcurrent: maxConcurrent,
		sem:           make(chan struct{}, maxConcurrent),
	}
}

// Query executes an SNMP query with connection pooling
func (p *Pool) Query(ctx context.Context, cfg Config, fn func(*gosnmp.GoSNMP) error) error {
	// Acquire semaphore
	select {
	case p.sem <- struct{}{}:
		defer func() { <-p.sem }()
	case <-ctx.Done():
		return ctx.Err()
	}

	// Create connection (SNMP is UDP, so we create per query)
	client := &gosnmp.GoSNMP{
		Target:    cfg.Host,
		Port:      cfg.Port,
		Community: cfg.Community,
		Version:   gosnmp.Version2c,
		Timeout:   cfg.Timeout,
		Retries:   cfg.Retries,
		MaxOids:   cfg.MaxOids,
	}

	if err := client.Connect(); err != nil {
		return err
	}
	defer client.Conn.Close()

	return fn(client)
}

// Stats returns current pool statistics
func (p *Pool) Stats() PoolStats {
	return PoolStats{
		MaxConcurrent:  p.maxConcurrent,
		ActiveConnections: len(p.sem),
		AvailableSlots: p.maxConcurrent - len(p.sem),
	}
}

// PoolStats represents pool statistics
type PoolStats struct {
	MaxConcurrent    int `json:"max_concurrent"`
	ActiveConnections int `json:"active_connections"`
	AvailableSlots   int `json:"available_slots"`
}

// Global pool instance
var globalPool *Pool
var poolOnce sync.Once

// GetPool returns the global SNMP pool
func GetPool() *Pool {
	poolOnce.Do(func() {
		globalPool = NewPool(100) // Max 100 concurrent SNMP requests
	})
	return globalPool
}

// SetPoolMax sets the max concurrent connections for the global pool
func SetPoolMax(max int) {
	globalPool = NewPool(max)
}

// QueryWithTimeout executes an SNMP query with timeout
func QueryWithTimeout(host string, port uint16, community string, timeout time.Duration, fn func(*gosnmp.GoSNMP) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return GetPool().Query(ctx, Config{
		Host:      host,
		Port:      port,
		Community: community,
		Timeout:   timeout,
		Retries:   1,
		MaxOids:   60,
	}, fn)
}
