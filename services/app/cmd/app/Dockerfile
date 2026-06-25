FROM golang:1.26.4-bookworm AS builder

LABEL author="Alexey Bobkov"
LABEL email="alexeybobkovdev@outlook.com"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o ./cmd/app/exe ./cmd/app

FROM golang:1.26.4-alpine
WORKDIR /app
COPY --from=builder /app/cmd/app/exe /app/
ENV PROJECT_ROOT=/app
ENTRYPOINT ["/app/exe"]
