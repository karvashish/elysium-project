package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Peer struct {
	ID         *uuid.UUID
	PublicKey  string
	AssignedIP string
	Status     bool
	IsGateway  bool
	CreatedOn  time.Time
}
