SHELL := /bin/bash
APP := go-bootstrap
LOCAL_INFRA_DIR := infra/local
LOCAL_ENV_FILE := $(LOCAL_INFRA_DIR)/.env
LOCAL_ENV_EXAMPLE := $(LOCAL_INFRA_DIR)/.env.example
COMPOSE_FILE := $(LOCAL_INFRA_DIR)/compose.dependencies.yaml
COMPOSE := docker compose --env-file $(LOCAL_ENV_FILE) -f $(COMPOSE_FILE)

.PHONY: help fmt fmt-check vet test race check build run local-env deps-up deps-down deps-logs \
	deps-ps docker-build zip clean

help:
	@awk 'BEGIN {FS = ":.*## "}; /^[a-zA-Z_-]+:.*## / {printf "%-14s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

fmt: ## Format all Go files
	@gofmt -w $$(find . -name '*.go')

fmt-check: ## Fail if any Go file is not formatted
	@test -z "$$(gofmt -l $$(find . -name '*.go'))"

vet: ## Run Go's built-in static analysis
	@go vet ./...

test: ## Run tests without the Go test cache
	@go test -race -count=1 ./...

check: fmt-check vet test ## Run the fast pre-submit checks

build: ## Build the API binary
	@mkdir -p bin
	@go build -trimpath -o bin/$(APP) ./cmd/server

run: ## Run the API
	@go run ./cmd/server

local-env: ## Create local dependency configuration if it does not exist
	@test -f "$(LOCAL_ENV_FILE)" || { \
		cp "$(LOCAL_ENV_EXAMPLE)" "$(LOCAL_ENV_FILE)"; \
		echo "Created $(LOCAL_ENV_FILE)"; \
	}

deps-up: local-env ## Start MySQL and DynamoDB Local (the API runs in your IDE)
	@$(COMPOSE) up --detach

deps-down: local-env ## Stop local dependency containers while preserving their data
	@$(COMPOSE) down

deps-logs: local-env ## Follow local dependency logs
	@$(COMPOSE) logs --follow

deps-ps: local-env ## Show local dependency status
	@$(COMPOSE) ps

docker-build: ## Build an optional local API image
	@docker build --tag $(APP):local .

zip: ## Create a clean submission archive under dist/
	@./scripts/package.sh

clean: ## Remove generated local artifacts
	@rm -f bin/$(APP) dist/$(APP).zip
