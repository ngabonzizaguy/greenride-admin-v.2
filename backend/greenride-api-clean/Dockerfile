###############################################
# Stage 1: Build the Go binary
###############################################
FROM golang:1.23-alpine AS builder

RUN apk --no-cache add git ca-certificates tzdata

WORKDIR /build

# Cache dependency downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build a statically-linked Linux binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o greenride-api ./main

###############################################
# Stage 2: Minimal runtime image
###############################################
FROM alpine:latest

# Install certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Set timezone
ENV TZ=Africa/Kigali

# Create non-root user
RUN addgroup -g 1001 appgroup && adduser -u 1001 -G appgroup -s /bin/sh -D appuser

# Setup directories
RUN mkdir -p /app/logs /logs

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /build/greenride-api ./greenride-api

# Copy required static assets and configs
COPY internal/locales ./internal/locales
COPY config.yaml ./config.yaml
# Copy prod.yaml if it exists (optional at build stage since compose mounts it)
COPY prod.yaml* ./prod.yaml

# Set permissions
RUN chmod +x ./greenride-api && \
    mkdir -p /app/logs && \
    chown -R appuser:appgroup /app && \
    chmod -R 777 /app/logs

USER appuser

# API ports
EXPOSE 8610 8611 8612

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8610/health || exit 1

# Start the application
CMD ["./greenride-api", "serve"]
