package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/bhartiyaanshul/envault/internal/models"
	"github.com/bhartiyaanshul/envault/internal/repository"
	vault "github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/rs/zerolog/log"
)

// CLIOrJWTAuth tries JWT auth first, falls back to Vault token auth for CLI clients.
// Vault token auth works by looking up the token's accessor via Vault, then finding
// the team member associated with that accessor.
func CLIOrJWTAuth(
	_ interface{},
	vaultClient *vault.Client,
	memberRepo *repository.TeamMemberRepository,
	userRepo *repository.UserRepository,
	jwtMiddleware func(http.Handler) http.Handler,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondUnauthorized(w, "missing authorization header")
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			// If token starts with "hvs." or "s.", it's a Vault token
			if strings.HasPrefix(token, "hvs.") || strings.HasPrefix(token, "s.") {
				handleVaultTokenAuth(w, r, next, vaultClient, memberRepo, userRepo, token)
				return
			}

			// Otherwise, try JWT auth (Supabase)
			jwtMiddleware(next).ServeHTTP(w, r)
		})
	}
}

func handleVaultTokenAuth(
	w http.ResponseWriter,
	r *http.Request,
	next http.Handler,
	vaultClient *vault.Client,
	memberRepo *repository.TeamMemberRepository,
	userRepo *repository.UserRepository,
	token string,
) {
	// Look up the token in Vault to get its accessor
	resp, err := vaultClient.Auth.TokenLookUp(r.Context(), schema.TokenLookUpRequest{Token: token})
	if err != nil {
		log.Debug().Err(err).Msg("vault token lookup failed")
		respondUnauthorized(w, "invalid token")
		return
	}

	accessor, ok := resp.Data["accessor"].(string)
	if !ok || accessor == "" {
		respondUnauthorized(w, "invalid token: no accessor")
		return
	}

	// Find the team member by accessor
	member, err := memberRepo.FindByAccessor(accessor)
	if err != nil {
		log.Debug().Err(err).Str("accessor", accessor).Msg("no team member for accessor")
		respondUnauthorized(w, "invalid token")
		return
	}

	// Load the user
	user, err := userRepo.FindByID(member.UserID)
	if err != nil {
		respondUnauthorized(w, "user not found")
		return
	}

	// Set user in context (same as JWT middleware does)
	ctx := context.WithValue(r.Context(), ctxKeyUser{}, &models.User{
		ID:          user.ID,
		SupabaseUID: user.SupabaseUID,
		Email:       user.Email,
	})

	next.ServeHTTP(w, r.WithContext(ctx))
}
