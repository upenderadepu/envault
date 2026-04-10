package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bhartiyaanshul/envault/internal/config"
	"github.com/bhartiyaanshul/envault/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type ctxKeyUser struct{}

// JWKS response structures
type jwksResponse struct {
	Keys []jwkKey `json:"keys"`
}

type jwkKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type keyCache struct {
	mu   sync.RWMutex
	keys map[string]*rsa.PublicKey
}

// JWTValidator validates Supabase JWTs using cached JWKS public keys.
func JWTValidator(cfg config.AuthConfig) func(http.Handler) http.Handler {
	cache := &keyCache{keys: make(map[string]*rsa.PublicKey)}

	// Skip routes that don't need auth
	skipPaths := map[string]bool{
		"/healthz": true,
		"/readyz":  true,
		"/metrics": true,
	}

	// Fetch JWKS on startup
	if cfg.JWKSURL != "" {
		if err := fetchJWKS(cfg.JWKSURL, cache); err != nil {
			log.Warn().Err(err).Msg("failed to fetch JWKS on startup, will retry")
		}

		// Background refresh every hour
		go func() {
			ticker := time.NewTicker(1 * time.Hour)
			defer ticker.Stop()
			for range ticker.C {
				if err := fetchJWKS(cfg.JWKSURL, cache); err != nil {
					log.Error().Err(err).Msg("JWKS refresh failed")
				} else {
					log.Debug().Msg("JWKS keys refreshed")
				}
			}
		}()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for health/metrics routes
			if skipPaths[r.URL.Path] {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				respondUnauthorized(w, "missing or invalid authorization header")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Parse and validate token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				kid, ok := token.Header["kid"].(string)
				if !ok {
					return nil, fmt.Errorf("missing kid in token header")
				}

				cache.mu.RLock()
				key, exists := cache.keys[kid]
				cache.mu.RUnlock()

				if !exists {
					return nil, fmt.Errorf("unknown key id: %s", kid)
				}
				return key, nil
			},
				jwt.WithValidMethods([]string{"RS256"}),
				jwt.WithIssuer(cfg.JWTIssuer),
			)

			if err != nil || !token.Valid {
				respondUnauthorized(w, "invalid token")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				respondUnauthorized(w, "invalid claims")
				return
			}

			// Validate audience if configured
			if cfg.JWTAudience != "" {
				aud, _ := claims["aud"].(string)
				if aud != cfg.JWTAudience {
					respondUnauthorized(w, "invalid audience")
					return
				}
			}

			sub, _ := claims["sub"].(string)
			email, _ := claims["email"].(string)

			if sub == "" {
				respondUnauthorized(w, "missing subject claim")
				return
			}

			user := &models.User{
				SupabaseUID: sub,
				Email:       email,
			}

			ctx := context.WithValue(r.Context(), ctxKeyUser{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext retrieves the authenticated user identity from context.
func GetUserFromContext(ctx context.Context) *models.User {
	if user, ok := ctx.Value(ctxKeyUser{}).(*models.User); ok {
		return user
	}
	return nil
}

func fetchJWKS(url string, cache *keyCache) error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	var jwks jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return fmt.Errorf("decode JWKS: %w", err)
	}

	newKeys := make(map[string]*rsa.PublicKey)
	for _, key := range jwks.Keys {
		if key.Kty != "RSA" {
			continue
		}

		nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
		if err != nil {
			continue
		}
		eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
		if err != nil {
			continue
		}

		n := new(big.Int).SetBytes(nBytes)
		e := int(new(big.Int).SetBytes(eBytes).Int64())

		newKeys[key.Kid] = &rsa.PublicKey{N: n, E: e}
	}

	cache.mu.Lock()
	cache.keys = newKeys
	cache.mu.Unlock()

	log.Debug().Int("key_count", len(newKeys)).Msg("JWKS keys loaded")
	return nil
}

func respondUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
