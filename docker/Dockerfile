# Stage 1: Build the Go binary
FROM golang:1.25 AS builder

WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o kart-challenge ./cmd/app/main.go

# Stage 2: Create a minimal image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/kart-challenge .

# Copy the configuration file
COPY config/config.docker.yaml .
COPY db /app/db

# Expose the port (optional)
EXPOSE 8080

# Run the binary
ENTRYPOINT exec ./kart-challenge --config config.docker.yaml
