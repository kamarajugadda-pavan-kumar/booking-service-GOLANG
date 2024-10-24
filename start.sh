#!/bin/sh

# Environment variables injected from Kubernetes
# Example variables: DB_HOST, DB_PORT, DB_USERNAME, DB_PASSWORD, DB_NAME, HTTP_PORT, HTTP_ADDRESS, API_GATEWAY

# Create or overwrite the production.yaml file in the config folder
cat <<EOL > /app/config/production.yaml
env: "production"

database:
  name: "${DB_NAME}"
  host: "${DB_HOST}"
  port: "${DB_PORT}"
  username: "${DB_UN}"
  password: "${DB_PW}"

http_server:
  address: "${HTTP_ADDRESS}"
  port: "${HTTP_PORT}"

api_gateway: "${API_GATEWAY}"
EOL

echo "Production config file created/overwritten at /app/config/production.yaml:"
cat /app/config/production.yaml

# Start the Golang application
./cmd/booking-service-GOLANG/main
