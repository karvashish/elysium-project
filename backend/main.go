package main

import (
	"elysium-backend/config"
	"elysium-backend/internal/routes"
	"elysium-backend/pkg/db"
	"log"
	"net/http"
)

func main() {
	config.LoadEnv()
	db.DBPool = db.InitializeDatabaseConnection()

	migrationDir := config.GetEnv("MIGRATION_PATH", "migrations")
	err := db.RunMigrations(migrationDir)
	if err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}

	port := ":" + config.GetEnv("PORT", "8080")

	server := &http.Server{Addr: port, Handler: routes.SetupRoutes()}

	server.ListenAndServe()

	db.CloseDatabaseConnection()
}
