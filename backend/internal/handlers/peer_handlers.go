package handlers

import (
	"elysium-backend/internal/models"
	"elysium-backend/internal/services"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetPeerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	res, err := services.GetPeer(&id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func PostPeerHandler(w http.ResponseWriter, r *http.Request) {

	var peer_request *models.Peer_Request

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&peer_request); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	exePath, err := services.CompileClient(*peer_request.PublicKey, peer_request.OSArch)
	if err != nil {
		http.Error(w, fmt.Sprintf("Compilation failed: %s", err), http.StatusInternalServerError)
		return
	}

	relativeDownloadLink := fmt.Sprintf("/downloads/%s", exePath)

	response := map[string]string{"download_link": relativeDownloadLink}
	json.NewEncoder(w).Encode(response)
}
