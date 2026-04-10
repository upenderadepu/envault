package middleware

import (
	"net/http"
	"strings"

	"github.com/bhartiyaanshul/envault/internal/config"
)

// CORSHandler sets CORS headers and handles preflight OPTIONS requests.
func CORSHandler(cfg config.CORSConfig) func(http.Handler) http.Handler {
	allowedSet := make(map[string]bool, len(cfg.AllowedOrigins))
	for _, o := range cfg.AllowedOrigins {
		allowedSet[o] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if origin != "" && allowedSet[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.Header().Set("Vary", "Origin")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// originAllowed checks if an origin matches any pattern (supports trailing wildcard).
func originAllowed(origin string, patterns []string) bool {
	for _, p := range patterns {
		if p == "*" || p == origin {
			return true
		}
		if strings.HasSuffix(p, "*") && strings.HasPrefix(origin, p[:len(p)-1]) {
			return true
		}
	}
	return false
}
