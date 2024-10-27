package services

import (
	"elysium-backend/internal/repositories"
)

func InsertPeer(newPeer repositories.Peer) error {

	repositories.InsertPeer(&newPeer)

	return nil
}
