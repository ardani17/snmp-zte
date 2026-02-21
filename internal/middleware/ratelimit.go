package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// RateLimiter implements rate limiting per IP
type RateLimiter struct {
	requests map[string]*clientInfo
	mu       sync.RWMutex
	rate     int           // requests per window
	window   time.Duration // time window
}

type clientInfo struct {
	count     int
	windowEnd time.Time
}

// NewRateLimiter creates a new rate limiter
// rate: number of requests allowed per window
// window: time window duration
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*clientInfo),
		rate:     rate,
		window:   window,
	}

	// Cleanup old entries every minute
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
			// New window
			rl.requests[ip] = &clientInfo{
				count:     1,
				windowEnd: now.Add(rl.window),
			}
			next.ServeHTTP(w, r)
			return
		}

		// Same window
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
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// Fall back to RemoteAddr
	return r.RemoteAddr
}
