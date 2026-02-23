package middleware

import (
	"net/http"
	"strings"
)

// Middleware CORS untuk permintaan lintas asal (cross-origin)
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			
			// Periksa apakah asal (origin) diizinkan
			allowed := false
			for _, ao := range allowedOrigins {
				if ao == "*" || ao == origin {
					allowed = true
					break
				}
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
			}

			// Tangani permintaan preflight
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// DefaultCORS mengembalikan middleware CORS yang mengizinkan semua asal
func DefaultCORS() func(http.Handler) http.Handler {
	return CORS([]string{"*"})
}

// StrictCORS mengembalikan middleware CORS yang hanya mengizinkan asal tertentu
func StrictCORS(origins string) func(http.Handler) http.Handler {
	allowedOrigins := strings.Split(origins, ",")
	for i, o := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(o)
	}
	return CORS(allowedOrigins)
}
