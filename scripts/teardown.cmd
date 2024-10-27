@echo off

cd containers
podman-compose down

REM Remove volumes created by Podman Compose
podman volume prune -f

echo Teardown completed successfully.
