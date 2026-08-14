include .env
export

# .SILENT
# .ONESHELL
# .NOTPARALLEL
# .DEFAULT_GOAL := target-name

export PROJECT_ROOT := $(CURDIR)
export MIGRATE_DIR := migrations
export LOGS_DIR := $(PROJECT_ROOT)/out/logs
export APP_SERVICE_NAME := app


# ===========================App===========================
.PHONY: app-run
app-run:
	@docker compose up --build $(APP_SERVICE_NAME)



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
MIGRATE_CMD = docker compose run --rm postgres-migrate
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
			-database=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable \
			$(action); \
	elif [ "$(action)" = "force" ]; then \
		if [ -z "$(v)" ]; then \
			echo "Version must be specified. E.g. v=3"; \
			exit 1; \
		fi; \
		$(MIGRATE_CMD) \
			-path=$(MIGRATE_DIR) \
			-database=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable \
			$(action) $(v); \
	fi;



# ===========================Logs===========================
.PHONY: cleanup-logs
cleanup-logs:
	@read -p "Would you like to cleanup all logs? DANGER OPERATION! [y/n]: " answer; \
	if [ "$$answer" = "y" ]; then \
		rm -rf $(LOGS_DIR); \
	fi;



# ===========================Env===========================
.PHONY: postgres-up
postgres-up:
	@docker compose up -d postgres

.PHONY: postgres-down
postgres-down:
	@docker compose down postgres

.PHONY: postgres-cleanup
postgres-cleanup:
	@read -p "Would you like to cleanup env? DANGER ACTION! [y/n]: " answer; \
	if [ "$$answer" = "y" ]; then \
		docker compose down postgres port-forwarder && \
		docker volume rm $(APP_SERVICE_NAME)_pgdata; \
	fi;



# ===========================Port forwarder===========================
.PHONY: postgres-port-forward
postgres-port-forward:
	@docker compose run -d postgres-port-forwarder

.PHONY: postgres-port-close
postgres-port-close:
	@docker compose stop postgres-port-forwarder



# ===========================Documentation===========================
.PHONY: swagger-gen
SWAGGER_CMD = docker compose run --build --rm swagger \
		init \
		-g ./cmd/$(APP_SERVICE_NAME)/main.go \
		-o docs \
		--parseInternal

swagger-gen: swagger-fmt
	@$(SWAGGER_CMD) || $(SWAGGER_CMD) --parseDependency

.PHONY: swagger-fmt
swagger-fmt:
	@docker compose run --rm --build swagger fmt \
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

.PHONY: test
test:
	@go test -v ./...
