package routes

import (
	"elysium-backend/internal/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func PeerRoutes(mux *mux.Router) {
	log.Println("routes.PeerRoutes -> called")

	mux.HandleFunc("/peer", func(w http.ResponseWriter, r *http.Request) {
		log.Println("------------------------------------------------------------------------------")
		log.Println("routes.PeerRoutes -> handling request for /peer")
		if r.Method == http.MethodPost {
			handlers.PostPeerHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/peers/{id}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("------------------------------------------------------------------------------")
		log.Println("routes.PeerRoutes -> handling request for /peers/{id}")
		if r.Method == http.MethodGet {
			handlers.GetPeerHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
