GOLANGCI_LINT ?= $(shell go env GOPATH)/bin/golangci-lint
GOCACHE ?= $(CURDIR)/.gocache
GOMODCACHE ?= $(CURDIR)/.gomodcache
GO_TEST ?= go test

OAPI_CODEGEN ?= oapi-codegen
OPENAPI_SPEC ?= ../api-schema/generated/openapi.yaml
OAPI_CONFIG ?= oapi-codegen.yaml
MIGRATE ?= migrate
MIGRATIONS_DIR ?= migrations
NAME ?=
DB_HOST ?= db
ENV_FILE ?= .env

ifneq (,$(wildcard .env))
include .env
export $(shell grep -E '^[A-Za-z0-9_]+=' .env | cut -d= -f1)
endif

# Build DATABASE_URL after .env is loaded, allowing override from environment.
DATABASE_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: lint
lint:
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) $(GOLANGCI_LINT) run ./...

.PHONY: test
test:
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) $(GO_TEST) ./...

.PHONY: sqlc
sqlc:
	@sqlc generate

.PHONY: oapi
oapi:
	@test -f $(OPENAPI_SPEC) || (echo "OpenAPI spec not found at $(OPENAPI_SPEC). Run api-schema/scripts/generate-openapi.sh first."; exit 1)
	$(OAPI_CODEGEN) --config $(OAPI_CONFIG) $(OPENAPI_SPEC)

.PHONY: build
build:
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go build ./...

.PHONY: migrate-up
migrate-up:
	@bash -c 'set -a; if [ -f $(ENV_FILE) ]; then source $(ENV_FILE); fi; set +a; \
	if [ -z "$$MIGRATION_DATABASE_URL" ]; then echo "Set MIGRATION_DATABASE_URL in $(ENV_FILE) or environment"; exit 1; fi; \
	"$(MIGRATE)" -path "$(MIGRATIONS_DIR)" -database "$$MIGRATION_DATABASE_URL" up'

.PHONY: migrate-down
migrate-down:
	@bash -c 'set -a; if [ -f $(ENV_FILE) ]; then source $(ENV_FILE); fi; set +a; \
	if [ -z "$$MIGRATION_DATABASE_URL" ]; then echo "Set MIGRATION_DATABASE_URL in $(ENV_FILE) or environment"; exit 1; fi; \
	"$(MIGRATE)" -path "$(MIGRATIONS_DIR)" -database "$$MIGRATION_DATABASE_URL" down 1'

.PHONY: migrate-new
migrate-new:
	@test -n "$(NAME)" || (echo "NAME is required, e.g. make migrate-new NAME=add_users"; exit 1)
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) $(NAME)
