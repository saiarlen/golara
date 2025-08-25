# Multi-stage build for production
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o golara main.go

# Production stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create app user for security
RUN addgroup -g 1001 -S golara && \
    adduser -S golara -u 1001 -G golara

WORKDIR /app

# Copy binary and required files
COPY --from=builder /app/golara .
COPY --from=builder /app/.env.yaml* ./
COPY --from=builder /app/resources ./resources/
COPY --from=builder /app/storage ./storage/

# Create necessary directories
RUN mkdir -p storage/logs storage/cache storage/app && \
    chown -R golara:golara /app

# Switch to app user
USER golara

# Expose port
EXPOSE 9000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:9000/health || exit 1

CMD ["./golara"]


