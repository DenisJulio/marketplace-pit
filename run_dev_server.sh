#!/bin/bash

# Function to stop the container
stop_container() {
    echo "Stopping compose project..."
    docker compose stop
}

echo "Starting compose project..."
docker compose start

# Trap SIGINT and call the stop_container function
trap stop_container SIGINT

# Start the dev-server
echo "Starting dev-server with air..."
air

# If the dev-server exits normally, stop the container
stop_container
