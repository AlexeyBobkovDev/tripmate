# TripMate — Protocol Buffer contracts and generated Go code
.DEFAULT_GOAL := help

SHELL := /bin/bash

# Display contract-generation, validation, and Go maintenance commands.
.PHONY: help
help:
	@printf '\n\033[1;36m╔══════════════════════════════════════════════════════════════╗\033[0m\n'
	@printf '\033[1;36m║                 TripMate · Contracts                         ║\033[0m\n'
	@printf '\033[1;36m╚══════════════════════════════════════════════════════════════╝\033[0m\n'
	@printf '\n\033[1;33mUSAGE\033[0m\n'
	@printf '  \033[1;32mmake\033[0m or \033[1;32mmake help\033[0m       Show this command reference.\n'
	@printf '\n\033[1;33mPROTOBUF / BUF\033[0m\n'
	@printf '  \033[1;32mbuf-lint\033[0m                      Lint .proto files and enforce Buf rules.\n'
	@printf '  \033[1;32mbuf-fmt\033[0m                       Format Protocol Buffer definitions in place.\n'
	@printf '  \033[1;32mbuf-generate\033[0m                  Regenerate Go code from the .proto definitions.\n'
	@printf '  \033[1;32mbuf-breaking\033[0m                  Check compatibility against contracts on main.\n'
	@printf '  \033[1;32mbuf-build\033[0m                     Compile and validate the Buf module.\n'
	@printf '\n\033[1;33mGO MODULE\033[0m\n'
	@printf '  \033[1;32mgo-tidy\033[0m                       Clean up and normalize go.mod / go.sum.\n'
	@printf '  \033[1;32mgo-fmt\033[0m                        Format and organize imports in Go files.\n'
	@printf '  \033[1;32mgo-lint\033[0m                       Run golangci-lint across the contracts module.\n'
	@printf '\n\033[1;33mEXAMPLES\033[0m\n'
	@printf '  make buf-lint && make buf-breaking\n'
	@printf '  make buf-fmt && make buf-generate\n\n'

# Lint Protocol Buffer definitions using the repository's Buf configuration.
.PHONY: buf-lint
buf-lint:
	@buf lint

# Apply the configured Buf formatting rules to all .proto files.
.PHONY: buf-fmt
buf-fmt:
	@buf format -w

# Remove stale generated files and regenerate Go bindings from the .proto sources.
.PHONY: buf-generate buf-gen
buf-generate buf-gen:
	@rm -rf gen/* && buf generate

# Detect backward-incompatible changes relative to the contracts on main.
.PHONY: buf-breaking
buf-breaking:
	@buf breaking --against '../.git#branch=main,subdir=contracts/proto'

# Build the Buf module to validate imports, schemas, and dependency resolution.
.PHONY: buf-build
buf-build:
	@buf build

# Reconcile this module's dependencies and update go.mod / go.sum.
.PHONY: go-tidy
go-tidy:
	@go mod tidy

# Format Go source and organize imports using the project's import grouping.
.PHONY: go-fmt
go-fmt:
	@set -e; \
	find . -type f -name '*.go' -print0 | while IFS= read -r -d '' file; do \
		goimports -w "$$file"; \
		gofumpt -w "$$file"; \
		gci write --custom-order -s standard -s default -s "prefix(github.com/AlexeyBobkovDev/tripmate)" "$$file"; \
	done

# Run golangci-lint against every package in the contracts Go module.
.PHONY: go-lint
go-lint:
	@golangci-lint run ./...
