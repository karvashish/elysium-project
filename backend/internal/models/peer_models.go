package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Peer struct {
	ID         *uuid.UUID
	PublicKey  string
	AssignedIP string
	Status     string
	IsGateway  bool
	CreatedOn  time.Time
}
