# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/expense-tracker

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/expense-tracker .

# Expose the application port
EXPOSE 8080

# Run the application with database flag
CMD ["./expense-tracker", "-db"]
