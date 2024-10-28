package services

import (
	"elysium-backend/internal/models"
	"elysium-backend/internal/repositories"
	"log"
)

func InsertPeer(newPeer *models.Peer) error {

	if err := repositories.InsertPeer(newPeer); err != nil {
		log.Printf("services.InsertPeer -> Error inserting peer : %v", err)
		return err
	}

	return nil
}
