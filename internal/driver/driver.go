package driver

import (
	"context"

	"github.com/ardani/snmp-zte/internal/model"
)

// Interface Driver mendefinisikan metode yang harus diimplementasikan oleh semua driver OLT
type Driver interface {
	// Metadata
	GetModelName() string
	GetModelInfo() ModelInfo

	// Operasi ONU
	GetONUList(ctx context.Context, boardID, ponID int) ([]model.ONUInfo, error)
	GetONUDetail(ctx context.Context, boardID, ponID, onuID int) (*model.ONUDetail, error)
	GetEmptySlots(ctx context.Context, boardID, ponID int) ([]model.ONUSlot, error)
	GetONUTraffic(ctx context.Context, boardID, ponID, onuID int) (*model.ONUTraffic, error)

	// Info OLT
	GetSystemInfo(ctx context.Context) (*SystemInfo, error)
	GetBoardInfo(ctx context.Context, boardID int) (*model.BoardInfo, error)
	GetAllBoards(ctx context.Context) ([]model.BoardInfo, error)
	GetInterfaceStats(ctx context.Context) ([]model.InterfaceStats, error)
	GetFanInfo(ctx context.Context) ([]map[string]interface{}, error)
	GetTemperatureInfo(ctx context.Context) (*model.TemperatureInfo, error)

	// Validasi
	ValidateBoardID(boardID int) bool
	ValidatePonID(ponID int) bool
	ValidateOnuID(onuID int) bool

	// Koneksi
	Connect() error
	Close() error
}

// ModelInfo merepresentasikan informasi model OLT
type ModelInfo struct {
	Name           string `json:"name"`
	Vendor         string `json:"vendor"`
	MaxBoards      int    `json:"max_boards"`
	MaxPonPerBoard int    `json:"max_pon_per_board"`
	MaxOnuPerPon   int    `json:"max_onu_per_pon"`
}

// SystemInfo merepresentasikan informasi sistem OLT
type SystemInfo struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	Uptime      string `json:"uptime"`
	Contact     string `json:"contact"`
	Location    string `json:"location"`
}

// BoardPonConfig merepresentasikan konfigurasi OID untuk kombinasi Board/PON tertentu
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
