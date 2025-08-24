# --- Stage 1: Build Go app ---
FROM golang:1.23.4-alpine AS builder

# Install git (needed for Go modules)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files and download dependencies first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the go binary
RUN go build -o myapp .

# --- Stage 2: Runtime ---
FROM alpine:latest

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates     # ← Fixed: "add" not "and"

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/myapp .

# Expose the port your app uses
EXPOSE 2112

# Command to run
CMD ["./myapp"]