# TripMate — Application Service development commands
.DEFAULT_GOAL := help

SHELL := /bin/bash

# Load repository-level environment settings when available.
-include ../../.env
export

export PROJECT_ROOT := $(CURDIR)
export MIGRATE_DIR := migrations
export LOGS_DIR := $(PROJECT_ROOT)/out/logs
export APP_SERVICE_NAME := app

ROOT_DIR := $(abspath ../..)

export UID := $(shell id -u)
export GID := $(shell id -g)

COMPOSE := docker compose --env-file $(ROOT_DIR)/.env -f $(ROOT_DIR)/deploy/docker-compose.dev.yml

# Show available Application Service migration, Go, testing, and maintenance commands.
.PHONY: help
help:
	@printf '\n\033[1;36m╔══════════════════════════════════════════════════════════════════╗\033[0m\n'
	@printf '\033[1;36m║                TripMate · Application Service                              ║\033[0m\n'
	@printf '\033[1;36m╚══════════════════════════════════════════════════════════════════╝\033[0m\n'
	@printf '\n\033[1;33mUSAGE\033[0m\n'
	@printf '  \033[1;32mmake\033[0m or \033[1;32mmake help\033[0m       Show this command reference.\n'
	@printf '\n\033[1;33mDATABASE MIGRATIONS\033[0m\n'
	@printf '  \033[1;32mmigrate-up\033[0m                    Apply all available up migrations.\n'
	@printf '  \033[1;32mmigrate-down\033[0m                  Roll back the most recent migration.\n'
	@printf '  \033[1;32mmigrate-create name=NAME\033[0m      Create a sequential up/down migration pair.\n'
	@printf '  \033[1;32mmigrate-action action=ACTION\033[0m  Run an internal migration action (up/down/create/force).\n'
	@printf '\n\033[1;33mAPI DOCUMENTATION\033[0m\n'
	@printf '  \033[1;32mswagger-gen\033[0m                    Generate Swagger/OpenAPI documentation.\n'
	@printf '  \033[1;32mswagger-fmt\033[0m                    Format generated Swagger annotations.\n'
	@printf '\n\033[1;33mGO DEVELOPMENT\033[0m\n'
	@printf '  \033[1;32mtidy\033[0m                          Reconcile Go module dependencies.\n'
	@printf '  \033[1;32mdownload\033[0m                       Download dependencies declared by the module.\n'
	@printf '  \033[1;32mfmt\033[0m                           Format Go code and organize imports.\n'
	@printf '  \033[1;32mlint\033[0m                          Run golangci-lint across all packages.\n'
	@printf '  \033[1;32mupdate-deps\033[0m                   Update dependencies, then tidy the module.\n'
	@printf '  \033[1;32mgenerate\033[0m                      Run go generate directives.\n'
	@printf '\n\033[1;33mTESTING\033[0m\n'
	@printf '  \033[1;32mtest-unit\033[0m                     Run unit tests with the race detector.\n'
	@printf '  \033[1;32mtest-integration\033[0m              Run integration tests with the race detector.\n'
	@printf '  \033[1;32mtest-e2e\033[0m                      Run end-to-end tests with the race detector.\n'
	@printf '  \033[1;32mtest\033[0m                          Run the complete test suite.\n'
	@printf '  \033[1;32mbenchmark\033[0m                    Run Go benchmarks and report memory allocations.\n'
	@printf '\n\033[1;33mEXAMPLES\033[0m\n'
	@printf '  make migrate-create name=add_places_table\n'
	@printf '  make migrate-up\n  make test-unit\n  make test\n'
	@printf '\n\033[1;90mRun commands from this service directory. Some targets require Docker and local .env configuration.\033[0m\n\n'

# Apply all pending database migrations.
.PHONY: migrate-up
migrate-up:
	@$(MAKE) migrate-action action=up

# Roll back the most recently applied database migration.
.PHONY: migrate-down
migrate-down:
	@$(MAKE) migrate-action action=down

# Create a sequential pair of SQL migration files; requires name=... .
.PHONY: migrate-create
migrate-create:
	@if [ -z "$(name)" ]; then \
		printf '\033[1;31mERROR:\033[0m specify a migration name, e.g. make migrate-create name=create_places\n'; \
		exit 1; \
	fi
	@$(MAKE) migrate-action action=create name="$(name)"

# Internal migration runner shared by migrate-up, migrate-down, and migrate-create.
.PHONY: migrate-action
MIGRATE_CMD = $(COMPOSE) run --rm --volume ./migrations:/migrations postgres-migrate
migrate-action:
	@if [ -z "$(action)" ]; then \
		printf '\033[1;31mERROR:\033[0m specify action=up, action=down, action=create, or action=force\n'; \
		exit 1; \
	fi
	@if [ "$(action)" = "create" ]; then \
		if [ -z "$(name)" ]; then printf '\033[1;31mERROR:\033[0m specify name=...\n'; exit 1; fi; \
		$(MIGRATE_CMD) create -dir $(MIGRATE_DIR) -ext sql -seq "$(name)"; \
	elif [ "$(action)" = "up" ] || [ "$(action)" = "down" ]; then \
		$(MIGRATE_CMD) -path=$(MIGRATE_DIR) -database=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(APP_POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable $(action); \
	elif [ "$(action)" = "force" ]; then \
		if [ -z "$(v)" ]; then printf '\033[1;31mERROR:\033[0m specify migration version with v=VERSION\n'; exit 1; fi; \
		$(MIGRATE_CMD) -path=$(MIGRATE_DIR) -database=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(APP_POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable force "$(v)"; \
	else \
		printf '\033[1;31mERROR:\033[0m unsupported migration action: %s\n' "$(action)"; exit 1; \
	fi

# Remove all files in this service's log directory after explicit confirmation.
.PHONY: cleanup-logs
cleanup-logs:
	@read -r -p "Delete all logs under $(LOGS_DIR)? [y/N] " answer; \
	if [[ "$$answer" == "y" || "$$answer" == "Y" ]]; then \
		rm -rf "$(LOGS_DIR)"/*; \
		printf '\033[1;32mLogs cleaned.\033[0m\n'; \
	else \
		printf '\033[1;33mCancelled; logs were preserved.\033[0m\n'; \
	fi

# Generate Swagger/OpenAPI files, retrying with dependency parsing when needed.
.PHONY: swagger-gen
SWAGGER_CMD = $(COMPOSE) run --build --rm --workdir /code/services/app swagger init -g ./cmd/$(APP_SERVICE_NAME)/main.go -o docs --parseInternal
swagger-gen: swagger-fmt
	@if $(SWAGGER_CMD) || $(SWAGGER_CMD) --parseDependency; then \
		printf '\033[1;32mSwagger documentation generated successfully.\033[0m\n'; \
	else \
		status=$$?; printf '\033[1;31mSwagger generation failed.\033[0m\n'; exit $$status; \
	fi

# Format Swagger annotations and comments in the App source tree.
.PHONY: swagger-fmt
swagger-fmt:
	@$(COMPOSE) run --rm --build --volume .:/code swagger fmt --dir ./cmd/$(APP_SERVICE_NAME),./internal

# Normalize module dependencies and remove unused requirements.
.PHONY: tidy
tidy:
	@go mod tidy

# Download all dependencies required by this service's Go module.
.PHONY: download
download:
	@go mod download

# Apply goimports, gofumpt, and the repository's import grouping rules.
.PHONY: fmt
fmt:
	@find . -type f -name '*.go' -print0 | while IFS= read -r -d '' file; do \
		start_sha=$$(sha256sum "$${file}" | cut -d ' ' -f 1) && \
		goimports -w "$${file}" && \
		gofumpt -w "$${file}" && \
		gci write --custom-order -s standard -s default -s "prefix(github.com/AlexeyBobkovDev/tripmate)" "$${file}" && \
		end_sha=$$(sha256sum "$${file}" | cut -d ' ' -f 1) && \
		if [ "$${start_sha}" != "$${end_sha}" ]; then \
			echo "$${file}"; \
		fi; \
	done

# Run golangci-lint against every package in this service.
.PHONY: lint
lint:
	@golangci-lint run ./...

# Upgrade module dependencies to their latest compatible versions and tidy the module.
.PHONY: update-deps
update-deps:
	@go get -u ./... && go mod tidy

# Execute all go:generate directives in this service.
.PHONY: generate
generate:
	@go generate ./...

# Run benchmarks and report memory allocations for each benchmark.
.PHONY: benchmark
benchmark:
	@go test -bench=. -benchmem ./...

# Run unit tests only, with race detection and uncached results.
.PHONY: test-unit
test-unit:
	@if go test -race -count=1 -v ./...; then \
		printf '\n\033[1;32mALL UNIT TESTS PASSED\033[0m\n'; \
	else \
		status=$$?; printf '\n\033[1;31mUNIT TESTS FAILED\033[0m\n'; exit $$status; \
	fi

# Run integration-tagged tests only, with race detection and uncached results.
.PHONY: test-integration
test-integration:
	@if go test -tags=integration -run 'Integration' -race -count=1 -v ./...; then \
		printf '\n\033[1;32mALL INTEGRATION TESTS PASSED\033[0m\n'; \
	else \
		status=$$?; printf '\n\033[1;31mINTEGRATION TESTS FAILED\033[0m\n'; exit $$status; \
	fi

# Run end-to-end tests only, with race detection and service log output enabled.
.PHONY: test-e2e
test-e2e:
	@if LOGGER_FOLDER=../../../out/logs go test -tags=e2e -run 'E2E' -race -count=1 -v ./tests/e2e/...; then \
		printf '\n\033[1;32mALL E2E TESTS PASSED\033[0m\n'; \
	else \
		status=$$?; printf '\n\033[1;31mE2E TESTS FAILED\033[0m\n'; exit $$status; \
	fi

# Run the full suite, including integration and E2E build-tagged tests.
.PHONY: test
test:
	@if LOGGER_FOLDER=../../../out/logs go test -tags=integration,e2e -race -count=1 -v ./...; then \
		printf '\n\033[1;32mALL TESTS PASSED\033[0m\n'; \
	else \
		status=$$?; printf '\n\033[1;31mTEST SUITE FAILED\033[0m\n'; exit $$status; \
	fi
