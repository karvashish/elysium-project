package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"elysium-backend/config"
	"elysium-backend/internal/routes"
	"elysium-backend/pkg/db"
	"elysium-backend/pkg/wgutil"
)

func main() {
	envFilePath := flag.String("env", "../.env", "Path to the env file")
	flag.Parse()

	config.LoadEnv(*envFilePath)
	db.DBPool = db.InitializeDatabaseConnection()

	migrationDir := config.GetEnv("MIGRATION_PATH", "migrations")
	err := db.RunMigrations(migrationDir)
	if err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}

	server_interface := config.GetEnv("BACKEND_WG_INTERFACE", "wg0")
	server_port, err := strconv.Atoi(config.GetEnv("BACKEND_WG_PORT", "51820"))
	if err != nil {
		log.Fatalf("invalid Port Provided: %v", err)
	}
	server_IP := config.GetEnv("BACKEND_WG_IP", "10.0.0.1")
	network_mask := config.GetEnv("WG_NETWORK_MASK", "/24")

	if err := wgutil.InitWireGuardInterface(server_interface, server_port, server_IP, network_mask); err != nil {
		log.Fatalf("Failed setup wireguard network: %v", err)
	}

	port := ":" + config.GetEnv("PORT", "8080")

	server := &http.Server{Addr: port, Handler: routes.SetupRoutes()}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to listen and serve: %v", err)
	}

	db.CloseDatabaseConnection()
}
