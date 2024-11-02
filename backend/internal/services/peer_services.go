package services

import (
	"elysium-backend/internal/models"
	"elysium-backend/internal/repositories"
	"log"

	"github.com/google/uuid"
)

func InsertPeer(newPeer *models.Peer) error {

	if err := repositories.InsertPeer(newPeer); err != nil {
		log.Printf("services.InsertPeer -> Error inserting peer : %v", err)
		return err
	}

	return nil
}

func GetPeer(peerID *uuid.UUID) (*models.Peer, error) {
	peer, err := repositories.GetPeer(*peerID)

	if err != nil {
		log.Printf("services.GetPeer -> Error retrieving peer : %v", err)
		return nil, err
	}
	return peer, nil
}
