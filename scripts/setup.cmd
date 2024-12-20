@echo off

cd containers
podman-compose up -d
cd ..

echo Setup completed successfully.
