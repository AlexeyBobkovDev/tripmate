# ===========================Buf===========================
.PHONY=buf-lint
buf-lint:
	@buf lint

.PHONY=buf-fmt
buf-fmt:
	@buf format -w

.PHONY=buf-generate
buf-gen:
	@rm -rf gen/* && \
	buf generate

.PHONY=buf-breaking
buf-breaking:
	@buf breaking --against '../.git#branch=main,subdir=contracts/proto'

.PHONY=buf-build
buf-build:
	@buf build


# ===========================Go===========================
.PHONY: go-tidy
go-tidy:
	@go mod tidy

.PHONY: go-fmt
go-fmt:
	@for file in $$(find . -type f -name '*.go'); do \
		goimports -w "$$file" && \
		gofumpt -w "$$file" && \
		gci write \
			--custom-order \
			-s standard \
			-s default \
			-s "prefix(github.com/AlexeyBobkovDev/tripmate)" \
			"$$file"; \
	done

.PHONY: go-lint
go-lint:
	@golangci-lint run ./...
