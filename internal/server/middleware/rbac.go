package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bhartiyaanshul/envault/internal/models"
	"github.com/bhartiyaanshul/envault/internal/repository"
	"github.com/go-chi/chi/v5"
)

type ctxKeyProject struct{}
type ctxKeyMember struct{}

// RBACEnforcer loads the project, verifies user membership, and enforces role permissions.
// Stores the project and team member in context for handlers.
func RBACEnforcer(
	projectRepo *repository.ProjectRepository,
	memberRepo *repository.TeamMemberRepository,
	userRepo *repository.UserRepository,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slug := chi.URLParam(r, "slug")
			if slug == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Load project
			project, err := projectRepo.FindBySlug(slug)
			if err != nil {
				respondForbidden(w, "project not found")
				return
			}

			// Get authenticated user from JWT context
			jwtUser := GetUserFromContext(r.Context())
			if jwtUser == nil {
				respondForbidden(w, "user not authenticated")
				return
			}

			// Find or create user in DB
			user, err := userRepo.FindOrCreate(jwtUser.SupabaseUID, jwtUser.Email)
			if err != nil {
				respondForbidden(w, "user lookup failed")
				return
			}

			// Check if user is owner or active team member
			var member *models.TeamMember

			if project.OwnerID == user.ID {
				// Owner has implicit admin access
				member = &models.TeamMember{
					ProjectID: project.ID,
					UserID:    user.ID,
					Role:      "admin",
					IsActive:  true,
				}
			} else {
				member, err = memberRepo.FindByProjectAndUser(project.ID, user.ID)
				if err != nil || !member.IsActive {
					respondForbidden(w, "access denied")
					return
				}
			}

			// Enforce role-based permissions
			if !isAllowed(member.Role, r.Method, r.URL.Path) {
				respondForbidden(w, "insufficient permissions")
				return
			}

			ctx := context.WithValue(r.Context(), ctxKeyProject{}, project)
			ctx = context.WithValue(ctx, ctxKeyMember{}, member)
			// Store the resolved DB user back in context
			ctx = context.WithValue(ctx, ctxKeyUser{}, &models.User{
				ID:          user.ID,
				SupabaseUID: user.SupabaseUID,
				Email:       user.Email,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetProjectFromContext(ctx context.Context) *models.Project {
	if p, ok := ctx.Value(ctxKeyProject{}).(*models.Project); ok {
		return p
	}
	return nil
}

func GetMemberFromContext(ctx context.Context) *models.TeamMember {
	if m, ok := ctx.Value(ctxKeyMember{}).(*models.TeamMember); ok {
		return m
	}
	return nil
}

// isAllowed checks if the role has permission for the given method/path.
//
// Permission matrix:
//
//	admin:     full access
//	developer: read/write secrets (non-prod), read members/audit
//	ci:        read secrets only
func isAllowed(role, method, path string) bool {
	switch role {
	case "admin":
		return true

	case "developer":
		// Can read anything
		if method == http.MethodGet {
			return true
		}
		// Can write secrets
		if (method == http.MethodPost || method == http.MethodPut) && strings.Contains(path, "/secrets") {
			return true
		}
		// Cannot delete secrets, manage members, or delete project
		return false

	case "ci":
		// Read-only access
		return method == http.MethodGet

	default:
		return false
	}
}

func respondForbidden(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
