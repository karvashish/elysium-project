package routes

import (
	"elysium-backend/internal/handlers"
	"net/http"
)

func PeerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/peers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetPeerHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
