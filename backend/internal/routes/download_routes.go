package routes

import (
	"elysium-backend/internal/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func DownloadRoutes(router *mux.Router) {
	log.Println("routes.DownloadRoutes -> called")

	router.HandleFunc("/downloads/{uniqueID}/{filename}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("------------------------------------------------------------------------------")
		log.Println("routes.DownloadRoutes -> handling request for /downloads/{uniqueID}/{filename}")

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		vars := mux.Vars(r)
		uniqueID := vars["uniqueID"]
		filename := vars["filename"]

		handlers.DownloadHandler(w, r, uniqueID, filename)
	})
}
