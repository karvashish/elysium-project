package main

import (
	"elysium-backend/config"
	"elysium-backend/internal/routes"
	"log"
	"net/http"
)

func main() {
	config.LoadEnv()
	config.DBPool = config.InitializeDatabaseConnection()

	migrationDir := config.GetEnv("MIGRATION_PATH", "migrations")
	err := config.RunMigrations(migrationDir)
	if err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}

	server := &http.Server{Addr: ":8080", Handler: routes.SetupRoutes()}

	server.ListenAndServe()

	config.CloseDatabaseConnection()
}
