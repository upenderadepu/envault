.PHONY: build build-server build-cli test lint clean \
       migrate-up migrate-down migrate-status \
       docker-up docker-down vault-init vault-unseal

# Go build
build: build-server build-cli

build-server:
	go build -o bin/server ./cmd/server

build-cli:
	go build -o bin/envault ./cmd/envault

# Testing & linting
test:
	go test ./... -v -cover

lint:
	golangci-lint run ./...

# Database migrations
GOOSE_DRIVER ?= postgres
GOOSE_DBSTRING ?= "postgres://envault:envault_dev_password@localhost:5432/envault?sslmode=disable"

migrate-up:
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir migrations up

migrate-down:
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir migrations down

migrate-status:
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir migrations status

# Docker
docker-up:
	docker compose up -d

docker-down:
	docker compose down

# Vault first-time setup
vault-init:
	docker compose exec vault vault operator init -key-shares=1 -key-threshold=1 -format=json

vault-unseal:
	@echo "Usage: make vault-unseal UNSEAL_KEY=<key>"
	docker compose exec vault vault operator unseal $(UNSEAL_KEY)

# Cleanup
clean:
	rm -rf bin/
