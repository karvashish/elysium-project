package models

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type Peer struct {
	ID         *uuid.UUID              `json:"id" db:"id"`
	PublicKey  string                  `json:"public_key" db:"public_key"`
	AssignedIP net.IP                  `json:"assigned_ip" db:"assigned_ip"`
	Status     string                  `json:"status" db:"status"`
	IsGateway  bool                    `json:"is_gateway" db:"is_gateway"`
	Metadata   *map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedOn  time.Time               `json:"created_on" db:"created_on"`
}
