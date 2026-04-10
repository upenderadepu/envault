package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Vault    VaultConfig
	Auth     AuthConfig
	CORS     CORSConfig
	Rate     RateConfig
}

type ServerConfig struct {
	Port     string
	Host     string
	LogLevel string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type VaultConfig struct {
	Addr        string
	Token       string
	MountPrefix string
}

type AuthConfig struct {
	JWKSURL     string
	JWTIssuer   string
	JWTAudience string
}

type CORSConfig struct {
	AllowedOrigins []string
}

type RateConfig struct {
	Auth  int
	Write int
	Read  int
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:     envOrDefault("SERVER_PORT", "8080"),
			Host:     envOrDefault("SERVER_HOST", "0.0.0.0"),
			LogLevel: envOrDefault("LOG_LEVEL", "info"),
		},
		Database: DatabaseConfig{
			Host:     envOrDefault("DATABASE_HOST", "localhost"),
			Port:     envOrDefault("DATABASE_PORT", "5432"),
			User:     envOrDefault("DATABASE_USER", "envault"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			Name:     envOrDefault("DATABASE_NAME", "envault"),
			SSLMode:  envOrDefault("DATABASE_SSLMODE", "disable"),
		},
		Vault: VaultConfig{
			Addr:        envOrDefault("VAULT_ADDR", "http://localhost:8200"),
			Token:       os.Getenv("VAULT_TOKEN"),
			MountPrefix: envOrDefault("VAULT_MOUNT_PREFIX", "envault"),
		},
		Auth: AuthConfig{
			JWKSURL:     os.Getenv("JWKS_URL"),
			JWTIssuer:   os.Getenv("JWT_ISSUER"),
			JWTAudience: os.Getenv("JWT_AUDIENCE"),
		},
		CORS: CORSConfig{
			AllowedOrigins: splitAndTrim(envOrDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),
		},
		Rate: RateConfig{
			Auth:  envOrDefaultInt("RATE_LIMIT_AUTH", 10),
			Write: envOrDefaultInt("RATE_LIMIT_WRITE", 30),
			Read:  envOrDefaultInt("RATE_LIMIT_READ", 100),
		},
	}

	if cfg.Vault.Addr == "" {
		return nil, fmt.Errorf("VAULT_ADDR is required")
	}

	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envOrDefaultInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
