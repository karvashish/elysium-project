@echo off

cd containers
podman-compose down
cd ..

podman volume prune -f

echo Teardown completed successfully.
