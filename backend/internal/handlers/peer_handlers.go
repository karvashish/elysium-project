package handlers

import (
	"elysium-backend/internal/models"
	"elysium-backend/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
	fmt.Printf("Processing request from %s\n", r.RemoteAddr)

	var peer_request *models.Peer_Request

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&peer_request); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if peer_request.PublicKey == nil || *peer_request.PublicKey == "" {
		fmt.Printf("PublicKey was not provided by %s\n", r.RemoteAddr)
		empty_str := ""
		peer_request.PublicKey = &empty_str
	} else {
		fmt.Printf("Received PublicKey from %s\n", r.RemoteAddr)
	}

	exePath, err := services.CompileClient(*peer_request.PublicKey, peer_request.OSArch)
	if err != nil {
		http.Error(w, fmt.Sprintf("Compilation failed: %s", err), http.StatusInternalServerError)
		return
	}

	new_peer := models.Peer{
		PublicKey:  *peer_request.PublicKey,
		AssignedIP: nil,
		Status:     "pending",
		IsGateway:  false,
		CreatedOn:  time.Now().UTC(),
	}

	if err := services.InsertPeer(&new_peer); err != nil {
		fmt.Printf("Error Inserting peer %s\n", err)
	}

	relativeDownloadLink := fmt.Sprintf("/downloads/%s", exePath)

	response := map[string]string{"download_link": relativeDownloadLink}
	json.NewEncoder(w).Encode(response)
}
