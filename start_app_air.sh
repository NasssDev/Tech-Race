#!/bin/bash

# Function to print usage instructions
print_usage() {
    echo "Usage: $0"
    echo "Runs docker compose up -d after ensuring ports are free"
}

# Check if docker-compose is available
if ! command -v docker compose &> /dev/null; then
    echo "Error: docker-compose is not installed"
    exit 1
fi

# Stop and disable mosquitto service
echo "Stopping mosquitto service..."
sudo systemctl stop mosquitto


# Kill processes on ports 1883 and 5432
echo "Killing processes on ports 1883 and 5432..."
sudo kill -9 $(sudo lsof -t -i:1883) 2>/dev/null || true
sudo kill -9 $(sudo lsof -t -i:5432) 2>/dev/null || true

# Start docker-compose
echo "Starting docker-compose..."
docker compose up -d

# Check if docker-compose started successfully
if [ $? -eq 0 ]; then
    echo "Docker Compose started successfully."
else
    echo "Error: Docker Compose failed to start. Exiting."
    exit 1
fi

npm run build
air