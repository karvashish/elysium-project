@echo off

cd containers
podman-compose down

REM Remove the OpenSSL image
podman rmi minimal_openssl_image

REM Remove volumes created by Podman Compose
podman volume prune -f

REM Clean up certificates folder
cd ../certificates
rmdir /s /q root
rmdir /s /q intermediate

echo Teardown completed successfully.
