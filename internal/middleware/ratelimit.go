package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// RateLimiter mengimplementasikan pembatasan rate (rate limiting) per IP
type RateLimiter struct {
	requests map[string]*clientInfo
	mu       sync.RWMutex
	rate     int           // jumlah request per jendela waktu
	window   time.Duration // jendela waktu
}

type clientInfo struct {
	count     int
	windowEnd time.Time
}

// NewRateLimiter membuat rate limiter baru
// rate: jumlah request yang diizinkan per jendela waktu
// window: durasi jendela waktu
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*clientInfo),
		rate:     rate,
		window:   window,
	}

	// Bersihkan entri lama setiap menit
	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, info := range rl.requests {
			if now.After(info.windowEnd) {
				delete(rl.requests, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getRealIP(r)

		rl.mu.Lock()
		defer rl.mu.Unlock()

		now := time.Now()
		info, exists := rl.requests[ip]

		if !exists || now.After(info.windowEnd) {
			// Jendela waktu baru
			rl.requests[ip] = &clientInfo{
				count:     1,
				windowEnd: now.Add(rl.window),
			}
			next.ServeHTTP(w, r)
			return
		}

		// Jendela waktu yang sama
		if info.count >= rl.rate {
			log.Warn().Str("ip", ip).Msg("Rate limit exceeded")
			http.Error(w, `{"code":429,"status":"ERROR","message":"rate limit exceeded"}`, http.StatusTooManyRequests)
			return
		}

		info.count++
		next.ServeHTTP(w, r)
	})
}

func getRealIP(r *http.Request) string {
	// Periksa header X-Forwarded-For
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	// Periksa header X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// Kembali ke RemoteAddr jika header di atas tidak ada
	return r.RemoteAddr
}
