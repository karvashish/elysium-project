@echo off

cd containers
podman-compose down
cd ..

podman volume prune -f

del /f /q containers\vault\vault_init_keys.txt
echo Teardown completed successfully.
