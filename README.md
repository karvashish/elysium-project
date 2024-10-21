## 1. Project Overview
**Elysium** automates secure peer-to-peer connections using **WireGuard**. It handles peer onboarding, key generation, and interface setup via containers for **PostgreSQL**, **Redis**, and **OpenSSL**. Security is enforced through mTLS and token-based authentication.

## 2. Core Components
- **Backend (Go)**: Controls WireGuard configuration, peer management, and communications.
- **Client (Rust)**: Interacts with the backend for establishing secure connections.
- **PostgreSQL**: Stores peer data, tokens, and network metadata.
- **Redis**: Manages workflows like peer inactivity and token expiration checks.

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
