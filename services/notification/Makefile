include .env
export

# .SILENT
# .ONESHELL
# .NOTPARALLEL
# .DEFAULT_GOAL := target-name

export PROJECT_ROOT := $(CURDIR)
export MIGRATE_DIR := $(PROJECT_ROOT)/migrations
export LOGS_DIR := $(PROJECT_ROOT)/out/logs
export NOTIFICATION_TARGET_NAME := notification
export NOTIFICATION_SERVICE_NAME := notification


# ===========================App===========================
.PHONY: $(NOTIFICATION_TARGET_NAME)-run
$(NOTIFICATION_TARGET_NAME)-run:
	@docker compose up --build $(NOTIFICATION_SERVICE_NAME)



# ===========================Migrate===========================
.PHONY: $(NOTIFICATION_TARGET_NAME)-migrate-up
$(NOTIFICATION_TARGET_NAME)-migrate-up:
	@$(MAKE) $(NOTIFICATION_SERVICE_NAME)-migrate-action action=up

.PHONY: $(NOTIFICATION_TARGET_NAME)-migrate-down
$(NOTIFICATION_TARGET_NAME)-migrate-down:
	@$(MAKE) $(NOTIFICATION_SERVICE_NAME)-migrate-action action=down

.PHONY: $(NOTIFICATION_TARGET_NAME)-migrate-create
$(NOTIFICATION_TARGET_NAME)-migrate-create:
	@$(MAKE) $(NOTIFICATION_TARGET_NAME)-migrate-action action=create

.PHONY: $(NOTIFICATION_TARGET_NAME)-migrate-action
MIGRATE_CMD = docker compose run --rm $(NOTIFICATION_SERVICE_NAME)-postgres-migrate
$(NOTIFICATION_TARGET_NAME)-migrate-action:
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
	else \
		$(MIGRATE_CMD) \
			-path=$(MIGRATE_DIR) \
			-database=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable \
			$(action); \
	fi;



# ===========================Logs===========================
.PHONY: $(NOTIFICATION_TARGET_NAME)-cleanup-logs
$(NOTIFICATION_TARGET_NAME)-cleanup-logs:
	@read -p "Would you like to cleanup all logs? DANGER OPERATION! [y/n]: " answer; \
	if [ "$$answer" = "y" ]; then \
		rm -rf $(LOGS_DIR); \
	fi;



# ===========================Env===========================
.PHONY: $(NOTIFICATION_TARGET_NAME)-env-up
$(NOTIFICATION_TARGET_NAME)-env-up:
	@docker compose up -d $(NOTIFICATION_SERVICE_NAME)-postgres

.PHONY: $(NOTIFICATION_TARGET_NAME)-env-down
$(NOTIFICATION_TARGET_NAME)-env-down:
	@docker compose down $(NOTIFICATION_SERVICE_NAME)-postgres

.PHONY: $(NOTIFICATION_TARGET_NAME)-env-cleanup
$(NOTIFICATION_TARGET_NAME)-env-cleanup:
	@read -p "Would you like to cleanup env? DANGER ACTION! [y/n]: " answer; \
	if [ "$$answer" = "y" ]; then \
		rm -rf out/pgdata; \
	fi;



# ===========================Port forwarder===========================
.PHONY: env-port-forward
env-port-forward:
	@docker compose run -d notification-postgres-port-forwarder

.PHONY: env-port-close
env-port-close:
	@docker compose stop notification-postgres-port-forwarder



# ===========================Documentation===========================
.PHONY: swagger-gen
SWAGGER_CMD = $(MAKE) swagger-action action=init \
		-g $(APP_PROJECT_ROOT)/cmd/$(APP_SERVICE_NAME)/main.go \
		-0 docs \
		--parseInternal

swagger-gen: swagger-fmt
	@$(SWAGGER_CMD) || $(SWAGGER_CMD) --parseDependency

.PHONY: swagger-fmt
swagger-fmt:
	@$(MAKE) swagger-action action=fmt

.PHONY: swagger-action
swagger-action:
	@if [ -z "$(action)" ]; then \
		echo "The parameter action must be provided. E.g. action=init"; \
		exit 1; \
	fi; \
	docker compose run --rm swagger $(action)



# ===========================Go===========================
.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: download
download:
	@go mod download
