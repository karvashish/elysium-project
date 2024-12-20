#!/bin/bash

echo "Choose the container engine:"
echo "1. Docker Compose"
echo "2. Podman Compose"
read -p "Enter your choice (1 or 2): " choice

if [ "$choice" == "1" ]; then
    engine="docker compose"
    cli="docker"
elif [ "$choice" == "2" ]; then
    engine="podman-compose"
    cli="podman"
else
    echo "Invalid choice. Defaulting to Podman Compose."
    engine="podman-compose"
    cli="podman"
fi

cd containers
sudo $engine -f podman-compose.yml down --volumes --remove-orphans
cd ..

echo "Removing all containers..."
sudo $cli rm -f $(sudo $cli ps -aq)

echo "Removing all volumes..."
sudo $cli volume rm -f $(sudo $cli volume ls -q)

echo "Removing all images..."
sudo $cli rmi -f $(sudo $cli images -q)

echo "Removing all networks..."
sudo $cli network rm $(sudo $cli network ls -q)

if [ "$cli" == "docker" ]; then
    echo "Cleaning up Docker build cache..."
    sudo $cli builder prune -a -f
fi

echo "Complete cleanup completed successfully."
