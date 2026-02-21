package c320

import (
	"strconv"

	"github.com/ardani/snmp-zte/internal/driver"
)

// OID constants for ZTE C320
const (
	BaseOID1 = ".1.3.6.1.4.1.3902.1082"
	BaseOID2 = ".1.3.6.1.4.1.3902.1012"
	BaseOID3 = ".1.3.6.1.4.1.3902.1015"

	// OID Prefixes for ONU
	OnuIDNamePrefix              = ".500.10.2.3.3.1.2"
	OnuTypePrefix                = ".3.50.11.2.1.17"
	OnuSerialNumberPrefix        = ".500.10.2.3.3.1.18"
	OnuRxPowerPrefix             = ".500.20.2.2.2.1.10"
	OnuTxPowerPrefix             = ".3.50.12.1.1.14"
	OnuStatusIDPrefix            = ".500.10.2.3.8.1.4"
	OnuIPAddressPrefix           = ".3.50.16.1.1.10"
	OnuDescriptionPrefix         = ".500.10.2.3.3.1.3"
	OnuLastOnlineTimePrefix      = ".500.10.2.3.8.1.5"
	OnuLastOfflineTimePrefix     = ".500.10.2.3.8.1.6"
	OnuLastOfflineReasonPrefix   = ".500.10.2.3.8.1.7"
	OnuGponOpticalDistancePrefix = ".500.10.2.3.10.1.2"

	// Card/Board OIDs (discovered via SNMP walk)
	CardCpuLoadOID  = ".2.1.1.3.1.9.1.1"    // CPU load per card
	CardMemUsageOID = ".2.1.1.3.1.10.1.1"   // Memory usage per card

	// PON Port OIDs (discovered via SNMP walk)
	// Index: 268501248 for Board 1 (268500992 + 256)
	PonRxPowerOID = ".1010.11.2.1.2"  // RX power in hundredths
	PonTxPowerOID = ".1010.11.1.1.5"  // TX power in hundredths

	// Board-PON ID Constants
	Board1OnuIDBase   = 285278464
	Board1OnuTypeBase = 268500992
	Board2OnuIDBase   = 285278720
	Board2OnuTypeBase = 268566528

	// PON Index base (for PON port stats)
	Board1PonIndexBase = 268501248
	Board2PonIndexBase = 268566784

	// Increment values
	OnuIDIncrement   = 1
	OnuTypeIncrement = 256

	// Limits
	MaxBoards       = 2
	MaxPonPerBoard  = 16
	MaxOnuPerPon    = 128
)

// ModelInfo returns C320 model information
func ModelInfo() driver.ModelInfo {
	return driver.ModelInfo{
		Name:           "ZTE C320",
		Vendor:         "ZTE",
		MaxBoards:      MaxBoards,
		MaxPonPerBoard: MaxPonPerBoard,
		MaxOnuPerPon:   MaxOnuPerPon,
	}
}

// GenerateBoardPonOID generates OID configuration for a specific Board/PON
func GenerateBoardPonOID(boardID, ponID int) *driver.BoardPonConfig {
	var baseOnuID, baseOnuType int

	if boardID == 1 {
		baseOnuID = Board1OnuIDBase
		baseOnuType = Board1OnuTypeBase
	} else {
		baseOnuID = Board2OnuIDBase
		baseOnuType = Board2OnuTypeBase
	}

	onuIDSuffix := baseOnuID + (ponID * OnuIDIncrement)
	onuTypeSuffix := baseOnuType + (ponID * OnuTypeIncrement)

	return &driver.BoardPonConfig{
		OnuIDNameOID:              OnuIDNamePrefix + "." + strconv.Itoa(onuIDSuffix),
		OnuTypeOID:                OnuTypePrefix + "." + strconv.Itoa(onuTypeSuffix),
		OnuSerialNumberOID:        OnuSerialNumberPrefix + "." + strconv.Itoa(onuIDSuffix),
		OnuRxPowerOID:             OnuRxPowerPrefix + "." + strconv.Itoa(onuIDSuffix),
		OnuTxPowerOID:             OnuTxPowerPrefix + "." + strconv.Itoa(onuTypeSuffix),
		OnuStatusOID:              OnuStatusIDPrefix + "." + strconv.Itoa(onuIDSuffix),
		OnuIPAddressOID:           OnuIPAddressPrefix + "." + strconv.Itoa(onuTypeSuffix),
		OnuDescriptionOID:         OnuDescriptionPrefix + "." + strconv.Itoa(onuIDSuffix),
		OnuLastOnlineOID:          OnuLastOnlineTimePrefix + "." + strconv.Itoa(onuIDSuffix),
		OnuLastOfflineOID:         OnuLastOfflineTimePrefix + "." + strconv.Itoa(onuIDSuffix),
		OnuLastOfflineReasonOID:   OnuLastOfflineReasonPrefix + "." + strconv.Itoa(onuIDSuffix),
		OnuGponOpticalDistanceOID: OnuGponOpticalDistancePrefix + "." + strconv.Itoa(onuIDSuffix),
	}
}

// GetPonIndexBase returns the PON index base for a board
func GetPonIndexBase(boardID int) int {
	if boardID == 1 {
		return Board1PonIndexBase
	}
	return Board2PonIndexBase
}
