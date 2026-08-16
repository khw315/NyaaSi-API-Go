# ──────────────────────────────────────────────────────────────
# NyaaSi-API Go — Docker image
# ──────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o nyaa-api .

# ──────────────────────────────────────────────────────────────
# Final stage: minimal runtime environment
# ──────────────────────────────────────────────────────────────
FROM alpine:3.20

WORKDIR /app

LABEL org.opencontainers.image.source=https://github.com/khw315/NyaaSi-API-Go
LABEL org.opencontainers.image.licenses=GPL-3.0
LABEL org.opencontainers.image.description="API for nyaa.si and sukebei.nyaa.si built with Go"

# Install ca-certificates and create a non-root user
RUN apk add --no-cache ca-certificates && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy binary from builder
COPY --from=builder /app/nyaa-api /app/nyaa-api

USER appuser

EXPOSE 88

CMD ["/app/nyaa-api"]
