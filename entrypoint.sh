#!/bin/sh

# Ensure the config directory exists
mkdir -p 

# Create or overwrite the production.yaml file in the config folder
cat <<EOL > production.yaml
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

echo "Production config file created/overwritten at production.yaml:"
cat production.yaml

# Start the Golang application
./main --config production.yaml
