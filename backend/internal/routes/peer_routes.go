package routes

import (
	"elysium-backend/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func PeerRoutes(mux *mux.Router) {
	mux.HandleFunc("/peers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.PostPeerHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/peers/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetPeerHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
