# --- Stage 1: Build --- #
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Download dependencies first (layer cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a statically linked binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o /out/server \
    ./cmd/server/main.go

# --- Stage 2: Runtime --- #
FROM alpine:3.21 AS runtime

RUN apk add --no-cache wget && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app

# Copy binary and SQL migrations
COPY --from=builder /out/server  ./server
COPY --from=builder /app/migrations ./migrations

EXPOSE 7129

HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
  CMD wget -qO- http://localhost:7129/health || exit 1

USER appuser

ENTRYPOINT ["/app/server"]