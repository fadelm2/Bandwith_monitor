# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the remaining application source code
COPY . .

# Build the Go application
# Pointing to cmd/app/main.go as the entry point
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app/main.go

# Stage 2: Create a small production image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy the configuration file (optional if using environment variables)
COPY config.json .

# Expose the application port (based on config.json)
EXPOSE 8082

# Run the application
CMD ["./main"]
