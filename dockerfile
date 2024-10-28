# Stage 1: Build the Go binary
FROM golang:1.23 AS builder
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/booking-service-GOLANG/main.go

# Stage 2: Run the Go binary with configuration
FROM alpine:latest
WORKDIR /app

# Copy the binary and entrypoint script from the builder stage
COPY --from=builder /app/main .
COPY entrypoint.sh /app/entrypoint.sh

# Ensure the entrypoint script is executable
RUN chmod +x /app/entrypoint.sh

# Set environment variables (these are passed in from Kubernetes ConfigMap)
ARG DB_HOST ARG DB_NAME ARG DB_PORT ARG DB_PW ARG DB_UN ARG HTTP_ADDRESS ARG HTTP_PORT ARG API_GATEWAY
ENV DB_HOST=${DB_HOST} DB_NAME=${DB_NAME} DB_PORT=${DB_PORT} DB_PW=${DB_PW} DB_UN=${DB_UN}
ENV HTTP_ADDRESS=${HTTP_ADDRESS} HTTP_PORT=${HTTP_PORT} API_GATEWAY=${API_GATEWAY}

# Expose the HTTP port
EXPOSE 50002

# Test to ensure the binary can execute in Alpine
RUN ./main --help || true

# Run the entrypoint script to generate the production.yaml and start the application
ENTRYPOINT ["/app/entrypoint.sh"]
