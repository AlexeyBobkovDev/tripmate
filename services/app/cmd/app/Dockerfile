FROM golang:1.26.4-bookworm AS builder

LABEL author="Alexey Bobkov"
LABEL email="alexeybobkovdev@outlook.com"

WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o ./cmd/app/exe ./cmd/app

FROM golang:1.26.4-alpine
WORKDIR /app
COPY --from=builder /app/cmd/app/exe /app/
ENV PROJECT_ROOT=/app
ENTRYPOINT ["/app/exe"]
