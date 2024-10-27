@echo off

REM Step 1: Start all services with Podman Compose
cd containers
podman-compose up -d
cd ..

REM Step 2: Check if Vault is already initialized; initialize if not
podman exec containers_vault_1 vault status | find "Initialized" | find "true" >nul
if errorlevel 1 (
    echo Initializing Vault...
    podman exec containers_vault_1 vault operator init -key-shares=3 -key-threshold=2 > vault_init_keys.txt
    echo Vault initialized. Unseal keys and root token saved to vault_init_keys.txt
) else (
    echo Vault is already initialized.
)

REM Step 3: Unseal Vault with saved keys (first-time unseal)
for /F "tokens=5" %%i in ('findstr /c:"Unseal Key" vault_init_keys.txt') do (
    podman exec containers_vault_1 vault operator unseal %%i
)

echo Setup completed successfully.
