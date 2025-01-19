package repositories

import (
	"context"
	"elysium-backend/config"
	"elysium-backend/internal/models"
	"elysium-backend/pkg/db"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

func InsertPeer(peer *models.Peer) error {
	if config.GetLogLevel() == "DEBUG" {
		log.Println("repositories.InsertPeer -> called")
	}

	query := `
		INSERT INTO peers (public_key, assigned_ip, status, is_gateway, created_on)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	ctx := context.Background()

	err := db.DBPool.QueryRowContext(ctx, query, peer.PublicKey, peer.AssignedIP, peer.Status, peer.IsGateway, peer.CreatedOn).Scan(&peer.ID)
	if err != nil {
		log.Println("repositories.InsertPeer -> Error inserting peer:", err)
		return err
	}
	log.Println("repositories.InsertPeer -> peer:", peer.ID)
	return nil
}

func GetPeer(id uuid.UUID) (*models.Peer, error) {
	if config.GetLogLevel() == "DEBUG" {
		log.Println("repositories.GetPeer -> called")
	}

	query := `
		SELECT id, public_key, assigned_ip, status, is_gateway, metadata, created_on
		FROM peers
		WHERE id = $1
	`
	ctx := context.Background()
	peer := &models.Peer{}

	var createdOnStr string

	row := db.DBPool.QueryRowContext(ctx, query, id)
	err := row.Scan(&peer.ID, &peer.PublicKey, &peer.AssignedIP, &peer.Status, &peer.IsGateway, &peer.Metadata, &createdOnStr)
	if err != nil {
		log.Println("repositories.GetPeer -> Error retrieving peer:", err)
		return nil, err
	}

	if createdOnStr == "" {
		log.Println("repositories.GetPeer -> created_on is empty or null")
		return nil, fmt.Errorf("created_on is empty or null")
	}

	peer.CreatedOn, err = time.Parse(time.RFC3339, strings.Replace(createdOnStr, " ", "T", 1))
	if err != nil {
		log.Println("repositories.GetPeer -> Error parsing created_on:", err)
		return nil, err
	}

	return peer, nil
}
