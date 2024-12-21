
# Step-by-Step Guide

## Clone the Repository:

git clone \<repository-url\>
cd \<repository-directory\>

## Step 1: Verify Environment Configuration
- **Description**: Check the `.env` file in the `project_root` directory.
- **Goal**: Ensure the configuration values are correct for variables such as `DB_HOST`, `DB_USER`, `DB_PASSWORD`, etc.

---

## Step 2: Initialize Services
- **Description**: Run the setup script to start all required services in containers.
- **Details**: This command will initialize containers for PostgreSQL, and other dependencies as defined.
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
