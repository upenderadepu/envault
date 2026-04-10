# Contributing to Envault

Thanks for your interest in contributing to Envault! This guide will help you get set up.

## Development Setup

### Prerequisites

- Go 1.22+
- Node.js 18+
- Docker and Docker Compose

### Getting Started

```bash
# Clone the repo
git clone https://github.com/bhartiyaanshul/envault.git
cd envault

# Start PostgreSQL and Vault
docker compose up -d postgres vault

# Enable Vault KV engine
docker compose exec -e VAULT_TOKEN=dev-root-token vault vault secrets enable -path=envault kv-v2

# Configure environment
cp .env.example .env
# Edit .env with your Supabase project credentials

# Start the API server
source .env && go run ./cmd/server/

# In a new terminal — start the dashboard
cd web && npm install && npm run dev

# In a new terminal — build the CLI
make build-cli
```

### Project Layout

| Directory | What it does |
|-----------|-------------|
| `cmd/server/` | API server entry point |
| `cmd/envault/` | CLI entry point |
| `internal/cli/` | CLI commands (cobra) |
| `internal/config/` | Environment-based configuration |
| `internal/db/` | Database connection |
| `internal/models/` | GORM models |
| `internal/repository/` | Data access layer |
| `internal/server/handlers/` | HTTP request handlers |
| `internal/server/middleware/` | JWT, RBAC, rate limiting, CORS, logging |
| `internal/service/` | Business logic |
| `internal/vault/` | Vault client wrapper |
| `migrations/` | SQL migrations (goose, embedded) |
| `web/` | Next.js 14 dashboard |

### Running Tests

```bash
# Go tests
go test ./... -v

# Frontend lint
cd web && npm run lint

# Build check
make build
```

## Making Changes

1. **Fork** the repo and create a branch from `main`
2. Make your changes
3. Run tests: `go test ./...`
4. Run lint: `cd web && npm run lint`
5. Ensure it builds: `make build`
6. Open a pull request

### Code Style

- **Go**: Follow standard Go conventions. Run `go vet ./...`
- **TypeScript/React**: Follow the existing patterns. ESLint is configured
- **Commits**: Use conventional commit messages (`feat:`, `fix:`, `docs:`, etc.)

### Key Design Decisions

- Secret values **never** go in PostgreSQL — only Vault
- The `AuditLog` table has no UPDATE or DELETE — it's append-only
- RBAC is enforced at the middleware level, not in handlers
- The dashboard uses React Query for server state, not Redux/Zustand
- All UI components are built manually (not shadcn CLI) for Tailwind v3 compatibility

## Reporting Issues

Found a bug? Open an issue with:
- Steps to reproduce
- Expected vs actual behavior
- Your environment (OS, Go version, Node version)

## Feature Requests

Open an issue tagged `enhancement` with:
- What problem it solves
- How you envision it working
- Whether you'd like to implement it yourself
