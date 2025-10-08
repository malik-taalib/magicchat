# ---------- Stage 1: Build ----------
FROM golang:1.23-alpine AS builder

# Install build deps
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files and download dependencies first (for caching)
COPY go.mod go.sum ./
ENV GOTOOLCHAIN=auto
RUN go mod download

# Copy the rest of the backend code
COPY . .

# Build the binary
ENV GOTOOLCHAIN=auto
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/server ./cmd/server

# ---------- Stage 2: Run ----------
FROM alpine:3.20

# Create non-root user
RUN adduser -D -g '' appuser

# Set up working directory
WORKDIR /home/appuser

# Copy built binary and any required files (like config)
COPY --from=builder /bin/server /usr/local/bin/server

# Change ownership
USER appuser

# Expose port for API
EXPOSE 8080

# Command
ENTRYPOINT ["/usr/local/bin/server"]
