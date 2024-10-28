package routes

import (
	"elysium-backend/internal/handlers"
	"net/http"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.BaseHandler)

	PeerRoutes(mux)

	return mux
}
