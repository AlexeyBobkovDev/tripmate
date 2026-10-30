FROM golang:1.27.3-bookworm AS builder

LABEL author="Alexey Bobkov"
LABEL email="alexeybobkovdev@outlook.com"

WORKDIR /workspace

COPY services/route_generator/go.mod services/route_generator/go.sum /workspace/services/route_generator/
COPY contracts/go.mod /workspace/contracts/

RUN --mount=type=cache,target=/go/pkg/mod \
    cd /workspace/services/route_generator && go mod download &&\
    cd /workspace/contracts && go mod download

COPY . .

WORKDIR /workspace/services/route_generator

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o ./out/exe ./cmd/route_generator

FROM golang:1.27.3-alpine

WORKDIR /workspace

RUN addgroup -S -g 10001 route_generator && \
    adduser -S \
            -u 10001 \
            -G route_generator \
            -H \
            -s /sbin/nologin route_generator \
            route_generator

COPY --from=builder --chown=10001:10001 /workspace/services/route_generator/out/exe /workspace/

USER 10001:10001

ENV PROJECT_ROOT=/workspace
ENTRYPOINT ["/workspace/exe"]
