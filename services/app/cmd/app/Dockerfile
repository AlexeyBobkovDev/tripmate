FROM golang:1.27.3-bookworm AS builder

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

FROM golang:1.27.3-alpine
WORKDIR /app

RUN addgroup -S -g 10001 app && \
    adduser -S \
            -u 10001 \
            -G app \
            -H \
            -s /sbin/nologin app \
            app

COPY --from=builder --chown=10001:10001 /app/cmd/app/exe /app/

USER 10001:10001

ENV PROJECT_ROOT=/app
ENTRYPOINT ["/app/exe"]
