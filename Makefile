# TripMate — root development commands
.DEFAULT_GOAL := help

SHELL := /bin/bash

# Load local development settings when available; help remains usable without .env.
-include .env
export

export APP_SERVICE_NAME := app
export ROUTE_GENERATOR_SERVICE_NAME := route-generator
export ROUTE_GENERATOR_VOLUME_NAME := route_generator
export UID := $(shell id -u)
export GID := $(shell id -g)

COMPOSE := docker compose -f ./deploy/docker-compose.dev.yml

# Display the root-level development command reference.
.PHONY: help
help:
	@printf '\n\033[1;36m╔══════════════════════════════════════════════════════════════════╗\033[0m\n'
	@printf '\033[1;36m║                     TripMate · Development                      ║\033[0m\n'
	@printf '\033[1;36m╚══════════════════════════════════════════════════════════════════╝\033[0m\n'
	@printf '\n\033[1;33mUSAGE\033[0m\n'
	@printf '  \033[1;32mmake\033[0m or \033[1;32mmake help\033[0m       Show this command reference.\n'
	@printf '  Run service-specific commands from services/app or services/route_generator.\n'
	@printf '  Run contract commands from contracts/.\n'
	@printf '\n\033[1;33mAPPLICATION SERVICE\033[0m\n'
	@printf '  \033[1;32mapp-run\033[0m                         Build and start the App service.\n'
	@printf '  \033[1;32mapp-postgres-up\033[0m                Start the App PostgreSQL container.\n'
	@printf '  \033[1;32mapp-postgres-down\033[0m              Stop and remove the App PostgreSQL container.\n'
	@printf '  \033[1;32mapp-postgres-cleanup\033[0m           Remove the App database volume after confirmation.\n'
	@printf '  \033[1;32mapp-postgres-port-forward\033[0m      Start the App PostgreSQL port-forwarder.\n'
	@printf '  \033[1;32mapp-postgres-port-close\033[0m        Stop the App PostgreSQL port-forwarder.\n'
	@printf '\n\033[1;33mROUTE GENERATOR SERVICE\033[0m\n'
	@printf '  \033[1;32mroute-generator-run\033[0m            Build and start the Route Generator service.\n'
	@printf '  \033[1;32mroute-generator-postgres-up\033[0m   Start the Route Generator PostgreSQL container.\n'
	@printf '  \033[1;32mroute-generator-postgres-down\033[0m Stop and remove its PostgreSQL container.\n'
	@printf '  \033[1;32mroute-generator-postgres-cleanup\033[0m Remove its database volume after confirmation.\n'
	@printf '  \033[1;32mroute-generator-postgres-port-forward\033[0m Start its PostgreSQL port-forwarder.\n'
	@printf '  \033[1;32mroute-generator-postgres-port-close\033[0m Stop its PostgreSQL port-forwarder.\n'
	@printf '\n\033[1;33mEXAMPLES\033[0m\n'
	@printf '  make app-run\n  make app-postgres-up\n'
	@printf '  make route-generator-run\n  make route-generator-postgres-up\n'
	@printf '\n\033[1;90mTip: run make help inside a component directory for its full command list.\033[0m\n\n'

# Build the App image if needed and stream its container logs in the foreground.
.PHONY: app-run
app-run:
	@$(COMPOSE) up --build $(APP_SERVICE_NAME)

# Start the App PostgreSQL container in the background.
.PHONY: app-postgres-up
app-postgres-up:
	@$(COMPOSE) up -d app-postgres

# Stop and remove the App PostgreSQL container without deleting its data volume.
.PHONY: app-postgres-down
app-postgres-down:
	@$(COMPOSE) rm -sf app-postgres

# Ask for confirmation, then remove the App database volume (destructive).
.PHONY: app-postgres-cleanup
app-postgres-cleanup: app-postgres-down
	@read -r -p "Delete the App PostgreSQL data volume? This cannot be undone. [y/N] " answer; \
	if [[ "$$answer" == "y" || "$$answer" == "Y" ]]; then \
		$(COMPOSE) rm -sf app-postgres-port-forwarder && \
		docker volume rm "$(APP_SERVICE_NAME)_pgdata"; \
	else \
		printf '\033[1;33mCancelled; database volume was preserved.\033[0m\n'; \
	fi

# Start the App PostgreSQL port-forwarder in the background.
.PHONY: app-postgres-port-forward
app-postgres-port-forward:
	@$(COMPOSE) up -d app-postgres-port-forwarder

# Stop and remove the App PostgreSQL port-forwarder.
.PHONY: app-postgres-port-close
app-postgres-port-close:
	@$(COMPOSE) rm -sf app-postgres-port-forwarder

# Build the Route Generator image if needed and stream its logs in the foreground.
.PHONY: route-generator-run
route-generator-run:
	@$(COMPOSE) up --build $(ROUTE_GENERATOR_SERVICE_NAME)

# Start the Route Generator PostgreSQL container in the background.
.PHONY: route-generator-postgres-up
route-generator-postgres-up:
	@$(COMPOSE) up -d route-generator-postgres

# Stop and remove the Route Generator PostgreSQL container without deleting its data volume.
.PHONY: route-generator-postgres-down
route-generator-postgres-down:
	@$(COMPOSE) rm -sf route-generator-postgres

# Ask for confirmation, then remove the Route Generator database volume (destructive).
.PHONY: route-generator-postgres-cleanup
route-generator-postgres-cleanup: route-generator-postgres-down
	@read -r -p "Delete the Route Generator PostgreSQL data volume? This cannot be undone. [y/N] " answer; \
	if [[ "$$answer" == "y" || "$$answer" == "Y" ]]; then \
		$(COMPOSE) rm -sf route-generator-postgres-port-forwarder && \
		docker volume rm "$(ROUTE_GENERATOR_VOLUME_NAME)_pgdata"; \
	else \
		printf '\033[1;33mCancelled; database volume was preserved.\033[0m\n'; \
	fi

# Start the Route Generator PostgreSQL port-forwarder in the background.
.PHONY: route-generator-postgres-port-forward
route-generator-postgres-port-forward:
	@$(COMPOSE) up -d route-generator-postgres-port-forwarder

# Stop and remove the Route Generator PostgreSQL port-forwarder.
.PHONY: route-generator-postgres-port-close
route-generator-postgres-port-close:
	@$(COMPOSE) rm -sf route-generator-postgres-port-forwarder
