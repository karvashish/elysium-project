@echo off

REM Step 1: Build and Run OpenSSL Dockerfile
cd ./containers/openssl

if not exist ../../certificates mkdir ../../certificates

podman build -t minimal_openssl_image .
podman run --rm -v "%cd%\..\..\certificates:/certificates" minimal_openssl_image

REM Step 2: Start all services with Podman Compose
cd ..
podman-compose up -d
