# Envault

**Self-hosted secrets management for teams.** Store secrets in HashiCorp Vault, manage access with roles, and inject environment variables anywhere — from a modern dashboard or CLI.

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev)
[![Next.js](https://img.shields.io/badge/Next.js-14-000000?style=flat&logo=next.js)](https://nextjs.org)
[![Vault](https://img.shields.io/badge/Vault-1.15-FFEC6E?style=flat&logo=vault)](https://www.vaultproject.io)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Why Envault?

Teams share secrets through Slack DMs, sticky notes in 1Password, or `.env` files committed to repos. Envault replaces all of that with:

- **Vault-backed storage** — secret values live in HashiCorp Vault, never in your database
- **Role-based access** — admin, developer, and CI roles with granular permissions
- **Full audit trail** — every read, write, and delete is logged
- **Modern dashboard** — manage everything from a clean UI with dark mode
- **CLI-first workflow** — `envault env pull` and you're ready to go

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Next.js Dashboard                       │
│                  (React, Tailwind, shadcn/ui)                │
│                     localhost:3000                           │
└──────────────────────────┬──────────────────────────────────┘
                           │ REST API
┌──────────────────────────▼──────────────────────────────────┐
│                      Go API Server                          │
│          (Chi router, RBAC middleware, Prometheus)           │
│                     localhost:8080                           │
├─────────────┬───────────────────────┬───────────────────────┤
│  Supabase   │    PostgreSQL         │   HashiCorp Vault     │
│  (JWT Auth) │  (metadata, teams,    │   (secret values,     │
│             │   audit logs)         │    KV-v2 engine)      │
└─────────────┴───────────────────────┴───────────────────────┘

CLI (envault) ──── REST API ──── Go API Server
```

**Data split:** PostgreSQL stores users, projects, environments, team members, secret metadata, and audit logs. Vault stores actual secret values. Secret values never touch PostgreSQL.

---

## Quick Start (Docker — recommended)

**Prerequisites:** Docker, a free [Supabase](https://supabase.com) project

### 1. Clone and configure

```bash
git clone https://github.com/bhartiyaanshul/envault.git
cd envault
cp .env.example .env
```

Edit `.env` and fill in your Supabase details:

```env
# From Supabase Dashboard > Settings > API
JWKS_URL=https://YOUR_PROJECT.supabase.co/auth/v1/.well-known/jwks.json
JWT_ISSUER=https://YOUR_PROJECT.supabase.co/auth/v1
NEXT_PUBLIC_SUPABASE_URL=https://YOUR_PROJECT.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=your-anon-key-here
```

### 2. Start everything

```bash
./setup.sh
```

Or manually:

```bash
docker compose up -d --build
```

That's it. Open http://localhost:3000, sign up, and create your first project.

### What gets started

| Service | Port | Description |
|---------|------|-------------|
| Dashboard | [localhost:3000](http://localhost:3000) | Next.js web UI |
| API Server | [localhost:8080](http://localhost:8080) | Go REST API |
| Vault | [localhost:8200](http://localhost:8200) | HashiCorp Vault (dev mode) |
| PostgreSQL | localhost:5432 | Metadata database |
| Metrics | [localhost:8080/metrics](http://localhost:8080/metrics) | Prometheus endpoint |

---

## Local Development (without Docker)

If you want to hack on Envault itself:

### Prerequisites

- Go 1.22+
- Node.js 18+
- Docker (for PostgreSQL and Vault only)

### Setup

```bash
# 1. Start infra
docker compose up -d postgres vault

# 2. Enable Vault KV engine
docker compose exec -e VAULT_TOKEN=dev-root-token vault vault secrets enable -path=envault kv-v2

# 3. Configure
cp .env.example .env
# Edit .env with your Supabase credentials

# 4. Start API server (auto-runs migrations)
source .env && go run ./cmd/server/

# 5. Start dashboard (new terminal)
cd web && npm install && npm run dev

# 6. Build CLI
make build-cli
./bin/envault --help
```

---

## CLI Usage

```bash
# Build
make build-cli

# Initialize a project
./bin/envault init my-app --api-url http://localhost:8080 --vault-token dev-root-token

# Manage secrets
./bin/envault secret set DATABASE_URL=postgresql://... --env production
./bin/envault secret get DATABASE_URL --env production
./bin/envault secret delete OLD_KEY --env staging

# Bulk operations
./bin/envault env push --env development -f .env
./bin/envault env pull --env production -o .env
./bin/envault env list --env production

# Team management
./bin/envault onboard alice@company.com --role developer
./bin/envault rotate
```

Config is stored in `~/.envault.yaml` after `init`:

```yaml
api_url: http://localhost:8080
project_slug: my-app
vault_token: hvs.CAESxxxxxx
```

### CLI Commands

| Command | Description |
|---------|-------------|
| `envault init <name>` | Create project, save config to `~/.envault.yaml` |
| `envault secret set KEY=VALUE --env <env>` | Set a secret |
| `envault secret get KEY --env <env>` | Get a secret value |
| `envault secret delete KEY --env <env>` | Delete a secret |
| `envault env pull --env <env> -o .env` | Download all secrets to a file |
| `envault env push --env <env> -f .env` | Upload a `.env` file |
| `envault env list --env <env>` | List all secret keys |
| `envault onboard <email> --role <role>` | Invite a team member |
| `envault rotate` | Rotate project credentials |

---

## Deploying to Production

### Option 1: VPS / Bare Metal (Docker Compose)

Works on any machine with Docker — DigitalOcean, Hetzner, AWS EC2, etc.

```bash
# On your server
git clone https://github.com/bhartiyaanshul/envault.git
cd envault
cp .env.example .env

# Edit .env:
# - Set a strong DATABASE_PASSWORD
# - Set your Supabase credentials
# - Set CORS_ALLOWED_ORIGINS to your domain
# - Set NEXT_PUBLIC_API_URL to your API's public URL

docker compose up -d --build
```

Put nginx or Caddy in front for HTTPS:

```
# Caddyfile example
envault.yourdomain.com {
    reverse_proxy localhost:3000
}

api.envault.yourdomain.com {
    reverse_proxy localhost:8080
}
```

### Option 2: Railway / Render / Fly.io

Deploy each service separately:

1. **PostgreSQL** — Use the platform's managed PostgreSQL
2. **Vault** — Deploy the Vault container or use HCP Vault (HashiCorp's managed service)
3. **API Server** — Deploy from the `Dockerfile` in root
4. **Dashboard** — Deploy from `web/Dockerfile`

Set all environment variables from `.env.example` in each service's config.

### Production Checklist

- [ ] Use a strong, unique `DATABASE_PASSWORD`
- [ ] Use a real Vault token (not `dev-root-token`)
- [ ] Set `CORS_ALLOWED_ORIGINS` to your actual domain
- [ ] Set `NEXT_PUBLIC_API_URL` to your API's public URL
- [ ] Enable HTTPS (TLS) via reverse proxy
- [ ] For production Vault, use file/raft storage instead of dev mode
- [ ] Set up Vault auto-unseal for production (AWS KMS, GCP, etc.)

---

## Supabase Setup

Envault uses [Supabase](https://supabase.com) for authentication. Here's how to set it up:

1. Create a free project at [supabase.com](https://supabase.com)
2. Go to **Project Settings > API**
3. Copy:
   - **Project URL** → `NEXT_PUBLIC_SUPABASE_URL`
   - **anon public key** → `NEXT_PUBLIC_SUPABASE_ANON_KEY`
   - **JWT Secret** is not needed (we use JWKS)
4. Set the JWKS URL:
   ```
   JWKS_URL=https://YOUR_PROJECT_REF.supabase.co/auth/v1/.well-known/jwks.json
   JWT_ISSUER=https://YOUR_PROJECT_REF.supabase.co/auth/v1
   ```
5. (Optional) Enable GitHub OAuth in **Authentication > Providers > GitHub**

---

## RBAC Roles

| Permission | Admin | Developer | CI |
|------------|-------|-----------|-----|
| Read secrets | Yes | Yes | Yes |
| Write secrets | Yes | Yes | No |
| Delete secrets | Yes | Yes | No |
| Manage members | Yes | No | No |
| Rotate credentials | Yes | No | No |
| View audit logs | Yes | Yes | Yes |
| Delete project | Yes | No | No |

---

## API Reference

All endpoints under `/api/v1` require `Authorization: Bearer <supabase-jwt>`.

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/healthz` | Health check |
| `GET` | `/readyz` | Readiness check |
| `GET` | `/metrics` | Prometheus metrics |
| `POST` | `/api/v1/projects` | Create a project |
| `GET` | `/api/v1/projects` | List projects |
| `GET` | `/api/v1/projects/{slug}` | Get project details |
| `DELETE` | `/api/v1/projects/{slug}` | Delete a project |
| `GET` | `/api/v1/projects/{slug}/secrets?environment=` | List secret keys |
| `POST` | `/api/v1/projects/{slug}/secrets` | Set a secret |
| `POST` | `/api/v1/projects/{slug}/secrets/bulk` | Bulk set secrets |
| `GET` | `/api/v1/projects/{slug}/secrets/{key}` | Get a secret value |
| `DELETE` | `/api/v1/projects/{slug}/secrets/{key}` | Delete a secret |
| `GET` | `/api/v1/projects/{slug}/members` | List team members |
| `POST` | `/api/v1/projects/{slug}/members` | Add a member |
| `DELETE` | `/api/v1/projects/{slug}/members/{id}` | Remove a member |
| `POST` | `/api/v1/projects/{slug}/rotate` | Rotate credentials |
| `GET` | `/api/v1/projects/{slug}/audit` | List audit logs |

---

## Security Model

- **Authentication**: Supabase JWT validated via JWKS (supports RS256 and ES256)
- **Authorization**: RBAC enforced per-project at the middleware level
- **Secret isolation**: Values stored only in Vault, never in PostgreSQL
- **Reveal on demand**: Dashboard masks values, auto-hides after 10 seconds
- **Audit trail**: Every secret access is logged with user identity
- **Rate limiting**: Configurable per-endpoint (auth, read, write)

---

## Prometheus Metrics

| Metric | Type | Labels |
|--------|------|--------|
| `envault_http_requests_total` | Counter | method, path, status |
| `envault_http_request_duration_seconds` | Histogram | method, path, status |
| `envault_vault_operations_total` | Counter | operation, status |

---

## Project Structure

```
envault/
├── cmd/
│   ├── server/main.go          # API server entry point
│   └── envault/main.go         # CLI entry point
├── internal/
│   ├── cli/                    # CLI commands (cobra)
│   ├── config/                 # Environment-based config
│   ├── db/                     # Database connection
│   ├── models/                 # GORM models
│   ├── repository/             # Data access layer
│   ├── server/
│   │   ├── handlers/           # HTTP handlers
│   │   └── middleware/         # JWT, RBAC, rate limiting, CORS
│   ├── service/                # Business logic
│   └── vault/                  # Vault client wrapper
├── migrations/                 # SQL migrations (embedded)
├── web/                        # Next.js 14 dashboard
├── setup.sh                    # One-command setup script
├── docker-compose.yml          # Full stack Docker config
├── Dockerfile                  # Go API multi-stage build
└── Makefile                    # Build, test, migrate targets
```

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

## License

MIT — see [LICENSE](LICENSE) for details.
