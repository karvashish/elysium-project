@echo off

:: Present choices to the user
echo Choose the container engine:
echo 1. Docker Compose
echo 2. Podman Compose
set /p choice=Enter your choice (1 or 2): 

:: Validate input
if "%choice%"=="1" (
    set engine=docker-compose
) else if "%choice%"=="2" (
    set engine=podman-compose
) else (
    echo Invalid choice. Defaulting to Podman Compose.
    set engine=podman-compose
)

cd containers
%engine% down
cd ..

if "%engine%"=="podman-compose" (
    podman volume prune -f
)

echo Teardown completed successfully.
