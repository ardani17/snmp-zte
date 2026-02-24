package c320

import (
	"strconv"

	"github.com/ardani/snmp-zte/internal/driver"
)

// Konstanta OID untuk ZTE C320
const (
	BaseOID1 = ".1.3.6.1.4.1.3902.1082"
	BaseOID2 = ".1.3.6.1.4.1.3902.1012"
	BaseOID3 = ".1.3.6.1.4.1.3902.1015"

	// Prefix OID untuk ONU
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

	// OID untuk Board/Card (di bawah BaseOID3)
	CardTypePrefix     = ".2.1.1.3.1.2"       // Configured Type (Integer)
	CardRealTypePrefix = ".2.1.1.3.1.4.1.1"  // Actual Type (String) - Use this
	CardStatusPrefix   = ".2.1.1.3.1.5.1.1"
	CardPortCountPrefix = ".2.1.1.3.1.7.1.1"
	CardCpuLoadPrefix  = ".2.1.1.3.1.9.1.1"
	CardMemUsagePrefix = ".2.1.1.3.1.11.1.1"
	CardSoftVerPrefix  = ".2.1.2.2.1.4.1.1"  // Software Version

	// OID untuk Port PON
	PonTxPowerPrefix = ".1010.11.1.1.5"
	PonRxPowerPrefix = ".1010.11.2.1.2"

	// OID untuk Bandwidth ONU (SLA) - Corrected from ANALISIS.md
	OnuAssuredUpstreamPrefix   = ".1010.1.7.7.1.2"   // Assured upstream (kbps)
	OnuAssuredDownstreamPrefix = ".1010.1.7.8.1.1"   // Assured downstream (kbps)
	OnuMaxUpstreamPrefix       = ".1010.1.7.7.1.3"   // Max upstream (kbps)
	OnuMaxDownstreamPrefix     = ".1010.1.7.8.1.2"   // Max downstream (kbps)

	// OID untuk Statistik PON Port
	PonPortRxBytesOID   = ".1010.1.9.1.8.1.3"   // RX bytes per PON port
	PonPortTxBytesOID   = ".1010.1.9.1.8.1.4"   // TX bytes per PON port
	PonPortRxPacketsOID = ".1010.1.9.1.8.1.5"   // RX packets per PON port
	PonPortTxPacketsOID = ".1010.1.9.1.8.1.6"   // TX packets per PON port
	PonPortStatusOID    = ".1010.11.1.1.6"      // Admin status per PON port

	// OID untuk Error Counter ONU
	OnuCrcErrorOID        = ".1010.1.9.1.4.1.3"   // CRC errors
	OnuFecErrorOID        = ".1010.1.9.1.4.1.4"   // FEC errors
	OnuDroppedFramesOID   = ".1010.1.9.1.4.1.5"   // Dropped frames
	OnuLostPacketsOID     = ".1010.1.9.1.4.1.6"   // Lost packets

	// OID untuk Voltage/Power Supply
	VoltageSystemOID = ".1010.11.1.1.3"   // System voltage
	VoltageCPUOID    = ".1010.11.1.2.3"   // CPU voltage

	// OID System (RFC 1213)
	SysDescrOID   = ".1.3.6.1.2.1.1.1.0"
	SysNameOID    = ".1.3.6.1.2.1.1.5.0"
	SysUptimeOID  = ".1.3.6.1.2.1.1.3.0"
	SysContactOID = ".1.3.6.1.2.1.1.4.0"
	SysLocationOID = ".1.3.6.1.2.1.1.6.0"

	// OID untuk Fan
	FanTableOID       = ".1.3.6.1.4.1.3902.1015.2.1.3.10.10.10"
	FanSpeedLevelOID  = ".1.3.6.1.4.1.3902.1015.2.1.3.10.10.10.1.3"
	FanStatusOID      = ".1.3.6.1.4.1.3902.1015.2.1.3.10.10.10.1.5"
	FanPresentOID     = ".1.3.6.1.4.1.3902.1015.2.1.3.10.10.10.1.6"

	// OID untuk Temperature
	TempSystemOID = ".1.3.6.1.4.1.3902.1015.2.1.3.10.10.11.0" // Sistem/Ambient
	TempCPUOID    = ".1.3.6.1.4.1.3902.1015.2.1.3.10.10.12.0" // CPU/Board

	// Konstanta ID Board-PON
	Board1OnuIDBase   = 285278464
	Board1OnuTypeBase = 268500992
	Board2OnuIDBase   = 285278720
	Board2OnuTypeBase = 268566528

	// Base Index PON (untuk statistik port PON)
	Board1PonIndexBase = 268501248
	Board2PonIndexBase = 268566784

	// Nilai Increment (Penambahan)
	OnuIDIncrement   = 1
	OnuTypeIncrement = 256

	// Batasan Maksimal
	MaxBoards       = 4  // Show all 4 slots including empty
	MaxPonPerBoard  = 16
	MaxOnuPerPon    = 128
)

// ModelInfo mengembalikan informasi model C320.
func ModelInfo() driver.ModelInfo {
	return driver.ModelInfo{
		Name:           "ZTE C320",
		Vendor:         "ZTE",
		MaxBoards:      MaxBoards,
		MaxPonPerBoard: MaxPonPerBoard,
		MaxOnuPerPon:   MaxOnuPerPon,
	}
}

// GenerateBoardPonOID membuat konfigurasi OID untuk kombinasi Board/PON tertentu.
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

// GetPonIndexBase mengembalikan base index PON untuk sebuah board.
func GetPonIndexBase(boardID int) int {
	if boardID == 1 {
		return Board1PonIndexBase
	}
	return Board2PonIndexBase
}
