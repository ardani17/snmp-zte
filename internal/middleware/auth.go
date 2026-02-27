package middleware

import (
	"net/http"
	"os"
)

// BasicAuth middleware untuk autentikasi sederhana
// Username dan password diambil dari environment variable
func BasicAuth() func(http.Handler) http.Handler {
	username := os.Getenv("AUTH_USER")
	password := os.Getenv("AUTH_PASS")

	// Default credentials kalau env tidak diset (untuk development)
	if username == "" {
		username = "admin"
	}
	if password == "" {
		password = "testing123"
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, p, ok := r.BasicAuth()
			if !ok || u != username || p != password {
				w.Header().Set("WWW-Authenticate", `Basic realm="SNMP-ZTE API"`)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"code":401,"status":"ERROR","message":"Unauthorized - Username atau password salah"}`))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
