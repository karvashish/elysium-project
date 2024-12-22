package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"strconv"

	"elysium-backend/config"
	"elysium-backend/internal/routes"
	"elysium-backend/pkg/db"
	"elysium-backend/pkg/wgutil"
)

func main() {

	envFilePath := flag.String("env", "../local.env", "Path to the env file")
	setupWg := flag.Bool("setupWg", true, "Setup wireguard network")
	flag.Parse()
	config.LoadEnv(*envFilePath)
	log.Println("main.setupConfig -> configuration loaded")

	setupDatabase()
	setupWireGuard(setupWg)
	startServer()
}

func setupDatabase() {
	if config.GetLogLevel() == "DEBUG" {
		log.Println("main.setupDatabase -> called")
	}
	db.DBPool = db.InitializeDatabaseConnection()
	migrationDir := config.GetEnv("MIGRATION_PATH", "migrations")

	if err := db.RunMigrations(migrationDir); err != nil {
		log.Fatalf("main.setupDatabase -> migrations failed: %v", err)
	}
	log.Println("main.setupDatabase -> database setup complete")
}

func setupWireGuard(setupWg *bool) {
	if config.GetLogLevel() == "DEBUG" {
		log.Println("main.setupWireGuard -> called")
	} else if !*setupWg {
		log.Println("main.setupWireGuard -> Skipping wireguard interface setup")
		return
	}
	serverInterface := config.GetEnv("BACKEND_WG_INTERFACE", "wg0")
	serverPort, err := strconv.Atoi(config.GetEnv("BACKEND_WG_PORT", "51820"))
	if err != nil {
		log.Fatalf("main.setupWireGuard -> invalid port provided: %v", err)
	}

	serverIP := config.GetEnv("BACKEND_WG_IP", "10.0.0.1")
	networkMask := config.GetEnv("WG_NETWORK_MASK", "/24")

	if err := wgutil.InitWireGuardInterface(serverInterface, serverPort, net.ParseIP(serverIP), networkMask); err != nil {
		log.Fatalf("main.setupWireGuard -> failed to set up WireGuard network: %v", err)
	}
	log.Println("main.setupWireGuard -> WireGuard setup complete")
}

func startServer() {
	if config.GetLogLevel() == "DEBUG" {
		log.Println("main.startServer -> called")
	}
	port := ":" + config.GetEnv("PORT", "8080")
	server := &http.Server{Addr: port, Handler: routes.SetupRoutes()}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("main.startServer -> failed to listen and serve: %v", err)
	}

	defer db.CloseDatabaseConnection()
	log.Println("main.startServer -> server shutdown gracefully")
}
