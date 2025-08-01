# Build stage
FROM golang:1.24-alpine AS builder

# Install git and other dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user
RUN addgroup -g 1001 appgroup && \
    adduser -D -s /bin/sh -u 1001 -G appgroup appuser

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Create directory for any potential config files
RUN mkdir -p /root/config

# Change ownership to non-root user
RUN chown -R appuser:appgroup /root

# Switch to non-root user
USER appuser

# Expose port 8080
EXPOSE 8080

# Set environment variables (can be overridden at runtime)
ENV GIN_MODE=release
ENV PORT=8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# Command to run the application
CMD ["./main"]