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

func GetPeer(peerID *uuid.UUID) error {
	test, err := repositories.GetPeer(*peerID)

	log.Println(test.AssignedIP)
	if err != nil {
		log.Printf("services.InsertPeer -> Error inserting peer : %v", err)
		return err
	}
	return nil
}
