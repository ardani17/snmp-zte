package snmp

import (
	"context"
	"sync"
	"time"

	"github.com/gosnmp/gosnmp"
)

// Pool mengelola koneksi SNMP dengan batasan konkurensi (jumlah pertanyaan bersamaan).
type Pool struct {
	maxConcurrent int
	sem           chan struct{}
	mu            sync.Mutex
}

// NewPool membuat pool koneksi SNMP baru.
func NewPool(maxConcurrent int) *Pool {
	return &Pool{
		maxConcurrent: maxConcurrent,
		sem:           make(chan struct{}, maxConcurrent),
	}
}

// Query menjalankan query SNMP menggunakan pooling koneksi.
func (p *Pool) Query(ctx context.Context, cfg Config, fn func(*gosnmp.GoSNMP) error) error {
	// 1. Mengambil izin (semaphore) dari antrean agar tidak melebihi batas maksimal.
	select {
	case p.sem <- struct{}{}:
		defer func() { <-p.sem }()
	case <-ctx.Done():
		return ctx.Err()
	}

	// 2. Membuat koneksi (SNMP menggunakan UDP, jadi kita buat setiap ada query).
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

// Stats mengembalikan statistik pool saat ini.
func (p *Pool) Stats() PoolStats {
	return PoolStats{
		MaxConcurrent:  p.maxConcurrent,
		ActiveConnections: len(p.sem),
		AvailableSlots: p.maxConcurrent - len(p.sem),
	}
}

// PoolStats merepresentasikan data statistik pool.
type PoolStats struct {
	MaxConcurrent    int `json:"max_concurrent"`
	ActiveConnections int `json:"active_connections"`
	AvailableSlots   int `json:"available_slots"`
}

// Instance pool global.
var globalPool *Pool
var poolOnce sync.Once

// GetPool mengembalikan instance pool SNMP global.
func GetPool() *Pool {
	poolOnce.Do(func() {
		globalPool = NewPool(100) // Maksimal 100 request SNMP bersamaan.
	})
	return globalPool
}

// SetPoolMax mengatur jumlah maksimal koneksi bersamaan untuk pool global.
func SetPoolMax(max int) {
	globalPool = NewPool(max)
}

// QueryWithTimeout menjalankan query SNMP dengan batas waktu (timeout).
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
