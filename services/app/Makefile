include ../../.env
export

# .SILENT
# .ONESHELL
# .NOTPARALLEL
# .DEFAULT_GOAL := target-name

SHELL := /bin/bash
export PROJECT_ROOT := $(CURDIR)
export MIGRATE_DIR := migrations
export LOGS_DIR := $(PROJECT_ROOT)/out/logs
export APP_SERVICE_NAME := app

ROOT_DIR := $(abspath ../..)

export APP_UID := $(shell id -u)
export APP_GID := $(shell id -g)

export ROUTE_GENERATOR_UID := $(shell id -u)
export ROUTE_GENERATOR_GID := $(shell id -g)

COMPOSE := docker compose \
	--env-file $(ROOT_DIR)/.env \
	-f $(ROOT_DIR)/deploy/docker-compose.dev.yml


# ===========================Migrate===========================
.PHONY: migrate-up
migrate-up:
	@$(MAKE) migrate-action action=up

.PHONY: migrate-down
migrate-down:
	@$(MAKE) migrate-action action=down

.PHONY: migrate-create
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Name of the migration must be specified! E.g. name=create_user"; \
		exit 1; \
	fi; \
	$(MAKE) migrate-action action=create name=$(name)

.PHONY: migrate-action
MIGRATE_CMD = $(COMPOSE) run --rm postgres-migrate
migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Action Must be specified! E.g. action=up"; \
		exit 1; \
	fi; \
	if [ "$(action)" = "create" ]; then \
		if [ -z "$(name)" ]; then \
			echo "Name of the migration must be specified! E.g. name=create_user"; \
			exit 1; \
		fi; \
		$(MIGRATE_CMD) $(action) \
			-dir $(MIGRATE_DIR) \
			-ext sql \
			-seq \
			$(name); \
	elif [ "$(action)" = "up" ] || [ "$(action)" = "down" ]; then \
		$(MIGRATE_CMD) \
			-path=$(MIGRATE_DIR) \
			-database=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(APP_POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable \
			$(action); \
	elif [ "$(action)" = "force" ]; then \
		if [ -z "$(v)" ]; then \
			echo "Version must be specified. E.g. v=3"; \
			exit 1; \
		fi; \
		$(MIGRATE_CMD) \
			-path=$(MIGRATE_DIR) \
			-database=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(APP_POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable \
			$(action) $(v); \
	fi;



# ===========================Logs===========================
.PHONY: cleanup-logs
cleanup-logs:
	@read -p "Would you like to cleanup all logs? DANGER OPERATION! [y/n]: " answer; \
	if [ "$$answer" = "y" ]; then \
		rm -rf $(LOGS_DIR)/*; \
	fi;


# ===========================Documentation===========================
.PHONY: swagger-gen
SWAGGER_CMD = $(COMPOSE) run --build --rm swagger \
		init \
		-g ./cmd/$(APP_SERVICE_NAME)/main.go \
		-o docs \
		--parseInternal

swagger-gen: swagger-fmt
	@$(SWAGGER_CMD) || $(SWAGGER_CMD) --parseDependency

.PHONY: swagger-fmt
swagger-fmt:
	@$(COMPOSE) run --rm --build swagger fmt \
		--dir ./cmd/$(APP_SERVICE_NAME),./internal



# ===========================Go===========================
.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: download
download:
	@go mod download

.PHONY: fmt
fmt:
	@goimports -w .; \
	gofumpt -w .; \
	gci write \
	--custom-order \
	-s standard \
	-s default \
	-s "prefix(github.com/AlexeyBobkovDev/tripmate)" \
	.;

.PHONY: lint
lint:
	@golangci-lint run ./...



# ===========================Test===========================
.PHONY: benchmark
benchmark:
	@go test -bench=. -benchmem ./...

.PHONY: test-unit
test-unit:
	@if go test -race -count=1 -v ./...; then \
		printf '\n\033[32m%s\033[0m\n' '******************** ALL UNIT TESTS PASSED ********************' ; \
	else \
		printf '\n\033[31m%s\033[0m\n' '******************** UNIT TESTS FAILED ********************'; \
	fi

.PHONY: test-integration
test-integration:
	@if go test -tags=integration -run 'Integration' -race -count=1 -v ./...; then \
		printf '\n\033[32m%s\033[0m\n' '******************** ALL INTEGRATION TESTS PASSED ********************' ; \
	else \
		printf '\n\033[31m%s\033[0m\n' '******************** INTEGRATION TESTS FAILED ********************'; \
	fi

.PHONY: test-e2e
test-e2e:
	@if LOGGER_FOLDER=../../../out/logs go test -tags=e2e -run 'E2E' -race -count=1 -v ./tests/e2e/...; then \
		printf '\n\033[32m%s\033[0m\n' '******************** ALL E2E TESTS PASSED ********************' ; \
	else \
		printf '\n\033[31m%s\033[0m\n' '******************** E2E TESTS FAILED ********************'; \
	fi

.PHONY: test
test:
	@if LOGGER_FOLDER=../../../out/logs go test -tags=integration,e2e -race -count=1 -v ./...; then \
		printf '\n\033[32m%s\033[0m\n' '******************** ALL TESTS PASSED ********************' ; \
	else \
		printf '\n\033[31m%s\033[0m\n' '******************** TESTS FAILED ********************'; \
	fi
