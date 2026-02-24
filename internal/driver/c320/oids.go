package c320

import (
	"strconv"

	"github.com/ardani/snmp-zte/internal/driver"
)

// Konstanta OID untuk ZTE C320 - UPDATED dari hasil riset 2026-02-24
// Community strings: public (RO), globalrw (RW)
const (
	// === BASE OID ===
	BaseOID1 = ".1.3.6.1.4.1.3902.1082" // Legacy (working)
	BaseOID2 = ".1.3.6.1.4.1.3902.1012" // zxGponService - MAIN TREE
	BaseOID3 = ".1.3.6.1.4.1.3902.1015" // zxAn - Traffic stats

	// === OID BARU dari Riset (.1012.3.28) ===
	// ONU Device Management - Base: .1012.3.28.1.1.{field}.{oltId}.{onuId}
	OnuMgmtBase        = ".28.1.1" // Under BaseOID2.3
	OnuTypeNameOID     = ".28.1.1.1"  // TypeName (ZTE-F609V2.0)
	OnuNameOID         = ".28.1.1.2"  // Name ✅ WRITEABLE
	OnuDescriptionNewOID = ".28.1.1.3" // Description ✅ WRITEABLE
	OnuRegisterIdOID   = ".28.1.1.4"  // RegisterId (CZTE)
	OnuSerialNumberNewOID = ".28.1.1.5" // SerialNumber (Hex)
	OnuPwModeOID       = ".28.1.1.6"  // PwMode
	OnuPasswordOID     = ".28.1.1.7"  // Password
	OnuTargetStateOID  = ".28.1.1.8"  // TargetState (1=offline, 2=online)
	OnuRowStatusOID    = ".28.1.1.9"  // RowStatus ✅ WRITEABLE (create/delete)
	OnuVportModeOID    = ".28.1.1.10" // VportMode
	OnuIsAutoUpdateOID = ".28.1.1.11" // IsAutoUpdate
	OnuRegModeOID      = ".28.1.1.12" // RegMode

	// === DISTANCE (.1012.3.11.4.1) ===
	DistanceBase       = ".11.4.1" // Under BaseOID2.3
	OnuEQDOID          = ".11.4.1.1" // Equalized Delay
	OnuDistanceOID     = ".11.4.1.2" // Distance in meters

	// === FEC CONFIG (.1012.3.11.3.1.1) ===
	FecConfigOID = ".11.3.1.1" // FEC Status (1=enabled)

	// === PON CONFIG (.1012.3.12.7.1) ===
	PonConfigBase = ".12.7.1" // PON Port config

	// === BANDWIDTH PROFILES (.1012.3.26) ===
	ProfileBase           = ".26.1.1" // Under BaseOID2.3
	ProfileNameOID        = ".26.1.1.2"  // Profile Name
	ProfileFixedBWOID     = ".26.1.1.3"  // Fixed Bandwidth (kbps)
	ProfileAssuredBWOID   = ".26.1.1.4"  // Assured Bandwidth (kbps)
	ProfileMaxBWOID       = ".26.1.1.5"  // Max Bandwidth (kbps)

	// Traffic Profile Table 2
	TrafficProfileNameOID     = ".26.2.1.2" // Traffic Profile Name
	TrafficProfileFixedBWOID  = ".26.2.1.3" // Fixed BW
	TrafficProfileAssuredBWOID = ".26.2.1.4" // Assured BW

	// === TRAFFIC STATISTICS (.1015.1010.5.4.1) ===
	TrafficStatsBase       = ".1010.5.4.1" // Under BaseOID3
	PonRxOctetsOID         = ".1010.5.4.1.2"  // RX Bytes (Counter64)
	PonRxPktsOID           = ".1010.5.4.1.3"  // RX Packets (Counter64)
	PonRxPktsDiscardOID    = ".1010.5.4.1.4"  // RX Discards
	PonRxPktsErrOID        = ".1010.5.4.1.5"  // RX Errors
	PonRxCRCAlignErrorsOID = ".1010.5.4.1.6"  // CRC Errors
	PonTxOctetsOID         = ".1010.5.4.1.17" // TX Bytes (Counter64)
	PonTxPktsOID           = ".1010.5.4.1.18" // TX Packets (Counter64)

	// === VLAN (.2.1.17.7.1.4.3) - Standard IF-MIB ===
	VlanNameBase = ".1.3.6.1.2.1.17.7.1.4.3.1.1" // VLAN Names
	VlanPVIDBase = ".1.3.6.1.2.1.17.7.1.4.5.1.1" // Port PVID

	// === LEGACY OID (masih bekerja) ===
	// Prefix OID untuk ONU - Legacy (1082.500.*)
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
	CardRealTypePrefix = ".2.1.1.3.1.4.1.1"  // Actual Type (String)
	CardStatusPrefix   = ".2.1.1.3.1.5.1.1"
	CardPortCountPrefix = ".2.1.1.3.1.7.1.1"
	CardCpuLoadPrefix  = ".2.1.1.3.1.9.1.1"
	CardMemUsagePrefix = ".2.1.1.3.1.11.1.1"
	CardSoftVerPrefix  = ".2.1.2.2.1.4.1.1"  // Software Version

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
	TempSystemOID = ".1.3.6.1.4.1.3902.1015.2.1.3.10.10.11.0"
	TempCPUOID    = ".1.3.6.1.4.1.3902.1015.2.1.3.10.10.12.0"

	// === INDEX CONSTANTS ===
	// Konstanta ID Board-PON
	Board1OnuIDBase   = 285278464
	Board1OnuTypeBase = 268500992
	Board2OnuIDBase   = 285278720
	Board2OnuTypeBase = 268566528

	// Base Index PON (untuk statistik port PON)
	// Formula: (1 << 28) | (0 << 24) | (board << 16) | (pon << 8)
	// Board 1, PON 1 = 268501248
	// Board 1, PON 2 = 268501504
	// Board 2, PON 1 = 268566784
	Board1PonIndexBase = 268501248
	Board2PonIndexBase = 268566784

	// Nilai Increment
	OnuIDIncrement   = 1
	OnuTypeIncrement = 256

	// Batasan Maksimal
	MaxBoards       = 4
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

// CalculateOltID menghitung OLT ID dari board dan PON
// Formula: (1 << 28) | (0 << 24) | (board << 16) | (pon << 8)
func CalculateOltID(boardID, ponID int) int {
	return (1 << 28) | (0 << 24) | (boardID << 16) | (ponID << 8)
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

// GetOnuMgmtOID mengembalikan OID untuk ONU Management (new format)
// Base: .1012.3.28.1.1.{field}.{oltId}.{onuId}
func GetOnuMgmtOID(field int, boardID, ponID, onuID int) string {
	oltID := CalculateOltID(boardID, ponID)
	return BaseOID2 + ".3" + strconv.Itoa(field) + "." + strconv.Itoa(oltID) + "." + strconv.Itoa(onuID)
}

// GetDistanceOID mengembalikan OID untuk Distance measurement
// Base: .1012.3.11.4.1.{field}.{oltId}.{onuId}
func GetDistanceOID(field int, boardID, ponID, onuID int) string {
	oltID := CalculateOltID(boardID, ponID)
	return BaseOID2 + ".3.11.4.1." + strconv.Itoa(field) + "." + strconv.Itoa(oltID) + "." + strconv.Itoa(onuID)
}

// GetTrafficStatsOID mengembalikan OID untuk PON port traffic stats
// Base: .1015.1010.5.4.1.{field}.{oltId}
func GetTrafficStatsOID(field int, boardID, ponID int) string {
	oltID := CalculateOltID(boardID, ponID)
	return BaseOID3 + ".1010.5.4.1." + strconv.Itoa(field) + "." + strconv.Itoa(oltID)
}

// GetProfileOID mengembalikan OID untuk Bandwidth Profile
// Base: .1012.3.26.{table}.1.{field}.{profileIndex}
func GetProfileOID(table, field int, profileIndex int) string {
	return BaseOID2 + ".3.26." + strconv.Itoa(table) + ".1." + strconv.Itoa(field) + "." + strconv.Itoa(profileIndex)
}

// GetVLANNameOID mengembalikan OID untuk VLAN name
// Base: .1.3.6.1.2.1.17.7.1.4.3.1.1.{vlanId}
func GetVLANNameOID(vlanID int) string {
	return VlanNameBase + "." + strconv.Itoa(vlanID)
}
