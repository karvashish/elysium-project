package repositories

import (
	"time"
)

type Peer struct {
	ID         int64
	PublicKey  string
	AssignedIP string
	Status     string
	IsGateway  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
