@echo off

REM Start all services with Podman Compose
cd ./containers
podman-compose up -d
