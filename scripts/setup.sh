#!/bin/bash

echo "Choose the container engine:"
echo "1. Docker Compose"
echo "2. Podman Compose"
read -p "Enter your choice (1 or 2): " choice

if [ "$choice" == "1" ]; then
    engine="docker-compose"
elif [ "$choice" == "2" ]; then
    engine="podman-compose"
else
    echo "Invalid choice. Defaulting to Podman Compose."
    engine="podman-compose"
fi

cd containers
$engine up -d
cd ..

echo "Setup completed successfully."
