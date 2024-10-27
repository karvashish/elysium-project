package repositories

import (
	"context"
	"elysium-backend/pkg/db"
	"log"
	"time"

	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

type Peer struct {
	ID         *uuid.UUID
	PublicKey  string
	AssignedIP string
	Status     bool
	IsGateway  bool
	CreatedOn  time.Time
}

func InsertPeer(peer *Peer) error {
	query := `
		INSERT INTO peers (public_key, assigned_ip, status, is_gateway, created_on)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	ctx := context.Background()

	err := db.DBPool.QueryRow(ctx, query, peer.PublicKey, peer.AssignedIP, peer.Status, peer.IsGateway, peer.CreatedOn).Scan(&peer.ID)
	if err != nil {
		log.Printf("Error inserting peer: %v", err)
		return err
	}
	return nil
}
