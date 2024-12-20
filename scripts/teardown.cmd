@echo off
echo Choose the container engine:
echo 1. Docker Compose
echo 2. Podman Compose
set /p choice=Enter your choice (1 or 2): 

if "%choice%"=="1" (
    set engine=docker-compose
    set cli=docker
) else if "%choice%"=="2" (
    set engine=podman-compose
    set cli=podman
) else (
    echo Invalid choice. Defaulting to Podman Compose.
    set engine=podman-compose
    set cli=podman
)

cd containers
%engine% -f podman-compose.yml down --volumes --remove-orphans
cd ..

echo Removing all containers...
for /f %%i in ('%cli% ps -aq') do %cli% rm -f %%i

echo Removing all volumes...
for /f %%i in ('%cli% volume ls -q') do %cli% volume rm -f %%i

echo Removing all images...
for /f %%i in ('%cli% images -q') do %cli% rmi -f %%i

echo Removing all networks...
for /f %%i in ('%cli% network ls -q') do %cli% network rm %%i

if "%cli%"=="docker" (
    echo Cleaning up Docker build cache...
    %cli% builder prune -a -f
)

echo Complete cleanup completed successfully.
