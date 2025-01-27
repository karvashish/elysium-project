package models

import (
  "errors"
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

type OSArch string

const (
  OSArchx86_64Linux  OSArch = "x86_64-unknown-linux-musl"
  OSArchAarch64Linux OSArch = "aarch64-unknown-linux-musl"
  OSArchWindows      OSArch = "x86_64-pc-windows-gnu"
)

func (o OSArch) Validate() error {
  switch o {
  case OSArchx86_64Linux, OSArchWindows, OSArchAarch64Linux:
    return nil
  default:
    return errors.New("invalid OS_Arch value")
  }
}

type Peer_Request struct {
  PublicKey *string `json:"public_key"`
  OSArch    OSArch  `json:"OS_Arch"`
}
