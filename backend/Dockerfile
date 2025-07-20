# Base stage for both development and production
FROM golang:1.24-alpine AS base

# Install git, ca-certificates and other tools needed for development
RUN apk add --no-cache git ca-certificates bash curl

# Set working directory
WORKDIR /app

# Development stage
FROM base AS development

# Install inotifyd for file watching
RUN apk add --no-cache inotify-tools

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# We'll mount the rest of the code as a volume in docker-compose

# Set environment variables for development
ENV GO_ENV=development

# The dev.sh script will be mounted from the host
CMD ["sh", "/app/dev.sh"]

# Build stage
FROM base AS builder

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Production stage
FROM alpine:latest AS production

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .

# Copy config files from builder stage
COPY --from=builder /app/configs ./configs

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Set environment variables for production
ENV GO_ENV=production

# Run the application
CMD ["./main"] 