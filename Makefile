include .env
export

# .SILENT
# .ONESHELL
# .NOTPARALLEL
# .DEFAULT_GOAL := target-name

SHELL := /bin/bash

export APP_SERVICE_NAME := app
export APP_UID := $(shell id -u)
export APP_GID := $(shell id -g)

export ROUTE_GENERATOR_SERVICE_NAME := route-generator
export ROUTE_GENERATOR_UID := $(shell id -u)
export ROUTE_GENERATOR_GID := $(shell id -g)

# ===========================App===========================
.PHONY: app-run
app-run:
	@docker compose -f ./deploy/docker-compose.dev.yml up --build $(APP_SERVICE_NAME)

.PHONY: app-postgres-up
app-postgres-up:
	@docker compose -f ./deploy/docker-compose.dev.yml up -d app-postgres

.PHONY: app-postgres-down
app-postgres-down:
	@docker compose -f ./deploy/docker-compose.dev.yml rm -sf app-postgres

.PHONY: app-postgres-cleanup
app-postgres-cleanup: app-postgres-down
	@read -p "Would you like to cleanup env? DANGER ACTION! [y/n]: " answer; \
	if [ "$$answer" = "y" ]; then \
		docker compose -f ./deploy/docker-compose.dev.yml down app-postgres-port-forwarder && \
		docker volume rm $(APP_SERVICE_NAME)_pgdata; \
	fi;


.PHONY: app-postgres-port-forward
app-postgres-port-forward:
	@docker compose -f ./deploy/docker-compose.dev.yml up -d app-postgres-port-forwarder

.PHONY: app-postgres-port-close
app-postgres-port-close:
	@docker compose -f ./deploy/docker-compose.dev.yml down app-postgres-port-forwarder


# ===========================Route-generator===========================
.PHONY: route-generator-run
route-generator-run:
	@docker compose -f ./deploy/docker-compose.dev.yml up --build $(ROUTE_GENERATOR_SERVICE_NAME)

.PHONY: route-generator-postgres-up
route-generator-postgres-up:
	@docker compose -f ./deploy/docker-compose.dev.yml up -d route-generator-postgres

.PHONY: route-generator-postgres-down
route-generator-postgres-down:
	@docker compose -f ./deploy/docker-compose.dev.yml rm -sf route-generator-postgres

.PHONY: route-generator-postgres-cleanup
route-generator-postgres-cleanup: route-generator-postgres-down
	@read -p "Would you like to cleanup env? DANGER ACTION! [y/n]: " answer; \
	if [ "$$answer" = "y" ]; then \
		docker compose -f ./deploy/docker-compose.dev.yml down postgres-port-forwarder && \
		docker volume rm $(ROUTE_GENERATOR_SERVICE_NAME)_pgdata; \
	fi;


.PHONY: route-generator-postgres-port-forward
route-generator-postgres-port-forward:
	@docker compose -f ./deploy/docker-compose.dev.yml up -d route-generator-postgres-port-forwarder

.PHONY: route-generator-postgres-port-close
route-generator-postgres-port-close:
	@docker compose -f ./deploy/docker-compose.dev.yml down route-generator-postgres-port-forwarder
