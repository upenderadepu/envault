package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bhartiyaanshul/envault/internal/config"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter applies per-IP rate limiting with different limits for auth, write, and read routes.
func RateLimiter(cfg config.RateConfig) func(http.Handler) http.Handler {
	var mu sync.Mutex
	limiters := make(map[string]*ipLimiter)

	// Background cleanup every 5 minutes
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			mu.Lock()
			for ip, l := range limiters {
				if time.Since(l.lastSeen) > 10*time.Minute {
					delete(limiters, ip)
				}
			}
			mu.Unlock()
		}
	}()

	getLimiter := func(ip string, rps int) *rate.Limiter {
		mu.Lock()
		defer mu.Unlock()

		key := ip + ":" + string(rune(rps))
		if l, exists := limiters[key]; exists {
			l.lastSeen = time.Now()
			return l.limiter
		}

		l := &ipLimiter{
			limiter:  rate.NewLimiter(rate.Limit(float64(rps)/60.0), rps/3+1),
			lastSeen: time.Now(),
		}
		limiters[key] = l
		return l.limiter
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)

			// Determine rate limit based on route/method
			rps := cfg.Read
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
				rps = cfg.Write
			}
			if strings.Contains(r.URL.Path, "/auth") {
				rps = cfg.Auth
			}

			limiter := getLimiter(ip, rps)
			if !limiter.Allow() {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", "60")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]string{"error": "rate limit exceeded"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.SplitN(xff, ",", 2)
		return strings.TrimSpace(parts[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
