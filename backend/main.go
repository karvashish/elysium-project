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
	setupConfig()
	setupDatabase()
	setupWireGuard()
	startServer()
}

func setupConfig() {
	envFilePath := flag.String("env", "../.env", "Path to the env file")
	flag.Parse()
	config.LoadEnv(*envFilePath)
}

func setupDatabase() {
	db.DBPool = db.InitializeDatabaseConnection()
	migrationDir := config.GetEnv("MIGRATION_PATH", "migrations")

	if err := db.RunMigrations(migrationDir); err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}
}

func setupWireGuard() {
	serverInterface := config.GetEnv("BACKEND_WG_INTERFACE", "wg0")
	serverPort, err := strconv.Atoi(config.GetEnv("BACKEND_WG_PORT", "51820"))
	if err != nil {
		log.Fatalf("Invalid port provided: %v", err)
	}

	serverIP := config.GetEnv("BACKEND_WG_IP", "10.0.0.1")
	networkMask := config.GetEnv("WG_NETWORK_MASK", "/24")

	if err := wgutil.InitWireGuardInterface(serverInterface, serverPort, serverIP, networkMask); err != nil {
		log.Fatalf("Failed to set up WireGuard network: %v", err)
	}
}

func startServer() {
	port := ":" + config.GetEnv("PORT", "8080")
	server := &http.Server{Addr: port, Handler: routes.SetupRoutes()}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to listen and serve: %v", err)
	}

	defer db.CloseDatabaseConnection()
}
