package routes

import (
	"elysium-backend/internal/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.BaseHandler)

	PeerRoutes(router)

	return router
}
