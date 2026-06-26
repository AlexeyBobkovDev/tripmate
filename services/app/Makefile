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
.PHONY: $(APP_SERVICE_NAME)-run
$(APP_SERVICE_NAME)-run:
	@export LOGGER_FOLDER=$(PROJECT_ROOT)/out/logs && \
	docker compose run --rm --build $(APP_SERVICE_NAME)
# 	go run ${PROJECT_ROOT}/cmd/app/main.go




# ===========================Migrate===========================
.PHONY: $(APP_SERVICE_NAME)-migrate-up
$(APP_SERVICE_NAME)-migrate-up:
	@$(MAKE) migrate-action action=up

.PHONY: $(APP_SERVICE_NAME)-migrate-down
$(APP_SERVICE_NAME)-migrate-down:
	@$(MAKE) migrate-action action=down

.PHONY: $(APP_SERVICE_NAME)-migrate-create
$(APP_SERVICE_NAME)-migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Name of the migration must be specified! E.g. name=create_user"; \
		exit 1; \
	fi; \
	$(MAKE) $(APP_SERVICE_NAME)-migrate-action action=create name=$(name)

.PHONY: $(APP_SERVICE_NAME)-migrate-action
MIGRATE_CMD = docker compose run --rm $(APP_SERVICE_NAME)-postgres-migrate
$(APP_SERVICE_NAME)-migrate-action:
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
		$(MIGRATE_CMD) $(action) \
			-path=$(MIGRATE_DIR) \
			-database=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)/$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable; \
	fi;



# ===========================Logs===========================
.PHONY: $(APP_SERVICE_NAME)-cleanup-logs
$(APP_SERVICE_NAME)-cleanup-logs:
	@read -p "Would you like to cleanup all logs? DANGER OPERATION! [y/n]: " answer; \
	if [ "$$answer" = "y" ]; then \
		rm -rf $(LOGS_DIR); \
	fi;



# ===========================Env===========================
.PHONY: $(APP_SERVICE_NAME)-env-up
$(APP_SERVICE_NAME)-env-up:
	@docker compose run -d $(APP_SERVICE_NAME)-postgres

.PHONY: $(APP_SERVICE_NAME)-env-down
$(APP_SERVICE_NAME)-env-down:
	@docker compose stop $(APP_SERVICE_NAME)-postgres

.PHONY: $(APP_SERVICE_NAME)-env-cleanup
$(APP_SERVICE_NAME)-env-cleanup:
	@read -p "Would you like to cleanup env? DANGER ACTION! [y/n]: " answer; \
	if [ "$$answer" = "y" ]; then \
		docker compose down $(APP_SERVICE_NAME)-postgres $(APP_SERVICE_NAME)-port-forwarder && \
		docker volume rm $(APP_SERVICE_NAME)_pgdata; \
	fi;



# ===========================Port forwarder===========================
.PHONY: env-port-forward
env-port-forward:
	@docker compose run port-forwarder

.PHONY: env-port-close
env-port-close:
	@docker compose stop port-forwarder



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
