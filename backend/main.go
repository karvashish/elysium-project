package main

import (
	"elysium-backend/config"
	"elysium-backend/internal/routes"
	"net/http"
)

func main() {
	config.LoadEnv()
	config.DBPool = config.InitializeDatabaseConnection()
	server := &http.Server{Addr: ":8080", Handler: routes.SetupRoutes()}

	server.ListenAndServe()

	config.CloseDatabaseConnection()
}
