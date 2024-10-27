package main

import (
	"flag"
	"log"
	"net/http"

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

	err = wgutil.InitWireGuardInterface()
	if err != nil {
		log.Fatalf("Failed setup wireguard network: %v", err)
	}

	port := ":" + config.GetEnv("PORT", "8080")

	server := &http.Server{Addr: port, Handler: routes.SetupRoutes()}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to listen and serve: %v", err)
	}

	db.CloseDatabaseConnection()
}
