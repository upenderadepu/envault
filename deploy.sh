#!/bin/bash
set -e

echo ""
echo "  ╔═══════════════════════════════════════════╗"
echo "  ║     Envault — Production Deployment       ║"
echo "  ╚═══════════════════════════════════════════╝"
echo ""

# Check Docker
command -v docker >/dev/null 2>&1 || { echo "Error: docker is required."; exit 1; }

# Check .env
if [ ! -f .env ]; then
  echo "Creating .env from production template..."
  cp .env.production.example .env
  echo ""
  echo "  ┌──────────────────────────────────────────────────────────────┐"
  echo "  │  ACTION REQUIRED: Edit .env with your production values     │"
  echo "  │                                                             │"
  echo "  │  1. Set DOMAIN to your domain (e.g. envault.example.com)    │"
  echo "  │  2. Set a strong DATABASE_PASSWORD                          │"
  echo "  │  3. Set a random VAULT_TOKEN (e.g. openssl rand -hex 32)   │"
  echo "  │  4. Set your Supabase credentials                          │"
  echo "  │                                                             │"
  echo "  │  Then run this script again.                                │"
  echo "  └──────────────────────────────────────────────────────────────┘"
  echo ""
  exit 0
fi

# Validate
source .env

if [ -z "$DOMAIN" ] || [[ "$DOMAIN" == *"yourdomain"* ]]; then
  echo "Error: Set DOMAIN in .env to your actual domain"
  exit 1
fi

if [[ "$DATABASE_PASSWORD" == *"CHANGE_ME"* ]]; then
  echo "Error: Set a strong DATABASE_PASSWORD in .env"
  exit 1
fi

if [[ "$VAULT_TOKEN" == *"CHANGE_ME"* ]]; then
  echo "Error: Set a random VAULT_TOKEN in .env"
  echo "  Generate one with: openssl rand -hex 32"
  exit 1
fi

if [[ "$JWKS_URL" == *"YOUR_PROJECT"* ]] || [ -z "$JWKS_URL" ]; then
  echo "Error: Set your Supabase JWKS_URL in .env"
  exit 1
fi

if [[ "$NEXT_PUBLIC_SUPABASE_URL" == *"YOUR_PROJECT"* ]] || [ -z "$NEXT_PUBLIC_SUPABASE_URL" ]; then
  echo "Error: Set NEXT_PUBLIC_SUPABASE_URL in .env"
  exit 1
fi

echo "Deploying Envault to $DOMAIN..."
echo ""

# Build and start
docker compose -f docker-compose.prod.yml up -d --build

echo ""
echo "Waiting for services..."
sleep 15

# Verify
if curl -sf http://localhost:8080/healthz > /dev/null 2>&1; then
  echo ""
  echo "  ┌─────────────────────────────────────────────────┐"
  echo "  │  Envault is live!                               │"
  echo "  │                                                 │"
  echo "  │  Dashboard:  https://$DOMAIN"
  echo "  │  API:        https://api.$DOMAIN"
  echo "  │                                                 │"
  echo "  │  Caddy will auto-provision SSL certificates.    │"
  echo "  │                                                 │"
  echo "  │  Logs:       docker compose -f docker-compose.prod.yml logs -f"
  echo "  │  Stop:       docker compose -f docker-compose.prod.yml down"
  echo "  └─────────────────────────────────────────────────┘"
  echo ""
else
  echo "Services are still starting..."
  echo "Check: docker compose -f docker-compose.prod.yml logs -f"
fi
