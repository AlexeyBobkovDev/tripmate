FROM golang:1.27.3-bookworm AS builder

LABEL author="Alexey Bobkov"
LABEL email="alexeybobkovdev@outlook.com"

WORKDIR /workspace

COPY services/app/go.mod services/app/go.sum /workspace/services/app/
COPY contracts/go.mod /workspace/contracts/

RUN --mount=type=cache,target=/go/pkg/mod \
    cd /workspace/services/app && go mod download &&\
    cd /workspace/contracts && go mod download

COPY . .

WORKDIR /workspace/services/app

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o ./out/exe ./cmd/app

FROM golang:1.27.3-alpine

WORKDIR /workspace

RUN addgroup -S -g 10001 app && \
    adduser -S \
            -u 10001 \
            -G app \
            -H \
            -s /sbin/nologin app \
            app

COPY --from=builder --chown=10001:10001 /workspace/services/app/out/exe /workspace/

USER 10001:10001

ENV PROJECT_ROOT=/workspace
ENTRYPOINT ["/workspace/exe"]
