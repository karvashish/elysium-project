package repositories

import (
	"context"
	"elysium-backend/internal/models"
	"elysium-backend/pkg/db"
	"log"

	"github.com/google/uuid"
)

func InsertPeer(peer *models.Peer) error {
	query := `
		INSERT INTO peers (public_key, assigned_ip, status, is_gateway, created_on)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	ctx := context.Background()

	err := db.DBPool.QueryRow(ctx, query, peer.PublicKey, peer.AssignedIP, peer.Status, peer.IsGateway, peer.CreatedOn).Scan(&peer.ID)
	if err != nil {
		log.Printf("repositories.InsertPeer ->Error inserting peer: %v", err)
		return err
	}
	return nil
}

func GetPeer(id uuid.UUID) (*models.Peer, error) {
	query := `
		SELECT id, public_key, assigned_ip, status, is_gateway, metadata, created_on
		FROM peers
		WHERE id = $1
	`
	ctx := context.Background()
	peer := &models.Peer{}

	row := db.DBPool.QueryRow(ctx, query, id)
	err := row.Scan(&peer.ID, &peer.PublicKey, &peer.AssignedIP, &peer.Status, &peer.IsGateway, &peer.Metadata, &peer.CreatedOn)
	if err != nil {
		log.Printf("repositories.GetPeer -> Error retrieving peer: %v", err)
		return nil, err
	}
	return peer, nil
}
