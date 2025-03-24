# Build stage
FROM golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -o task-management-api ./cmd/api

# Final stage
FROM alpine:3.19

# Add ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/task-management-api .

# Copy configuration files
COPY config/config.yaml /app/config/config.yaml
COPY config/rbac_model.conf /app/config/rbac_model.conf

# Expose API port
EXPOSE 8080

# Run the application
CMD ["./task-management-api"] 