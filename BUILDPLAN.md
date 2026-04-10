# Envault — Complete Build Plan

> **Version**: 1.0 | **Created**: 2026-04-10 | **Repo**: github.com/bhartiyaanshul/envault
>
> This document is the single source of truth for building Envault from scratch.
> Hand this to Claude Code in any new session and it has full project context.

---

## Table of Contents

1. [What is Envault](#1-what-is-envault)
2. [Tech Stack](#2-tech-stack)
3. [Prerequisites — Before Writing Code](#3-prerequisites--before-writing-code)
4. [Final Directory Structure](#4-final-directory-structure)
5. [Phase 0: Scaffold & Infrastructure](#phase-0-scaffold--infrastructure-13-files)
6. [Phase 1: Config & GORM Models](#phase-1-config--gorm-models-7-files)
7. [Phase 2: Database & Vault Clients](#phase-2-database--vault-clients-2-files)
8. [Phase 3: Repository Layer](#phase-3-repository-layer-6-files)
9. [Phase 4: Service Layer](#phase-4-service-layer-4-files)
10. [Phase 5: Middleware Stack](#phase-5-middleware-stack-7-files)
11. [Phase 6: HTTP Handlers](#phase-6-http-handlers-6-files)
12. [Phase 7: Router, Server & Entry Points](#phase-7-router-server--entry-points-4-files)
13. [Phase 8: Docker Integration](#phase-8-docker-integration)
14. [Phase 9: CLI Commands](#phase-9-cli-commands-10-files)
15. [Phase 10: Prometheus Metrics](#phase-10-prometheus-metrics)
16. [Phase 11: Next.js Dashboard](#phase-11-nextjs-dashboard)
17. [Phase 12: README & Documentation](#phase-12-readme--documentation)
18. [Phase 13: End-to-End Verification](#phase-13-end-to-end-verification)
19. [Security Rules — Non-Negotiable](#security-rules--non-negotiable)
20. [Key Design Decisions](#key-design-decisions)
21. [Session Handoff Notes](#session-handoff-notes)

---

## 1. What is Envault

Developers waste days onboarding because secrets live in Slack and `.env` files. Envault stores all secrets in HashiCorp Vault, provides a CLI to inject them into any environment, and gives teams audited, role-scoped access.

**Core security invariant**: Secret values NEVER touch the metadata database (Supabase/Postgres). They live ONLY in Vault. The database stores key names, versions, and audit trails.

---

## 2. Tech Stack

| Component | Technology |
|-----------|-----------|
| CLI | Go — cobra, viper, vault-client-go |
| API Server | Go — chi router, zerolog, golang.org/x/time/rate |
| ORM | GORM with pgx/v5 driver |
| Migrations | Goose v3 — embedded SQL files |
| Validation | go-playground/validator/v10 |
| Vault | HashiCorp Vault — KV-v2, AppRole auth, dynamic policies |
| Service Mesh | HashiCorp Consul (optional, service discovery) |
| Metadata DB | PostgreSQL (Supabase-compatible) |
| Auth | JWT validation via Supabase JWKS endpoint |
| Dashboard | Next.js 14, TypeScript, Tailwind CSS, Supabase JS client |
| Observability | zerolog JSON logs + Prometheus metrics |
| Orchestration | Docker Compose |

---

## 3. Prerequisites — Before Writing Code

### 3.1 Software to Install

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.22+ | Backend language |
| Node.js | 18+ | Next.js dashboard |
| Docker Desktop | Latest | Runs Vault, Postgres, Consul |
| Git | Latest | Version control |
| goose CLI | `go install github.com/pressly/goose/v3/cmd/goose@latest` | Run migrations manually |
| golangci-lint | Latest (optional) | Go linter |

### 3.2 Accounts to Create

| Service | Action | Values Needed |
|---------|--------|---------------|
| Supabase | Create free project at supabase.com | Project URL, `anon` key, `service_role` key |
| GitHub | Repo already exists at github.com/bhartiyaanshul/envault | Push access |

### 3.3 Supabase SQL Setup

Run this in Supabase SQL Editor before first migration:
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

### 3.4 Local Environment Setup Steps

```bash
# 1. Clone the repo
git clone https://github.com/bhartiyaanshul/envault.git
cd envault

# 2. Initialize Go module
go mod init github.com/bhartiyaanshul/envault

# 3. Copy env file and fill in Supabase values
cp .env.example .env

# 4. Start infrastructure
docker compose up -d postgres vault

# 5. Initialize Vault (first time only)
docker compose exec vault vault operator init -key-shares=1 -key-threshold=1 -format=json
# SAVE the unseal_key and root_token from output

# 6. Unseal Vault
docker compose exec vault vault operator unseal <unseal_key>

# 7. Set VAULT_TOKEN in .env to the root_token

# 8. Run migrations
make migrate-up

# 9. Build and run
make build
go run ./cmd/server
```

### 3.5 Go Dependencies (exact import paths)

```
github.com/spf13/cobra
github.com/spf13/viper
github.com/go-chi/chi/v5
github.com/go-chi/chi/v5/middleware
github.com/rs/zerolog
github.com/rs/zerolog/log
golang.org/x/time/rate
gorm.io/gorm
gorm.io/driver/postgres
gorm.io/datatypes
github.com/google/uuid
github.com/go-playground/validator/v10
github.com/hashicorp/vault-client-go
github.com/hashicorp/vault-client-go/schema
github.com/pressly/goose/v3
github.com/prometheus/client_golang/prometheus
github.com/prometheus/client_golang/prometheus/promauto
github.com/prometheus/client_golang/prometheus/promhttp
github.com/manifoldco/promptui
github.com/golang-jwt/jwt/v5
```

### 3.6 Internal Package Dependency Graph (Build Order)

```
Layer 0: internal/config          (standalone)
Layer 1: internal/models          (depends on gorm, uuid, datatypes)
Layer 2: internal/db              (depends on config, gorm)
         internal/vault           (depends on config, vault-client-go)
Layer 3: internal/repository      (depends on models, gorm)
Layer 4: internal/service         (depends on repository, vault, models)
Layer 5: internal/server/middleware (depends on config, models, repository)
Layer 6: internal/server/handlers (depends on service, models, validator)
Layer 7: internal/server          (depends on handlers, middleware, chi)
         internal/cli             (depends on config, cobra, viper)
Layer 8: cmd/server/main.go      (depends on config, db, vault, server, goose)
         cmd/envault/main.go     (depends on cli)
```

---

## 4. Final Directory Structure

```
envault/
  go.mod
  go.sum
  .gitignore
  .env.example
  Dockerfile
  docker-compose.yml
  Makefile
  README.md
  BUILDPLAN.md
  deploy/
    vault/
      config.hcl
  migrations/
    00001_create_users.sql
    00002_create_projects.sql
    00003_create_environments.sql
    00004_create_team_members.sql
    00005_create_secret_metadata.sql
    00006_create_audit_logs.sql
  internal/
    config/
      config.go
    models/
      user.go
      project.go
      environment.go
      team_member.go
      secret_metadata.go
      audit_log.go
    db/
      db.go
    vault/
      vault.go
    repository/
      user_repo.go
      project_repo.go
      environment_repo.go
      team_member_repo.go
      secret_metadata_repo.go
      audit_log_repo.go
    service/
      project_service.go
      secret_service.go
      member_service.go
      audit_service.go
    server/
      middleware/
        request_id.go
        logger.go
        recovery.go
        rate_limiter.go
        cors.go
        jwt_validator.go
        rbac.go
      handlers/
        helpers.go
        health.go
        project_handler.go
        secret_handler.go
        member_handler.go
        audit_handler.go
      router.go
      server.go
    cli/
      root.go
      init_cmd.go
      env_pull.go
      env_push.go
      env_list.go
      secret_set.go
      secret_get.go
      secret_delete.go
      onboard.go
      rotate.go
  cmd/
    server/
      main.go
    envault/
      main.go
  web/                            # Next.js 14 dashboard
    package.json
    next.config.js
    tsconfig.json
    tailwind.config.ts
    postcss.config.js
    .env.local.example
    public/
      logo.svg
    src/
      app/
        layout.tsx
        page.tsx
        login/
          page.tsx
        dashboard/
          layout.tsx
          page.tsx
          projects/
            page.tsx
            [slug]/
              page.tsx
              secrets/
                page.tsx
              members/
                page.tsx
              audit/
                page.tsx
              settings/
                page.tsx
          settings/
            page.tsx
      components/
        ui/
          button.tsx
          input.tsx
          card.tsx
          badge.tsx
          dialog.tsx
          dropdown-menu.tsx
          table.tsx
          toast.tsx
          tabs.tsx
          skeleton.tsx
        layout/
          sidebar.tsx
          header.tsx
          nav-links.tsx
        projects/
          project-card.tsx
          create-project-dialog.tsx
        secrets/
          secrets-table.tsx
          set-secret-dialog.tsx
          bulk-import-dialog.tsx
          secret-value-cell.tsx
        members/
          members-table.tsx
          invite-member-dialog.tsx
        audit/
          audit-log-table.tsx
          audit-filters.tsx
      lib/
        supabase/
          client.ts
          server.ts
          middleware.ts
        api.ts
        utils.ts
        types.ts
      hooks/
        use-projects.ts
        use-secrets.ts
        use-members.ts
        use-audit.ts
      middleware.ts
```

---

## Phase 0: Scaffold & Infrastructure (13 files)

### Task List

| # | Task | File | Details |
|---|------|------|---------|
| 0.1 | Initialize git repo | — | `git init` |
| 0.2 | Initialize Go module | `go.mod` | `go mod init github.com/bhartiyaanshul/envault` |
| 0.3 | Create .gitignore | `.gitignore` | Go binaries, `.env`, `/bin/`, `/vendor/`, IDE files, `node_modules/`, `.next/` |
| 0.4 | Create Vault config | `deploy/vault/config.hcl` | File storage at `/vault/data`, TCP listener `0.0.0.0:8200`, TLS disabled, UI enabled, `disable_mlock=true` |
| 0.5 | Create Docker Compose | `docker-compose.yml` | postgres:16-alpine (port 5432), vault:1.15 (port 8200), envault-api (port 8080, commented out initially), envault-web (port 3000, commented out) |
| 0.6 | Create env template | `.env.example` | All env vars: SERVER_*, DATABASE_*, VAULT_*, JWKS_URL, JWT_*, CORS_*, RATE_LIMIT_* |
| 0.7 | Create Makefile | `Makefile` | Targets: build-server, build-cli, build, test, lint, migrate-up/down/status, docker-up/down, vault-init, vault-unseal, clean |
| 0.8 | Create Dockerfile | `Dockerfile` | Multi-stage: golang:1.22-alpine builder -> alpine:3.19 runtime. Copies binary + migrations/ |
| 0.9 | Migration: users | `migrations/00001_create_users.sql` | `CREATE EXTENSION uuid-ossp`. Table: id UUID PK, supabase_uid UNIQUE, email UNIQUE, created_at |
| 0.10 | Migration: projects | `migrations/00002_create_projects.sql` | Soft delete (deleted_at), FK to users, slug unique index |
| 0.11 | Migration: environments | `migrations/00003_create_environments.sql` | FK to projects CASCADE, CHECK (name IN ('development','staging','production')), UNIQUE(project_id, name) |
| 0.12 | Migration: team_members | `migrations/00004_create_team_members.sql` | CHECK (role IN ('admin','developer','ci')), UNIQUE(project_id, user_id) |
| 0.13 | Migration: secret_metadata | `migrations/00005_create_secret_metadata.sql` | UNIQUE(project_id, environment_id, key_name). NO value column. |
| 0.14 | Migration: audit_logs | `migrations/00006_create_audit_logs.sql` | Index on (project_id, created_at DESC). **`REVOKE UPDATE, DELETE ON audit_logs FROM envault;`** |

### Checkpoint Commands
```bash
docker compose up -d postgres vault
docker compose exec vault vault operator init -key-shares=1 -key-threshold=1 -format=json
docker compose exec vault vault operator unseal <key>
make migrate-up
make migrate-status  # all 6 applied
```

### Migration SQL Details

**00001_create_users.sql**:
```sql
-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE users (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    supabase_uid TEXT NOT NULL,
    email        TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX idx_users_supabase_uid ON users(supabase_uid);
CREATE UNIQUE INDEX idx_users_email ON users(email);

-- +goose Down
DROP TABLE IF EXISTS users;
```

**00002_create_projects.sql**:
```sql
-- +goose Up
CREATE TABLE projects (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name             TEXT NOT NULL,
    slug             TEXT NOT NULL,
    vault_mount_path TEXT NOT NULL,
    owner_id         UUID NOT NULL REFERENCES users(id),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at       TIMESTAMPTZ
);
CREATE UNIQUE INDEX idx_projects_slug ON projects(slug);
CREATE INDEX idx_projects_owner_id ON projects(owner_id);
CREATE INDEX idx_projects_deleted_at ON projects(deleted_at);

-- +goose Down
DROP TABLE IF EXISTS projects;
```

**00003_create_environments.sql**:
```sql
-- +goose Up
CREATE TABLE environments (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id    UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name          TEXT NOT NULL CHECK (name IN ('development', 'staging', 'production')),
    is_production BOOLEAN NOT NULL DEFAULT false,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(project_id, name)
);
CREATE INDEX idx_environments_project_id ON environments(project_id);

-- +goose Down
DROP TABLE IF EXISTS environments;
```

**00004_create_team_members.sql**:
```sql
-- +goose Up
CREATE TABLE team_members (
    id                   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id           UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id              UUID NOT NULL REFERENCES users(id),
    role                 TEXT NOT NULL CHECK (role IN ('admin', 'developer', 'ci')),
    vault_policy_name    TEXT,
    vault_token_accessor TEXT,
    is_active            BOOLEAN NOT NULL DEFAULT true,
    invited_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    joined_at            TIMESTAMPTZ,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX idx_team_members_project_user ON team_members(project_id, user_id);
CREATE INDEX idx_team_members_project_id ON team_members(project_id);
CREATE INDEX idx_team_members_user_id ON team_members(user_id);

-- +goose Down
DROP TABLE IF EXISTS team_members;
```

**00005_create_secret_metadata.sql**:
```sql
-- +goose Up
CREATE TABLE secret_metadata (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id       UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    environment_id   UUID NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    key_name         TEXT NOT NULL,
    vault_path       TEXT NOT NULL,
    created_by_id    UUID NOT NULL REFERENCES users(id),
    vault_version    INT NOT NULL DEFAULT 1,
    last_modified_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX idx_secret_metadata_unique ON secret_metadata(project_id, environment_id, key_name);
CREATE INDEX idx_secret_metadata_project_env ON secret_metadata(project_id, environment_id);

-- +goose Down
DROP TABLE IF EXISTS secret_metadata;
```

**00006_create_audit_logs.sql**:
```sql
-- +goose Up
CREATE TABLE audit_logs (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id    UUID NOT NULL REFERENCES projects(id),
    user_id       UUID REFERENCES users(id),
    action        TEXT NOT NULL,
    resource_path TEXT NOT NULL,
    ip_address    TEXT,
    user_agent    TEXT,
    request_id    TEXT,
    metadata      JSONB DEFAULT '{}',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_audit_logs_project_created ON audit_logs(project_id, created_at DESC);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);

REVOKE UPDATE, DELETE ON audit_logs FROM envault;

-- +goose Down
GRANT UPDATE, DELETE ON audit_logs TO envault;
DROP TABLE IF EXISTS audit_logs;
```

---

## Phase 1: Config & GORM Models (7 files)

### Task List

| # | Task | File | Exports |
|---|------|------|---------|
| 1.1 | Config loader | `internal/config/config.go` | `Config`, `ServerConfig`, `DatabaseConfig`, `VaultConfig`, `AuthConfig`, `CORSConfig`, `RateConfig` structs. `Load() (*Config, error)`, `DSN() string`. Pure os.Getenv, no viper. |
| 1.2 | User model | `internal/models/user.go` | `User{ID uuid, SupabaseUID, Email, CreatedAt}` |
| 1.3 | Project model | `internal/models/project.go` | `Project{ID, Name, Slug, VaultMountPath, OwnerID, Owner *User, Environments [], TeamMembers [], CreatedAt, UpdatedAt, DeletedAt gorm.DeletedAt}` |
| 1.4 | Environment model | `internal/models/environment.go` | `Environment{ID, ProjectID, Name, IsProduction, CreatedAt}` |
| 1.5 | TeamMember model | `internal/models/team_member.go` | `TeamMember{ID, ProjectID, UserID, User *User, Role, VaultPolicyName, VaultTokenAccessor, IsActive, InvitedAt, JoinedAt *time.Time}` |
| 1.6 | SecretMetadata model | `internal/models/secret_metadata.go` | `SecretMetadata{ID, ProjectID, EnvironmentID, Environment, KeyName, VaultPath, CreatedByID, VaultVersion, LastModifiedAt}` — **NO Value field** |
| 1.7 | AuditLog model | `internal/models/audit_log.go` | `AuditLog{ID, ProjectID, UserID *uuid.UUID, Action, ResourcePath, IPAddress, UserAgent, RequestID, Metadata datatypes.JSON, CreatedAt}` + Action constants |

### GORM Tag Pattern
All UUID PKs use: `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

### Action Constants (in audit_log.go)
```go
const (
    ActionSecretRead    = "secret.read"
    ActionSecretWrite   = "secret.write"
    ActionSecretDelete  = "secret.delete"
    ActionSecretRotate  = "secret.rotate"
    ActionMemberInvite  = "member.invite"
    ActionMemberRemove  = "member.remove"
    ActionProjectCreate = "project.create"
)
```

### Checkpoint
```bash
go mod tidy
go build ./internal/config/...
go build ./internal/models/...
```

---

## Phase 2: Database & Vault Clients (2 files)

### Task List

| # | Task | File | Details |
|---|------|------|---------|
| 2.1 | DB connection | `internal/db/db.go` | `Connect(cfg DatabaseConfig) (*gorm.DB, error)`. Opens GORM+postgres, configures pool (10 idle, 100 max, 1h lifetime). Does NOT call AutoMigrate. |
| 2.2 | Vault service | `internal/vault/vault.go` | `VaultService` struct. All Vault interactions in one file. |

### VaultService Methods (vault.go is the most complex file)

| Method | Signature | Purpose |
|--------|-----------|---------|
| `NewVaultService` | `(cfg VaultConfig) (*VaultService, error)` | Creates vault client, sets token |
| `StartTokenRenewal` | `(ctx context.Context)` | Background goroutine, renews every 30min via TokenRenewSelf |
| `EnableKV2` | `(ctx, mountPath) error` | Check-then-create. Mount creation is NOT idempotent — must check MountsReadConfiguration first, create only if 404 |
| `WritePolicy` | `(ctx, name, hcl) error` | Writes HCL policy via PoliciesWriteAclPolicy |
| `BuildPolicies` | `(slug, mount, role, envs) string` | Generates HCL. Admin: full `{mount}/data/*` + `{mount}/metadata/*`. Developer: dev+staging only. CI: read-only prod+staging |
| `CreateAppRole` | `(ctx, name, policies, ttl, maxTTL) error` | For CI pipelines |
| `CreateUserToken` | `(ctx, policies, ttl, maxTTL) (token, accessor, error)` | **Returns token+accessor. Store ONLY accessor in DB.** |
| `RevokeTokenByAccessor` | `(ctx, accessor) error` | Revokes without knowing the token |
| `RevokeAllProjectCredentials` | `(ctx, accessors []string) error` | Batch revoke on project delete |
| `ReadSecret` | `(ctx, mount, path) (map[string]interface{}, error)` | KV-v2 read via KvV2Read |
| `WriteSecret` | `(ctx, mount, path, data) error` | **Read-merge-write**: reads existing map, merges new keys (preserves others), writes back |
| `DeleteSecretKey` | `(ctx, mount, path, key) error` | Read, delete key from map, write back |

### KV-v2 Path Convention
- Mount per project: `envault-{slug}` (e.g., `envault-my-project`)
- Environment as path: `development`, `staging`, `production`
- All secrets for an environment = one KV entry at that path

### Token TTLs
- Human (admin/developer): 8h renewable to 24h max
- CI: 1h renewable to 8h max

### Checkpoint
```bash
go mod tidy
go build ./internal/db/...
go build ./internal/vault/...
```

---

## Phase 3: Repository Layer (6 files)

All in `internal/repository/`. Each takes `*gorm.DB` in constructor. No business logic.

### Task List

| # | Task | File | Key Methods |
|---|------|------|-------------|
| 3.1 | User repo | `user_repo.go` | `FindOrCreate(supabaseUID, email)`, `FindBySupabaseUID`, `FindByID`, `FindByEmail` |
| 3.2 | Project repo | `project_repo.go` | `Create`, `FindBySlug` (preloads Envs+Members), `ListForUser` (JOIN team_members OR owner_id), `Delete` (soft) |
| 3.3 | Environment repo | `environment_repo.go` | `CreateBatch`, `FindByProjectID`, `FindByProjectAndName` |
| 3.4 | TeamMember repo | `team_member_repo.go` | `Create`, `FindByProjectAndUser`, `FindActiveByProjectID` (preloads User), `GetAccessorsByProjectID` (returns []string), `Update` |
| 3.5 | SecretMetadata repo | `secret_metadata_repo.go` | `Upsert` (ON CONFLICT project_id,env_id,key_name), `FindByProjectAndEnv`, `FindByKey`, `DeleteByKey` |
| 3.6 | AuditLog repo | `audit_log_repo.go` | `Create`, `FindByProjectID(projectID, limit, offset) ([]AuditLog, int64, error)` — paginated. **NO Update or Delete methods** |

### Checkpoint
```bash
go build ./internal/repository/...
```

---

## Phase 4: Service Layer (4 files)

All in `internal/service/`. Business logic lives here. Handlers are thin.

### Task List

| # | Task | File | Key Methods |
|---|------|------|-------------|
| 4.1 | Project service | `project_service.go` | `CreateProject(ctx, name, ownerID) (*Project, vaultToken, error)`, `GetProject`, `ListProjects`, `DeleteProject` |
| 4.2 | Secret service | `secret_service.go` | `SetSecret`, `GetSecret` (ONLY method returning values), `ListKeys`, `DeleteSecret`, `BulkSetSecrets` |
| 4.3 | Member service | `member_service.go` | `AddMember(ctx, slug, email, role, inviterID) (*TeamMember, vaultToken, error)`, `RemoveMember`, `RotateCredentials` |
| 4.4 | Audit service | `audit_service.go` | `ListAuditLogs(ctx, slug, action, limit, offset)` |

### CreateProject Flow (most complex method)
```
1. Generate slug from name (lowercase, hyphens, truncate 50)
2. vaultMountPath = "envault-{slug}"
3. Create Project record in DB
4. Create 3 default environments (dev, staging, prod)
5. vaultSvc.EnableKV2(ctx, vaultMountPath)
6. Build admin policy HCL -> WritePolicy to Vault
7. CreateUserToken(policies, 8h, 24h) -> get (token, accessor)
8. Create TeamMember record (role=admin, store accessor)
9. Audit log: project.create
10. Return (project, token, nil) -- token given to user ONCE
```

### Secret Write Flow
```
1. Resolve project by slug, environment by name
2. vaultSvc.WriteSecret(mount, envName, {key: value}) -- read-merge-write
3. SecretMetadataRepo.Upsert (increment vault_version)
4. AuditLogRepo.Create (key name in metadata, NEVER the value)
```

### Member Invite Flow
```
1. Resolve project
2. FindOrCreate user by email
3. Determine env access by role (admin=all, dev=dev+staging, ci=prod+staging read-only)
4. Build + write Vault policy
5. CreateUserToken -> get (token, accessor)
6. Create TeamMember (store accessor ONLY)
7. Audit log: member.invite
8. Return (member, token) -- token given once
```

### Checkpoint
```bash
go build ./internal/service/...
```

---

## Phase 5: Middleware Stack (7 files)

All in `internal/server/middleware/`. Applied in this exact order:

```
Request -> RequestID -> Logger -> Recoverer -> RateLimiter -> CORS -> JWTValidator -> RBACEnforcer -> Handler
```

### Task List

| # | Task | File | Details |
|---|------|------|---------|
| 5.1 | Request ID | `request_id.go` | Check X-Request-ID header, generate UUID if absent, store in context, set response header |
| 5.2 | Logger | `logger.go` | Wrap ResponseWriter for status capture. Log: method, path, status, duration_ms, request_id, remote_addr. Record Prometheus metrics. |
| 5.3 | Recovery | `recovery.go` | defer recover(), log panic + stack trace, return 500 JSON |
| 5.4 | Rate limiter | `rate_limiter.go` | sync.Map of `*rate.Limiter` per IP. Auth: 10/min burst 5. Write: 30/min burst 10. Read: 100/min burst 20. Return 429 + Retry-After. Background cleanup every 5min. |
| 5.5 | CORS | `cors.go` | Set headers from config. Handle OPTIONS preflight with 204. |
| 5.6 | JWT validator | `jwt_validator.go` | Fetch JWKS from Supabase, cache RSA keys, refresh every 1h. Validate: signature + expiry + issuer + audience. Extract sub + email. Store User in context. Skip: /healthz, /readyz, /metrics. Return 401 on failure. |
| 5.7 | RBAC enforcer | `rbac.go` | Extract {slug} from URL. Load project. FindOrCreate user. Check team membership + role. Store project+member in context. Permission matrix below. Return 403 on failure. |

### RBAC Permission Matrix

| Permission | admin | developer | ci |
|-----------|-------|-----------|-----|
| secrets:read | All envs | Non-prod only | Prod + staging |
| secrets:write | All envs | Non-prod only | No |
| secrets:delete | All envs | No | No |
| members:invite | Yes | No | No |
| members:remove | Yes | No | No |
| project:delete | Yes | No | No |

### Checkpoint
```bash
go build ./internal/server/middleware/...
```

---

## Phase 6: HTTP Handlers (6 files)

All in `internal/server/handlers/`.

### Task List

| # | Task | File | Details |
|---|------|------|---------|
| 6.1 | Helpers | `helpers.go` | `RespondJSON(w, status, data)`, `RespondError(w, status, msg)`, `DecodeAndValidate(r, dst) error`. Init validator with custom `alphanum_underscore` rule for key names. |
| 6.2 | Health | `health.go` | `Healthz` (always 200), `Readyz` (pings DB, 200/503), `Metrics` (promhttp.Handler) |
| 6.3 | Projects | `project_handler.go` | `CreateProjectRequest{Name validate:"required,min=2,max=100"}`. Create returns 201 + vault_token. List, Get, Delete. |
| 6.4 | Secrets | `secret_handler.go` | `SetSecretRequest{Environment, Key, Value}`, `BulkSetSecretsRequest`. **List returns metadata only. Get is ONLY endpoint returning a value.** |
| 6.5 | Members | `member_handler.go` | `AddMemberRequest{Email validate:"email", Role validate:"oneof=admin developer ci"}`. Add returns 201 + vault_token. Rotate returns new token. |
| 6.6 | Audit | `audit_handler.go` | List with ?action=&limit=&offset=. Response: {data, total, limit, offset}. |

### API Routes Summary

```
GET    /healthz                                  # 200 always
GET    /readyz                                   # DB connectivity check
GET    /metrics                                  # Prometheus

POST   /api/v1/projects                          # create project
GET    /api/v1/projects                          # list user's projects
GET    /api/v1/projects/{slug}                   # get project
DELETE /api/v1/projects/{slug}                   # admin only

GET    /api/v1/projects/{slug}/secrets           # list key names (no values)
POST   /api/v1/projects/{slug}/secrets           # set secret
GET    /api/v1/projects/{slug}/secrets/{key}     # get value (audit logged)
DELETE /api/v1/projects/{slug}/secrets/{key}     # delete secret
POST   /api/v1/projects/{slug}/secrets/bulk      # bulk push

GET    /api/v1/projects/{slug}/members           # list members
POST   /api/v1/projects/{slug}/members           # invite member
DELETE /api/v1/projects/{slug}/members/{id}      # remove member

POST   /api/v1/projects/{slug}/rotate            # rotate credentials

GET    /api/v1/projects/{slug}/audit             # paginated audit log
```

**Note**: chi uses `{slug}` syntax, not `:slug`. Extract via `chi.URLParam(r, "slug")`.

### Checkpoint
```bash
go build ./internal/server/handlers/...
```

---

## Phase 7: Router, Server & Entry Points (4 files)

### Task List

| # | Task | File | Details |
|---|------|------|---------|
| 7.1 | Router | `internal/server/router.go` | `NewRouter(deps RouterDeps) *chi.Mux`. Applies middleware in order. Defines all route groups. Health routes skip auth. /api/v1 group uses JWT. /projects/{slug} sub-group uses RBAC. |
| 7.2 | Server | `internal/server/server.go` | `Start(router, addr) error`. HTTP server with timeouts. Graceful shutdown on SIGINT/SIGTERM. |
| 7.3 | API entry | `cmd/server/main.go` | Startup: zerolog config -> config.Load -> db.Connect -> goose.Up (embedded migrations) -> vault.NewVaultService -> StartTokenRenewal -> init all repos -> NewRouter -> Start |
| 7.4 | CLI entry | `cmd/envault/main.go` | Simply calls `cli.Execute()` |

### cmd/server/main.go Startup Sequence
```go
1. Configure zerolog
2. cfg := config.Load()
3. database := db.Connect(cfg.Database)
4. sqlDB, _ := database.DB()
   goose.SetDialect("postgres")
   goose.Up(sqlDB, "migrations")  // embedded via //go:embed
5. vaultSvc := vault.NewVaultService(cfg.Vault)
6. ctx, cancel := context.WithCancel(context.Background())
   vaultSvc.StartTokenRenewal(ctx)
7. Initialize all 6 repositories
8. router := server.NewRouter(deps)
9. server.Start(router, addr)
```

### Checkpoint
```bash
go mod tidy
go build -o bin/server ./cmd/server
go build -o bin/envault ./cmd/envault
go run ./cmd/server
curl http://localhost:8080/healthz   # {"status":"ok"}
curl http://localhost:8080/readyz    # {"status":"ok"}
```

---

## Phase 8: Docker Integration

### Task List

| # | Task | Details |
|---|------|---------|
| 8.1 | Uncomment api service | In docker-compose.yml, enable the envault-api block |
| 8.2 | Build Docker image | `docker compose build api` |
| 8.3 | Run full stack | `docker compose up -d` |
| 8.4 | Verify health | `curl localhost:8080/healthz` + `curl localhost:8080/readyz` |

---

## Phase 9: CLI Commands (10 files)

All in `internal/cli/`.

### Task List

| # | Task | File | Details |
|---|------|------|---------|
| 9.1 | Root command | `root.go` | Viper setup: reads ~/.envault.yaml, env prefix ENVAULT_. Binds: api_url, vault_addr, vault_token, project_slug. Registers all subcommands. |
| 9.2 | Init | `init_cmd.go` | `envault init <project>`. POST /api/v1/projects. Writes ~/.envault.yaml with slug + token. |
| 9.3 | Env pull | `env_pull.go` | `envault env pull --env`. GET secrets list + GET each value. Writes .env file. |
| 9.4 | Env push | `env_push.go` | `envault env push --env --file .env`. Parse .env, POST /secrets/bulk. |
| 9.5 | Env list | `env_list.go` | `envault env list --env`. GET /secrets. Print table (key, version, modified). |
| 9.6 | Secret set | `secret_set.go` | `envault secret set KEY=VALUE --env`. Empty value -> promptui masked input. POST /secrets. |
| 9.7 | Secret get | `secret_get.go` | `envault secret get KEY --env`. GET /secrets/{key}. Print raw value. |
| 9.8 | Secret delete | `secret_delete.go` | `envault secret delete KEY --env`. Confirm with promptui. DELETE /secrets/{key}. |
| 9.9 | Onboard | `onboard.go` | `envault onboard <email> --role`. POST /members. Print VAULT_ADDR + VAULT_TOKEN once. |
| 9.10 | Rotate | `rotate.go` | `envault rotate`. POST /rotate. Update ~/.envault.yaml. |

### Checkpoint
```bash
go build -o bin/envault ./cmd/envault
./bin/envault --help
./bin/envault env --help
./bin/envault secret --help
```

---

## Phase 10: Prometheus Metrics

### Task List

| # | Task | File | Details |
|---|------|------|---------|
| 10.1 | HTTP metrics | `middleware/logger.go` | Add `envault_http_requests_total` CounterVec (method, path, status) + `envault_http_request_duration_seconds` HistogramVec |
| 10.2 | Vault metrics | `vault/vault.go` | Add `envault_vault_operations_total` CounterVec (operation, status) |

### Checkpoint
```bash
curl http://localhost:8080/metrics | grep envault_
```

---

## Phase 11: Next.js Dashboard

### 11.1 — Scaffold

| # | Task | Details |
|---|------|---------|
| 11.1.1 | Create Next.js app | `cd web && npx create-next-app@14 . --typescript --tailwind --eslint --app --src-dir --import-alias "@/*"` |
| 11.1.2 | Install deps | `npm install @supabase/supabase-js @supabase/ssr` |
| 11.1.3 | Install UI deps | `npm install clsx tailwind-merge lucide-react` |
| 11.1.4 | Create .env.local.example | `NEXT_PUBLIC_SUPABASE_URL`, `NEXT_PUBLIC_SUPABASE_ANON_KEY`, `NEXT_PUBLIC_API_URL=http://localhost:8080` |

### 11.2 — Supabase Auth Integration

| # | Task | File | Details |
|---|------|------|---------|
| 11.2.1 | Browser client | `src/lib/supabase/client.ts` | `createBrowserClient()` using @supabase/ssr |
| 11.2.2 | Server client | `src/lib/supabase/server.ts` | `createServerClient()` for server components |
| 11.2.3 | Auth middleware | `src/lib/supabase/middleware.ts` | Refreshes session, redirects unauthenticated users to /login |
| 11.2.4 | Next.js middleware | `src/middleware.ts` | Calls supabase middleware, protects /dashboard/* routes |

### 11.3 — API Client

| # | Task | File | Details |
|---|------|------|---------|
| 11.3.1 | API wrapper | `src/lib/api.ts` | `ApiClient` class. Methods: `createProject`, `listProjects`, `getProject`, `deleteProject`, `listSecrets`, `getSecret`, `setSecret`, `deleteSecret`, `bulkSetSecrets`, `listMembers`, `addMember`, `removeMember`, `rotateCredentials`, `listAuditLogs`. Adds `Authorization: Bearer <supabase-access-token>` header to all requests. |
| 11.3.2 | Types | `src/lib/types.ts` | TypeScript interfaces matching Go models: `User`, `Project`, `Environment`, `TeamMember`, `SecretMetadata`, `AuditLog`, `PaginatedResponse` |
| 11.3.3 | Utilities | `src/lib/utils.ts` | `cn()` (clsx+twMerge), `formatDate`, `formatRelativeTime`, `roleColor` |

### 11.4 — React Hooks (SWR-style with useEffect+useState or @tanstack/react-query)

| # | Task | File | Details |
|---|------|------|---------|
| 11.4.1 | Projects hook | `src/hooks/use-projects.ts` | `useProjects()` -> list. `useProject(slug)` -> single. `useCreateProject()`, `useDeleteProject()` |
| 11.4.2 | Secrets hook | `src/hooks/use-secrets.ts` | `useSecrets(slug, env)` -> list metadata. `useSetSecret()`, `useDeleteSecret()`, `useBulkSet()` |
| 11.4.3 | Members hook | `src/hooks/use-members.ts` | `useMembers(slug)`. `useAddMember()`, `useRemoveMember()` |
| 11.4.4 | Audit hook | `src/hooks/use-audit.ts` | `useAuditLogs(slug, action, limit, offset)` |

### 11.5 — UI Components

| # | Task | File | Details |
|---|------|------|---------|
| 11.5.1 | Button | `src/components/ui/button.tsx` | Variants: default, destructive, outline, ghost. Sizes: sm, md, lg. |
| 11.5.2 | Input | `src/components/ui/input.tsx` | With label, error state, password toggle |
| 11.5.3 | Card | `src/components/ui/card.tsx` | Header, content, footer sections |
| 11.5.4 | Badge | `src/components/ui/badge.tsx` | Role badges (admin=blue, developer=green, ci=orange) |
| 11.5.5 | Dialog | `src/components/ui/dialog.tsx` | Modal with overlay, title, description, actions |
| 11.5.6 | Table | `src/components/ui/table.tsx` | Sortable columns, pagination |
| 11.5.7 | Toast | `src/components/ui/toast.tsx` | Success/error/info notifications |
| 11.5.8 | Tabs | `src/components/ui/tabs.tsx` | Environment tab switching (dev/staging/prod) |
| 11.5.9 | Skeleton | `src/components/ui/skeleton.tsx` | Loading state placeholders |
| 11.5.10 | Dropdown | `src/components/ui/dropdown-menu.tsx` | Action menus |

### 11.6 — Layout Components

| # | Task | File | Details |
|---|------|------|---------|
| 11.6.1 | Sidebar | `src/components/layout/sidebar.tsx` | Project list, nav links (Secrets, Members, Audit, Settings), user profile, logout |
| 11.6.2 | Header | `src/components/layout/header.tsx` | Breadcrumbs, project name, environment selector |
| 11.6.3 | Nav links | `src/components/layout/nav-links.tsx` | Active state highlighting, icons |

### 11.7 — Feature Components

| # | Task | File | Details |
|---|------|------|---------|
| 11.7.1 | Project card | `src/components/projects/project-card.tsx` | Name, slug, env count, member count, last updated |
| 11.7.2 | Create project dialog | `src/components/projects/create-project-dialog.tsx` | Name input, submit, shows vault_token once |
| 11.7.3 | Secrets table | `src/components/secrets/secrets-table.tsx` | Key name, version, last modified, "Reveal" button, copy button, delete action |
| 11.7.4 | Set secret dialog | `src/components/secrets/set-secret-dialog.tsx` | Key input, value textarea (masked), environment select |
| 11.7.5 | Bulk import dialog | `src/components/secrets/bulk-import-dialog.tsx` | Paste .env content or upload .env file, preview keys, submit |
| 11.7.6 | Secret value cell | `src/components/secrets/secret-value-cell.tsx` | Masked by default (dots), "Reveal" fetches value via API, auto-hide after 10s, copy to clipboard |
| 11.7.7 | Members table | `src/components/members/members-table.tsx` | Email, role badge, invited/joined dates, remove action |
| 11.7.8 | Invite member dialog | `src/components/members/invite-member-dialog.tsx` | Email input, role select, shows vault_token once |
| 11.7.9 | Audit log table | `src/components/audit/audit-log-table.tsx` | Timestamp, user, action, resource, expandable metadata row |
| 11.7.10 | Audit filters | `src/components/audit/audit-filters.tsx` | Filter by action type, date range, user |

### 11.8 — Pages

| # | Task | File | Details |
|---|------|------|---------|
| 11.8.1 | Root layout | `src/app/layout.tsx` | HTML shell, global styles, toast provider |
| 11.8.2 | Landing page | `src/app/page.tsx` | Redirect to /login or /dashboard based on auth state |
| 11.8.3 | Login page | `src/app/login/page.tsx` | Supabase Auth UI: email/password + magic link + GitHub OAuth |
| 11.8.4 | Dashboard layout | `src/app/dashboard/layout.tsx` | Sidebar + header + main content area. Protected by middleware. |
| 11.8.5 | Dashboard home | `src/app/dashboard/page.tsx` | Project grid with create button. Shows recent activity. |
| 11.8.6 | Projects list | `src/app/dashboard/projects/page.tsx` | All projects with search/filter |
| 11.8.7 | Project detail | `src/app/dashboard/projects/[slug]/page.tsx` | Overview: env count, member count, recent secrets, quick actions |
| 11.8.8 | Secrets page | `src/app/dashboard/projects/[slug]/secrets/page.tsx` | Environment tabs (dev/staging/prod). Secrets table. Set/bulk import buttons. |
| 11.8.9 | Members page | `src/app/dashboard/projects/[slug]/members/page.tsx` | Members table + invite button |
| 11.8.10 | Audit page | `src/app/dashboard/projects/[slug]/audit/page.tsx` | Audit log table with filters + pagination |
| 11.8.11 | Project settings | `src/app/dashboard/projects/[slug]/settings/page.tsx` | Rename, rotate credentials, danger zone (delete project) |
| 11.8.12 | User settings | `src/app/dashboard/settings/page.tsx` | Profile, API tokens, sign out |

### 11.9 — Docker Integration

Update `docker-compose.yml` to uncomment the envault-web service:
```yaml
envault-web:
  build: ./web
  ports: ["3000:3000"]
  env_file: ./web/.env.local
  depends_on:
    envault-api: { condition: service_healthy }
```

Create `web/Dockerfile`:
```dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./
EXPOSE 3000
CMD ["npm", "start"]
```

### Dashboard Checkpoint
```bash
cd web
npm install
npm run dev     # http://localhost:3000
npm run build   # verify production build
npm run lint    # no errors
```

---

## Phase 12: README & Documentation

### Task List

| # | Task | File | Details |
|---|------|------|---------|
| 12.1 | README | `README.md` | One-paragraph intro. ASCII architecture diagram. Prerequisites. Step-by-step setup. CLI reference table. Model overview (Supabase vs Vault). Migration workflow. Security model. Env vars reference. |

---

## Phase 13: End-to-End Verification

### 13.1 — Infrastructure Health
```bash
docker compose exec postgres pg_isready -U envault -d envault
docker compose exec vault vault status          # Sealed=false
curl http://localhost:8080/healthz               # {"status":"ok"}
curl http://localhost:8080/readyz                # {"status":"ok"}
curl http://localhost:8080/metrics | grep envault_
curl http://localhost:3000                       # Dashboard loads
```

### 13.2 — Migration Verification
```sql
-- docker compose exec postgres psql -U envault -d envault
\dt   -- 7 tables

-- CHECK constraint works
INSERT INTO environments (project_id, name) VALUES (uuid_generate_v4(), 'invalid');
-- ERROR: violates check constraint

-- Audit immutability works
UPDATE audit_logs SET action = 'x' WHERE true;
-- ERROR: permission denied
DELETE FROM audit_logs WHERE true;
-- ERROR: permission denied
```

### 13.3 — Full API Flow
```bash
export TOKEN="<jwt>" API="http://localhost:8080/api/v1"

# Unauthorized
curl -s "$API/projects" | jq .                    # 401

# Create project
curl -s -X POST "$API/projects" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"My Startup"}' | jq .               # 201 + vault_token

# Set secret
curl -s -X POST "$API/projects/my-startup/secrets" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"environment":"development","key":"DB_URL","value":"postgres://..."}' | jq .

# Get secret (ONLY endpoint with value)
curl -s "$API/projects/my-startup/secrets/DB_URL?environment=development" \
  -H "Authorization: Bearer $TOKEN" | jq .value    # "postgres://..."

# List secrets (NO values)
curl -s "$API/projects/my-startup/secrets?environment=development" \
  -H "Authorization: Bearer $TOKEN" | jq .

# Bulk push
curl -s -X POST "$API/projects/my-startup/secrets/bulk" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"environment":"development","secrets":{"API_KEY":"abc","SECRET":"xyz"}}'

# Verify in Vault directly
docker compose exec vault vault kv get -mount=envault-my-startup development

# Add member
curl -s -X POST "$API/projects/my-startup/members" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"dev@example.com","role":"developer"}' | jq .

# Audit log
curl -s "$API/projects/my-startup/audit?limit=10" \
  -H "Authorization: Bearer $TOKEN" | jq .

# Rate limiting
for i in $(seq 1 15); do
  curl -s -o /dev/null -w "%{http_code}\n" -X POST "$API/projects" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" -d '{"name":"rate-test"}'
done
# Last 5 should return 429
```

### 13.4 — CLI Verification
```bash
./bin/envault --help
./bin/envault init my-startup
echo "FOO=bar" > .env
./bin/envault env push --env development
rm .env
./bin/envault env pull --env development
cat .env                                          # FOO=bar
./bin/envault secret set MY_KEY=val --env development
./bin/envault secret get MY_KEY --env development  # val
./bin/envault env list --env development
```

### 13.5 — Security Verification
```bash
# No values in audit metadata
docker compose exec postgres psql -U envault -d envault \
  -c "SELECT metadata FROM audit_logs WHERE action='secret.write' LIMIT 5;"
# Key name only, NO value

# No values in server logs
docker compose logs api 2>&1 | grep -i "postgres://"
# No output
```

### 13.6 — Dashboard Verification
1. Open http://localhost:3000
2. Login with Supabase credentials
3. Create a project -> verify vault_token shown once
4. Navigate to secrets -> set a secret -> verify masked by default
5. Click "Reveal" -> verify value fetched and shown
6. Add a team member -> verify vault_token shown once
7. Check audit log page -> verify all actions logged

---

## Security Rules — Non-Negotiable

1. Secret values NEVER in: Supabase DB, API responses (except GET /secrets/{key}), logs, error messages
2. Store Vault token ACCESSORS in DB — never the token itself
3. All tokens have TTL: human 8h renewable to 24h, CI 1h renewable to 8h
4. Token renewal goroutine runs in API server process
5. Audit log table: `REVOKE UPDATE, DELETE` in migration SQL
6. RBAC at both API middleware AND Vault policy layer (two independent enforcement points)
7. Rate limiting: per user-ID post-auth, per IP pre-auth
8. All request bodies validated by go-playground/validator before service layer
9. JWKS keys cached on startup, background refresh every 1h
10. `.env` in `.gitignore`. No real values in committed config files.

---

## Key Design Decisions

1. **Vault runs in server mode, not dev mode** — Dev mode gives KV-v1 at `secret/` and in-memory storage. Server mode gives KV-v2 at custom mounts + persistent storage.
2. **All secrets for an environment = one Vault KV entry** — WriteSecret does read-merge-write to avoid wiping other keys.
3. **Mount creation is NOT idempotent** — EnableKV2 must check-then-create (read mount first, create only on 404).
4. **Dual RBAC enforcement** — API middleware checks HTTP-level access. Vault policies independently restrict what each token can do. Even if API has a bug, Vault blocks unauthorized access.
5. **Token accessor pattern** — CreateUserToken returns (token, accessor). Token given to user once. Accessor stored in DB for later revocation via RevokeTokenByAccessor.
6. **CLI talks through the API**, not directly to Vault — ensures all operations are audit-logged through the same pipeline.
7. **chi uses `{slug}` not `:slug`** for URL parameters — extract via `chi.URLParam(r, "slug")`.
8. **Dashboard reveals secrets on-demand** — Values are never sent in list responses. The "Reveal" button makes a separate GET /secrets/{key} call, fetching the value only when explicitly requested, and auto-hides after 10 seconds.

---

## Session Handoff Notes

### Current Status
- **Completed**: Full architecture research, dependency verification, detailed build plan
- **Not started**: No code written yet. Repo is empty.

### How to Resume in a New Session
1. Give Claude Code this file: `@BUILDPLAN.md`
2. Tell it which phase to start from (Phase 0 if fresh start)
3. Each phase has checkpoint commands — run them to verify before moving on

### Phase Dependencies
```
Phase 0 (scaffold)     -> Phase 1 (config/models)
Phase 1                -> Phase 2 (db/vault clients)
Phase 2                -> Phase 3 (repositories)
Phase 3                -> Phase 4 (services)
Phase 4                -> Phase 5 (middleware)
Phase 5                -> Phase 6 (handlers)
Phase 6                -> Phase 7 (router/server/entry)
Phase 7                -> Phase 8 (docker integration)
Phase 7                -> Phase 9 (CLI commands)   [parallel with 8]
Phase 7                -> Phase 10 (metrics)        [parallel with 8,9]
Phase 8                -> Phase 11 (dashboard)
All phases             -> Phase 12 (README)
All phases             -> Phase 13 (E2E verification)
```

### Critical Files (highest complexity, review carefully)
- `internal/vault/vault.go` — All Vault interactions, token renewal, read-merge-write, policy generation
- `internal/service/project_service.go` — Orchestrates DB + Vault mount + policy + token + audit in CreateProject
- `internal/server/middleware/jwt_validator.go` — JWKS fetch/cache, JWT validation
- `internal/server/router.go` — Wires all middleware + handlers into chi router
- `cmd/server/main.go` — Startup: migrations, dependency construction, graceful shutdown
