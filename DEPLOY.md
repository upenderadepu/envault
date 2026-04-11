# Deploying Envault for Free

Deploy Envault using **Vercel** (frontend), **Render** (API + Vault), and **Supabase** (database + auth).

```
Browser / CLI
     |
     +----> Vercel (Next.js dashboard)
     |           |
     +---------->+----> Render (Go API) ----> Render (Vault)
                             |
                        Supabase (PostgreSQL + Auth)
```

**Total cost: $0** (free tiers). Render services sleep after 15 min of inactivity on free tier.

---

## Prerequisites

- A GitHub account (to connect repos to Vercel/Render)
- [Supabase](https://supabase.com) account
- [Render](https://render.com) account
- [Vercel](https://vercel.com) account

---

## Step 1: Supabase Setup

1. Go to [supabase.com](https://supabase.com) and create a new project
2. Save your **database password** (you'll need it later)
3. Once created, go to **Settings > API** and note:

| Value | Where to find it | Example |
|-------|-------------------|---------|
| Project URL | Settings > API | `https://abcdefgh.supabase.co` |
| anon key | Settings > API | `eyJhbGciOiJI...` |
| Project Ref | The `abcdefgh` in your URL | `abcdefgh` |

4. Go to **Settings > Database** and note:

| Value | Example |
|-------|---------|
| Host | `db.abcdefgh.supabase.co` |
| Port | `5432` |
| User | `postgres` |
| Password | (your database password) |
| Database | `postgres` |

> **Important:** Use the **connection pooler** host if you see one (`aws-0-...pooler.supabase.com` on port `6543` with user `postgres.abcdefgh`). It works better with Render.

---

## Step 2: Deploy Vault on Render

1. Go to [render.com](https://render.com) > **New > Web Service**
2. Choose **Deploy an image from a registry**
3. Image: `hashicorp/vault:1.15`
4. Configure:
   - **Name:** `envault-vault`
   - **Plan:** Free
5. Add environment variables:

```
VAULT_DEV_ROOT_TOKEN_ID = <run: openssl rand -hex 32>
VAULT_DEV_LISTEN_ADDRESS = 0.0.0.0:10000
VAULT_ADDR = http://0.0.0.0:10000
```

> Generate the token locally: `openssl rand -hex 32` and paste the result.

6. Click **Deploy** and wait until it's live
7. Note the URL (e.g. `https://envault-vault.onrender.com`)

### Initialize Vault

After deployment, open the Render **Shell** tab for your Vault service and run:

```bash
vault secrets enable -path=envault kv-v2
```

If the Shell tab is not available on the free tier, you can run this via curl from your local machine:

```bash
curl -X POST https://envault-vault.onrender.com/v1/sys/mounts/envault \
  -H "X-Vault-Token: YOUR_VAULT_TOKEN" \
  -d '{"type":"kv","options":{"version":"2"}}'
```

---

## Step 3: Deploy the API on Render

1. Go to Render > **New > Web Service**
2. Connect your GitHub repo
3. Configure:
   - **Name:** `envault-api`
   - **Root Directory:** (leave empty)
   - **Plan:** Free
4. Add environment variables:

```
SERVER_PORT               = 10000
SERVER_HOST               = 0.0.0.0
LOG_LEVEL                 = info

DATABASE_HOST             = db.abcdefgh.supabase.co
DATABASE_PORT             = 5432
DATABASE_USER             = postgres
DATABASE_PASSWORD         = your-database-password
DATABASE_NAME             = postgres
DATABASE_SSLMODE          = require

VAULT_ADDR                = https://envault-vault.onrender.com
VAULT_TOKEN               = your-vault-token-from-step-2
VAULT_MOUNT_PREFIX        = envault

JWKS_URL                  = https://abcdefgh.supabase.co/auth/v1/.well-known/jwks.json
JWT_ISSUER                = https://abcdefgh.supabase.co/auth/v1
JWT_AUDIENCE              =

CORS_ALLOWED_ORIGINS      = https://your-app.vercel.app

RATE_LIMIT_AUTH            = 10
RATE_LIMIT_WRITE           = 30
RATE_LIMIT_READ            = 100
```

> Replace all placeholder values with your actual credentials.

5. Click **Deploy**
6. Once live, test: `curl https://envault-api.onrender.com/healthz`

---

## Step 4: Deploy the Dashboard on Vercel

1. Go to [vercel.com](https://vercel.com) > **New Project**
2. Import your GitHub repo
3. Configure:
   - **Framework:** Next.js
   - **Root Directory:** `web`
4. Add environment variables:

```
NEXT_PUBLIC_SUPABASE_URL      = https://abcdefgh.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY = your-anon-key
NEXT_PUBLIC_API_URL           = https://envault-api.onrender.com
```

5. Click **Deploy**
6. Note your Vercel URL (e.g. `https://envault-app.vercel.app`)

---

## Step 5: Connect Everything

### Update Render API

Go to your `envault-api` on Render > **Environment** and update:

```
CORS_ALLOWED_ORIGINS = https://envault-app.vercel.app
```

### Update Supabase

Go to **Authentication > URL Configuration**:

- **Site URL:** `https://envault-app.vercel.app`
- **Redirect URLs:** Add `https://envault-app.vercel.app/auth/callback`

---

## Step 6: Verify

1. Open your Vercel URL
2. Sign up with email
3. Create a project
4. Add a secret
5. Check the audit log

### Test the CLI

```bash
export ENVAULT_API_URL=https://envault-api.onrender.com
export ENVAULT_SUPABASE_URL=https://abcdefgh.supabase.co
export ENVAULT_SUPABASE_ANON_KEY=your-anon-key

envault login
envault projects
```

---

## Troubleshooting

| Problem | Fix |
|---------|-----|
| API returns 401 | Check `JWKS_URL` and `JWT_ISSUER` match your Supabase project |
| CORS errors | Update `CORS_ALLOWED_ORIGINS` to your exact Vercel URL (no trailing slash) |
| Vault connection refused | Check Vault service is running on Render, verify `VAULT_ADDR` |
| Database error | Verify Supabase credentials, ensure `DATABASE_SSLMODE=require` |
| Slow first load | Free Render services sleep after 15 min, first request takes ~30s |
| Prepared statement error | Use connection pooler URL with port `6543` instead of direct connection |

---

## Cost

| Service | Free Tier | Paid (when ready) |
|---------|-----------|-------------------|
| Vercel | 100GB bandwidth | $20/mo |
| Render (API) | Sleeps after 15min | $7/mo |
| Render (Vault) | Sleeps after 15min | $7/mo |
| Supabase | 500MB DB, 50K users | $25/mo |

Free tier works great for development and small teams. Upgrade Render ($14/mo total) to avoid cold starts.
