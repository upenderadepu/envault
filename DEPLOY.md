# Deploying Envault (Free)

Deploy Envault for free using Vercel, Render, and Supabase.

## Architecture

```
Users
  │
  ├── https://envault.vercel.app ──── Vercel (Next.js frontend)
  │                                      │
  │                                      ▼
  └── CLI ────────────────────────── Render (Go API)
                                         │
                                    ┌────┴────┐
                                    │         │
                              Supabase DB   Render (Vault)
                              (PostgreSQL)  (HashiCorp Vault)
```

| Service | Platform | Free Tier |
|---------|----------|-----------|
| Frontend | Vercel | Unlimited |
| Go API | Render | 750 hrs/month |
| PostgreSQL | Supabase | 500MB, 2 projects |
| Vault | Render | 750 hrs/month |
| Auth | Supabase | 50K MAU |

---

## Step 1: Supabase Setup (5 min)

You probably already have a Supabase project. If not:

1. Go to [supabase.com](https://supabase.com) → **New Project**
2. Save your:
   - **Project URL**: `https://xxxxx.supabase.co`
   - **anon key**: from Settings → API
   - **Database password**: you set this during creation

3. Get your **database connection details** from **Settings → Database**:
   - Host: `db.xxxxx.supabase.co`
   - Port: `5432`
   - User: `postgres`
   - Password: (your database password)
   - Database: `postgres`

4. Go to **Authentication → URL Configuration**:
   - Site URL: `https://your-app.vercel.app` (update after Vercel deploy)
   - Add Redirect URL: `https://your-app.vercel.app/auth/callback`

---

## Step 2: Deploy Vault on Render (5 min)

1. Go to [render.com](https://render.com) → **New → Web Service**
2. Select **Deploy an image from a registry**
3. Image URL: `hashicorp/vault:1.15`
4. Settings:
   - **Name**: `envault-vault`
   - **Plan**: Free
   - **Environment Variables**:
     ```
     VAULT_DEV_ROOT_TOKEN_ID = your-random-token-here
     VAULT_DEV_LISTEN_ADDRESS = 0.0.0.0:8200
     VAULT_ADDR = http://0.0.0.0:8200
     ```
   - Generate a random token: run `openssl rand -hex 32` in your terminal
5. Click **Deploy**
6. Once running, note the URL (e.g. `https://envault-vault.onrender.com`)
7. Enable the KV engine — open the Render **Shell** tab and run:
   ```bash
   vault secrets enable -path=envault kv-v2
   ```

---

## Step 3: Deploy Go API on Render (5 min)

1. Go to [render.com](https://render.com) → **New → Web Service**
2. Connect your **GitHub repo** (`bhartiyaanshul/envault`)
3. Settings:
   - **Name**: `envault-api`
   - **Root Directory**: (leave empty — Dockerfile is at root)
   - **Plan**: Free
   - **Environment Variables**:

     ```
     SERVER_PORT = 8080
     SERVER_HOST = 0.0.0.0
     LOG_LEVEL = info

     # Supabase PostgreSQL (from Step 1)
     DATABASE_HOST = db.xxxxx.supabase.co
     DATABASE_PORT = 5432
     DATABASE_USER = postgres
     DATABASE_PASSWORD = your-supabase-db-password
     DATABASE_NAME = postgres
     DATABASE_SSLMODE = require

     # Vault (from Step 2)
     VAULT_ADDR = https://envault-vault.onrender.com
     VAULT_TOKEN = your-vault-token-from-step-2
     VAULT_MOUNT_PREFIX = envault

     # Supabase Auth
     JWKS_URL = https://xxxxx.supabase.co/auth/v1/.well-known/jwks.json
     JWT_ISSUER = https://xxxxx.supabase.co/auth/v1
     JWT_AUDIENCE =

     # CORS (update after Vercel deploy)
     CORS_ALLOWED_ORIGINS = https://your-app.vercel.app

     RATE_LIMIT_AUTH = 10
     RATE_LIMIT_WRITE = 30
     RATE_LIMIT_READ = 100
     ```

4. Click **Deploy**
5. Once running, note the URL (e.g. `https://envault-api.onrender.com`)
6. Test: `curl https://envault-api.onrender.com/healthz`

---

## Step 4: Deploy Frontend on Vercel (3 min)

1. Go to [vercel.com](https://vercel.com) → **New Project**
2. Import your GitHub repo (`bhartiyaanshul/envault`)
3. Settings:
   - **Framework**: Next.js
   - **Root Directory**: `web`
   - **Environment Variables**:
     ```
     NEXT_PUBLIC_SUPABASE_URL = https://xxxxx.supabase.co
     NEXT_PUBLIC_SUPABASE_ANON_KEY = your-anon-key
     NEXT_PUBLIC_API_URL = https://envault-api.onrender.com
     ```
4. Click **Deploy**
5. Note your Vercel URL (e.g. `https://envault.vercel.app`)

---

## Step 5: Connect Everything (2 min)

Now update the services to know about each other:

### Update Render (Go API)
- Go to your `envault-api` service on Render → **Environment**
- Update `CORS_ALLOWED_ORIGINS` to your Vercel URL:
  ```
  CORS_ALLOWED_ORIGINS = https://envault.vercel.app
  ```

### Update Supabase
- Go to **Authentication → URL Configuration**
- Set **Site URL**: `https://envault.vercel.app`
- Add **Redirect URL**: `https://envault.vercel.app/auth/callback`

### (Optional) Custom Domain
- On Vercel: **Settings → Domains → Add** your domain
- On Render: **Settings → Custom Domain** for the API
- Update CORS and Supabase redirect URLs to match

---

## Step 6: Verify

1. Open `https://envault.vercel.app` — you should see the landing page
2. Click **Get Started** → sign up with email
3. Create a project — you should get a Vault token
4. Add a secret — it gets stored in Vault
5. Check the audit log — your action is recorded

### Test CLI

```bash
make build-cli
./bin/envault init my-project \
  --api-url https://envault-api.onrender.com \
  --vault-token your-vault-token

./bin/envault secret set API_KEY=test123 --env development
./bin/envault env list --env development
```

---

## Troubleshooting

### API returns 401
- Check that `JWKS_URL` and `JWT_ISSUER` are correct in Render env vars
- Verify Supabase project URL matches

### CORS errors in browser
- Update `CORS_ALLOWED_ORIGINS` on Render to match your exact Vercel URL
- Make sure there's no trailing slash

### Vault connection refused
- Check that the Vault service on Render is running
- Verify `VAULT_ADDR` in the API's env vars matches the Vault service URL

### Database connection error
- Verify Supabase database credentials in Render env vars
- Make sure `DATABASE_SSLMODE=require` is set

### Render free tier spins down
- Free Render services sleep after 15 min of inactivity
- First request after sleep takes ~30 seconds (cold start)
- Upgrade to paid ($7/mo) to keep it always on

---

## Cost Summary

| Service | Free Tier Limits | Paid (if needed) |
|---------|-----------------|-------------------|
| Vercel | 100GB bandwidth | $20/mo |
| Render (API) | Sleeps after 15min | $7/mo |
| Render (Vault) | Sleeps after 15min | $7/mo |
| Supabase | 500MB DB, 50K users | $25/mo |

**For a real product with users**, you'll eventually want Render paid ($14/mo for API + Vault) to avoid cold starts. But free tier is perfect for launching and getting initial users.
