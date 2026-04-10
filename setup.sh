#!/bin/bash
set -e

echo ""
echo "  ╔═══════════════════════════════════════╗"
echo "  ║         Envault — Quick Setup         ║"
echo "  ╚═══════════════════════════════════════╝"
echo ""

# Check prerequisites
command -v docker >/dev/null 2>&1 || { echo "Error: docker is required but not installed."; exit 1; }
command -v docker compose version >/dev/null 2>&1 || { echo "Error: docker compose is required but not installed."; exit 1; }

# Check if .env exists
if [ ! -f .env ]; then
  echo "Creating .env from template..."
  cp .env.example .env
  echo ""
  echo "  ┌─────────────────────────────────────────────────────────────┐"
  echo "  │  ACTION REQUIRED: Configure your .env file                 │"
  echo "  │                                                            │"
  echo "  │  You need a free Supabase project for authentication.      │"
  echo "  │  1. Go to https://supabase.com and create a project        │"
  echo "  │  2. Copy your project URL and anon key from Settings > API │"
  echo "  │  3. Edit .env and set:                                     │"
  echo "  │     - JWKS_URL (your Supabase JWKS endpoint)               │"
  echo "  │     - JWT_ISSUER (your Supabase JWT issuer)                │"
  echo "  │     - NEXT_PUBLIC_SUPABASE_URL                             │"
  echo "  │     - NEXT_PUBLIC_SUPABASE_ANON_KEY                        │"
  echo "  │                                                            │"
  echo "  │  Then run this script again.                               │"
  echo "  └─────────────────────────────────────────────────────────────┘"
  echo ""
  exit 0
fi

# Validate required vars
source .env
if [[ "$JWKS_URL" == *"<project>"* ]] || [ -z "$JWKS_URL" ]; then
  echo "Error: JWKS_URL is not configured in .env"
  echo "Set it to: https://<your-project>.supabase.co/auth/v1/.well-known/jwks.json"
  exit 1
fi

if [ -z "$NEXT_PUBLIC_SUPABASE_URL" ]; then
  echo "Error: NEXT_PUBLIC_SUPABASE_URL is not configured in .env"
  exit 1
fi

if [ -z "$NEXT_PUBLIC_SUPABASE_ANON_KEY" ]; then
  echo "Error: NEXT_PUBLIC_SUPABASE_ANON_KEY is not configured in .env"
  exit 1
fi

echo "Starting Envault..."
echo ""

# Start all services
docker compose up -d --build

echo ""
echo "Waiting for services to be healthy..."
sleep 10

# Check health
if curl -sf http://localhost:${SERVER_PORT:-8080}/healthz > /dev/null 2>&1; then
  echo ""
  echo "  ┌─────────────────────────────────────────────────┐"
  echo "  │  Envault is running!                            │"
  echo "  │                                                 │"
  echo "  │  Dashboard:  http://localhost:${WEB_PORT:-3000}              │"
  echo "  │  API:        http://localhost:${SERVER_PORT:-8080}              │"
  echo "  │  Metrics:    http://localhost:${SERVER_PORT:-8080}/metrics      │"
  echo "  │  Vault UI:   http://localhost:8200              │"
  echo "  │                                                 │"
  echo "  │  To stop:    docker compose down                │"
  echo "  │  Logs:       docker compose logs -f             │"
  echo "  └─────────────────────────────────────────────────┘"
  echo ""
else
  echo ""
  echo "Services are starting... check status with: docker compose ps"
  echo "View logs with: docker compose logs -f"
fi
