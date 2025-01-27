package handlers

import (
  "elysium-backend/internal/models"
  "elysium-backend/internal/services"
  "encoding/json"
  "log"
  "net/http"
  "time"

  "github.com/google/uuid"
  "github.com/gorilla/mux"
)

func GetAllPeersHandler(w http.ResponseWriter, r *http.Request) {
  log.Println("handlers.GetAllPeersHandler  -> Processing request from", r.RemoteAddr)

  res, err := services.GetAllPeer()
  if err != nil {
    http.Error(w, "Internal server error", http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")

  if err := json.NewEncoder(w).Encode(res); err != nil {
    http.Error(w, "Failed to encode response", http.StatusInternalServerError)
  }

}

func GetPeerHandler(w http.ResponseWriter, r *http.Request) {
  log.Println("handlers.GetPeerHandler  -> Processing request from", r.RemoteAddr)

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
  log.Println("handlers.PostPeerHandler -> Processing request from", r.RemoteAddr)

  var peer_request *models.Peer_Request

  decoder := json.NewDecoder(r.Body)

  if err := decoder.Decode(&peer_request); err != nil {
    http.Error(w, "Invalid Request", http.StatusBadRequest)
    return
  }

  if peer_request.PublicKey == nil || *peer_request.PublicKey == "" {
    log.Println("handlers.PostPeerHandler -> PublicKey was not provided by", r.RemoteAddr)
    empty_str := ""
    peer_request.PublicKey = &empty_str
  } else {
    log.Println("handlers.PostPeerHandler -> Received PublicKey from", r.RemoteAddr)
  }

  new_peer := models.Peer{
    PublicKey:  *peer_request.PublicKey,
    AssignedIP: nil,
    Status:     "pending",
    IsGateway:  false,
    CreatedOn:  time.Now().UTC(),
  }

  log.Println("handlers.PostPeerHandler -> requesting new IP")
  if err := services.AssignNewIP(&new_peer); err != nil {
    http.Error(w, "Unable to assign IP", http.StatusInternalServerError)
  }

  exePath, err := services.CompileClient(*peer_request.PublicKey, peer_request.OSArch, new_peer.AssignedIP)
  if err != nil {
    http.Error(w, "Compilation failed: "+err.Error(), http.StatusInternalServerError)
    return
  }

  if err := services.InsertPeer(&new_peer); err != nil {
    http.Error(w, "Error creating peer", http.StatusInternalServerError)
    return
  }

  relativeDownloadLink := "/downloads/" + exePath

  response := map[string]string{"download_link": relativeDownloadLink}
  json.NewEncoder(w).Encode(response)
}
