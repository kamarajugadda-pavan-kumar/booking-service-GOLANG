# Stage 1: Build the Go binary
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/booking-service-GOLANG/main.go

# Stage 2: Run the Go binary with configuration
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY entrypoint.sh /app/entrypoint.sh

# Create config directory
RUN mkdir -p /app/config && chmod +x /app/entrypoint.sh

# Pass the build arguments to runtime environment
ARG DB_HOST ARG DB_NAME ARG DB_PORT ARG DB_PW ARG DB_UN
ENV DB_HOST=${DB_HOST} DB_NAME=${DB_NAME} DB_PORT=${DB_PORT} DB_PW=${DB_PW} DB_UN=${DB_UN} 

# Expose the HTTP port
EXPOSE 50002

ENTRYPOINT ["/app/entrypoint.sh"]