GOLANGCI_LINT ?= $(shell go env GOPATH)/bin/golangci-lint
GOCACHE ?= $(CURDIR)/.gocache
GOMODCACHE ?= $(CURDIR)/.gomodcache
GO_TEST ?= go test

OAPI_CODEGEN ?= oapi-codegen
OPENAPI_SPEC ?= ../api-schema/generated/openapi.yaml
OAPI_CONFIG ?= oapi-codegen.yaml

.PHONY: lint
lint:
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) $(GOLANGCI_LINT) run ./...

.PHONY: test
test:
	@GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) $(GO_TEST) ./...

.PHONY: oapi
oapi:
	@test -f $(OPENAPI_SPEC) || (echo "OpenAPI spec not found at $(OPENAPI_SPEC). Run api-schema/scripts/generate-openapi.sh first."; exit 1)
	$(OAPI_CODEGEN) --config $(OAPI_CONFIG) $(OPENAPI_SPEC)
