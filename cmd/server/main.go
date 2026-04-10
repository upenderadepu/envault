package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/bhartiyaanshul/envault/internal/config"
	"github.com/bhartiyaanshul/envault/internal/db"
	"github.com/bhartiyaanshul/envault/internal/repository"
	"github.com/bhartiyaanshul/envault/internal/server"
	"github.com/bhartiyaanshul/envault/internal/vault"
	"github.com/bhartiyaanshul/envault/migrations"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configure zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	// Set log level
	level, err := zerolog.ParseLevel(cfg.Server.LogLevel)
	if err == nil {
		zerolog.SetGlobalLevel(level)
	}

	// Connect to database
	database, err := db.Connect(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	// Run migrations
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get sql.DB")
	}
	if err := runMigrations(sqlDB); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	// Initialize Vault service
	vaultSvc, err := vault.NewVaultService(cfg.Vault)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize vault service")
	}

	// Start token renewal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	vaultSvc.StartTokenRenewal(ctx)

	// Initialize repositories
	userRepo := repository.NewUserRepository(database)
	projectRepo := repository.NewProjectRepository(database)
	envRepo := repository.NewEnvironmentRepository(database)
	memberRepo := repository.NewTeamMemberRepository(database)
	secretMetaRepo := repository.NewSecretMetadataRepository(database)
	auditRepo := repository.NewAuditLogRepository(database)

	// Build router
	router := server.NewRouter(server.RouterDeps{
		Config:         cfg,
		DB:             database,
		VaultSvc:       vaultSvc,
		ProjectRepo:    projectRepo,
		UserRepo:       userRepo,
		EnvRepo:        envRepo,
		MemberRepo:     memberRepo,
		SecretMetaRepo: secretMetaRepo,
		AuditRepo:      auditRepo,
	})

	// Start HTTP server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := server.Start(router, addr); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}

func runMigrations(sqlDB *sql.DB) error {
	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(sqlDB, ".")
}
