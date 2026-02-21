package driver

import (
	"context"

	"github.com/ardani/snmp-zte/internal/model"
)

// Driver interface defines methods that all OLT drivers must implement
type Driver interface {
	// Metadata
	GetModelName() string
	GetModelInfo() ModelInfo

	// ONU Operations
	GetONUList(ctx context.Context, boardID, ponID int) ([]model.ONUInfo, error)
	GetONUDetail(ctx context.Context, boardID, ponID, onuID int) (*model.ONUDetail, error)
	GetEmptySlots(ctx context.Context, boardID, ponID int) ([]model.ONUSlot, error)

	// OLT Info
	GetSystemInfo(ctx context.Context) (*SystemInfo, error)
	GetBoardInfo(ctx context.Context, boardID int) (*model.BoardInfo, error)
	GetPONInfo(ctx context.Context, boardID, ponID int) (*model.PONInfo, error)

	// Validation
	ValidateBoardID(boardID int) bool
	ValidatePonID(ponID int) bool
	ValidateOnuID(onuID int) bool

	// Connection
	Connect() error
	Close() error
}

// ModelInfo represents OLT model information
type ModelInfo struct {
	Name         string `json:"name"`
	Vendor       string `json:"vendor"`
	MaxBoards    int    `json:"max_boards"`
	MaxPonPerBoard int  `json:"max_pon_per_board"`
	MaxOnuPerPon int    `json:"max_onu_per_pon"`
}

// SystemInfo represents OLT system information
type SystemInfo struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	Uptime      string `json:"uptime"`
	Contact     string `json:"contact"`
	Location    string `json:"location"`
}

// OIDConfig represents OID configuration for a specific OLT model
type OIDConfig struct {
	BaseOID1 string
	BaseOID2 string
}

// BoardPonConfig represents OID configuration for a specific Board/PON combination
type BoardPonConfig struct {
	OnuIDNameOID              string
	OnuTypeOID                string
	OnuSerialNumberOID        string
	OnuRxPowerOID             string
	OnuTxPowerOID             string
	OnuStatusOID              string
	OnuIPAddressOID           string
	OnuDescriptionOID         string
	OnuLastOnlineOID          string
	OnuLastOfflineOID         string
	OnuLastOfflineReasonOID   string
	OnuGponOpticalDistanceOID string
}
