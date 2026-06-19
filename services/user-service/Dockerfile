# ---- Build Stage ----
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies (needed for some Go packages)
RUN apk add --no-cache git

# Copy go mod files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy full source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o user-service ./cmd/main.go


# ---- Run Stage (small image) ----
FROM alpine:latest

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/user-service .

# Expose service port (change if needed)
EXPOSE 8080

# Run the app
CMD ["./user-service"]