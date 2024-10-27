# Stage 1: Build the Go binary
FROM golang:1.18-alpine AS builder

# Set working directory in container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o myapp .

# Stage 2: Run the Go binary
FROM alpine:latest

# Set working directory in container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/myapp .

# Expose the port your app runs on (adjust as necessary)
EXPOSE 50002

# Run the Go application
CMD ["./myapp"]
