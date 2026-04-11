package server

import (
	"github.com/bhartiyaanshul/envault/internal/config"
	"github.com/bhartiyaanshul/envault/internal/repository"
	"github.com/bhartiyaanshul/envault/internal/server/handlers"
	mw "github.com/bhartiyaanshul/envault/internal/server/middleware"
	"github.com/bhartiyaanshul/envault/internal/service"
	"github.com/bhartiyaanshul/envault/internal/vault"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

type RouterDeps struct {
	Config         *config.Config
	DB             *gorm.DB
	VaultSvc       *vault.VaultService
	ProjectRepo    *repository.ProjectRepository
	UserRepo       *repository.UserRepository
	EnvRepo        *repository.EnvironmentRepository
	MemberRepo     *repository.TeamMemberRepository
	SecretMetaRepo *repository.SecretMetadataRepository
	AuditRepo      *repository.AuditLogRepository
}

func NewRouter(deps RouterDeps) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware (applied to every request)
	r.Use(mw.RequestID)
	r.Use(mw.StructuredLogger)
	r.Use(mw.Recoverer)
	r.Use(mw.RateLimiter(deps.Config.Rate))
	r.Use(mw.CORSHandler(deps.Config.CORS))
	r.Use(chimw.Compress(5))

	// Build services
	projectSvc := service.NewProjectService(
		deps.ProjectRepo, deps.EnvRepo, deps.MemberRepo,
		deps.AuditRepo, deps.SecretMetaRepo, deps.VaultSvc,
		deps.Config.Vault.MountPrefix,
	)
	secretSvc := service.NewSecretService(
		deps.SecretMetaRepo, deps.EnvRepo, deps.ProjectRepo,
		deps.AuditRepo, deps.VaultSvc,
	)
	memberSvc := service.NewMemberService(
		deps.MemberRepo, deps.UserRepo, deps.ProjectRepo,
		deps.AuditRepo, deps.VaultSvc,
	)
	auditSvc := service.NewAuditService(deps.AuditRepo, deps.ProjectRepo)

	// Build handlers
	healthHandler := handlers.NewHealthHandler(deps.DB)
	projectHandler := handlers.NewProjectHandler(projectSvc, deps.UserRepo)
	secretHandler := handlers.NewSecretHandler(secretSvc, deps.UserRepo)
	memberHandler := handlers.NewMemberHandler(memberSvc)
	auditHandler := handlers.NewAuditHandler(auditSvc)

	// Health routes (no auth required)
	r.Get("/healthz", healthHandler.Healthz)
	r.Get("/readyz", healthHandler.Readyz)
	r.Handle("/metrics", promhttp.Handler())

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		jwtMw := mw.JWTValidator(deps.Config.Auth)
		r.Use(mw.CLIOrJWTAuth(deps.Config.Auth, deps.VaultSvc.Client(), deps.MemberRepo, deps.UserRepo, jwtMw))

		// Accept invite (no RBAC — user just needs to be authenticated)
		r.Post("/invite/accept", memberHandler.AcceptInvite)

		// Project CRUD (no RBAC — Create and List don't have a slug)
		r.Post("/projects", projectHandler.Create)
		r.Get("/projects", projectHandler.List)

		// Project-scoped routes (RBAC enforced)
		r.Route("/projects/{slug}", func(r chi.Router) {
			r.Use(mw.RBACEnforcer(deps.ProjectRepo, deps.MemberRepo, deps.UserRepo))

			r.Get("/", projectHandler.Get)
			r.Delete("/", projectHandler.Delete)

			// Secrets
			r.Get("/secrets", secretHandler.List)
			r.Post("/secrets", secretHandler.Set)
			r.Post("/secrets/bulk", secretHandler.BulkSet)
			r.Get("/secrets/{key}", secretHandler.Get)
			r.Delete("/secrets/{key}", secretHandler.Delete)

			// Members
			r.Get("/members", memberHandler.List)
			r.Post("/members", memberHandler.Add)
			r.Delete("/members/{id}", memberHandler.Remove)

			// Rotate credentials
			r.Post("/rotate", memberHandler.Rotate)

			// Audit logs
			r.Get("/audit", auditHandler.List)
		})
	})

	return r
}
