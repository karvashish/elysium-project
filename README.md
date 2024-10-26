## 1. Project Overview
**Elysium** automates secure peer-to-peer connections using **WireGuard**. It handles peer onboarding, key generation, and interface setup via containers for **PostgreSQL**, **Redis**, and **OpenSSL**. Security is enforced through mTLS and token-based authentication.

## 2. Core Components
- **Backend (Go)**: Controls WireGuard configuration, peer management, and communications.
- **Client (Rust)**: Interacts with the backend for establishing secure connections.
- **PostgreSQL**: Stores peer data, tokens, and network metadata.
- **Redis**: Manages workflows like peer inactivity and token expiration checks.

Step-by-Step Guide

Clone the Repository:

git clone <repository-url>
cd <repository-directory>

# Environment Setup Instructions

## Step 1: Verify Environment Configuration
- **Description**: Check the `.env` file in the `project_root` directory.
- **Goal**: Ensure the configuration values are correct for variables such as `DB_HOST`, `DB_USER`, `DB_PASSWORD`, etc.

---

## Step 2: Initialize Services
- **Description**: Run the setup script to start all required services in containers.
- **Details**: This command will initialize containers for PostgreSQL, Redis, OpenSSL, and other dependencies as defined.
```
./scripts/setup.cmd
```

---

## Step 3: Build the Backend Docker Image
- **Description**: Build the Docker image for the backend server using Podman.
- **Location**: Run this command from `project_root`.
```
podman build -t elysium-backend -f backend/Dockerfile .
```

---

## Step 4: Run the Backend Server
- **Description**: Start the backend server in a container on the same network as other services.
- **Goal**: Maps the backend server to port `8080` and connects it to other containers on the `containers_default` network.
```
podman run --rm --cap-add=NET_ADMIN --cap-add=SYS_MODULE --network=containers_default -p 8080:8080 elysium-backend
```




## 3. Structure Overview
- **Backend**: Initializes network services, manages API, database, and WireGuard through programmatic interfaces like `vishvananda/netlink` for creating interfaces and `wgtypes` for key management.
  
- **Client**: A Rust-based tool that securely connects peers to the hub using WireGuard parameters retrieved from the backend.

- **Containers**: Podman Compose setup for PostgreSQL, Redis, and OpenSSL services.

## 4. WireGuard Configuration
- **Interface Creation**: WireGuard interfaces (`wg0`, `wg1`) are created dynamically.
- **Key Generation**: Private keys are generated via `GeneratePrivateKey()` and distributed securely to peers.
- **Network Setup**: IP address assignment and routing are automated for each peer. The backend handles route updates as peers join/leave.

## 5. Automation Scripts
- **Setup**: Automates OpenSSL container build, certificate creation, container orchestration, and database initialization.
- **Teardown**: Stops all services, removes containers and volumes, and cleans up certificates.

## 6. Deployment Process
- **Containerized Services**: Podman Compose deploys PostgreSQL, Redis, and OpenSSL in isolated containers.
- **Dynamic Workflow Management**: Redis queues workflows (e.g., checking peer status), while PostgreSQL holds persistent data (peer, token, audit logs).
- **WireGuard Management**: Backend dynamically manages WireGuard interfaces, peer connections, and communication security using API tokens.
