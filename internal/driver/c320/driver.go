package c320

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ardani/snmp-zte/internal/driver"
	"github.com/ardani/snmp-zte/internal/model"
	"github.com/gosnmp/gosnmp"
)

// Driver implements driver.Driver for ZTE C320
type Driver struct {
	client    *gosnmp.GoSNMP
	snmpHost  string
	snmpPort  uint16
	community string
	connected bool
}

// New creates a new C320 driver
func New(host string, port uint16, community string) *Driver {
	return &Driver{
		snmpHost:  host,
		snmpPort:  port,
		community: community,
	}
}

// GetModelName returns the model name
func (d *Driver) GetModelName() string {
	return "C320"
}

// GetModelInfo returns model information
func (d *Driver) GetModelInfo() driver.ModelInfo {
	return ModelInfo()
}

// Connect establishes SNMP connection
func (d *Driver) Connect() error {
	d.client = &gosnmp.GoSNMP{
		Target:    d.snmpHost,
		Port:      d.snmpPort,
		Community: d.community,
		Version:   gosnmp.Version2c,
		Timeout:   5 * time.Second,
		Retries:   2,
		MaxOids:   60,
	}

	if err := d.client.Connect(); err != nil {
		return fmt.Errorf("SNMP connect failed: %w", err)
	}

	d.connected = true
	return nil
}

// Close closes the SNMP connection
func (d *Driver) Close() error {
	if d.client != nil && d.client.Conn != nil {
		d.connected = false
		return d.client.Conn.Close()
	}
	return nil
}

// ValidateBoardID validates board ID
func (d *Driver) ValidateBoardID(boardID int) bool {
	return boardID >= 1 && boardID <= MaxBoards
}

// ValidatePonID validates PON ID
func (d *Driver) ValidatePonID(ponID int) bool {
	return ponID >= 1 && ponID <= MaxPonPerBoard
}

// ValidateOnuID validates ONU ID
func (d *Driver) ValidateOnuID(onuID int) bool {
	return onuID >= 1 && onuID <= MaxOnuPerPon
}

// GetONUList returns list of ONUs for a Board/PON
func (d *Driver) GetONUList(ctx context.Context, boardID, ponID int) ([]model.ONUInfo, error) {
	if !d.connected {
		if err := d.Connect(); err != nil {
			return nil, err
		}
	}

	cfg := GenerateBoardPonOID(boardID, ponID)
	
	var onuList []model.ONUInfo
	onuMap := make(map[int]*model.ONUInfo)

	// Walk to get all ONU IDs and Names
	oid := BaseOID1 + cfg.OnuIDNameOID
	err := d.client.Walk(oid, func(pdu gosnmp.SnmpPDU) error {
		onuID := extractOnuID(pdu.Name)
		if onuID == 0 {
			return nil
		}

		info := &model.ONUInfo{
			Board: boardID,
			PON:   ponID,
			ID:    onuID,
			Name:  extractString(pdu.Value),
		}
		onuMap[onuID] = info
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("SNMP walk failed: %w", err)
	}

	// Get additional info for each ONU
	for onuID, info := range onuMap {
		onuIDStr := strconv.Itoa(onuID)
		
		// Get ONU Type
		if val, err := d.snmpGet(BaseOID2 + cfg.OnuTypeOID + "." + onuIDStr); err == nil {
			info.Type = extractString(val)
		}

		// Get Serial Number
		if val, err := d.snmpGet(BaseOID1 + cfg.OnuSerialNumberOID + "." + onuIDStr); err == nil {
			info.SerialNumber = extractSerialNumber(val)
		}

		// Get RX Power
		if val, err := d.snmpGet(BaseOID1 + cfg.OnuRxPowerOID + "." + onuIDStr + ".1"); err == nil {
			info.RXPower = convertPower(val)
		}

		// Get Status
		if val, err := d.snmpGet(BaseOID1 + cfg.OnuStatusOID + "." + onuIDStr); err == nil {
			info.Status = convertStatus(val)
		}

		onuList = append(onuList, *info)
	}

	return onuList, nil
}

// GetONUDetail returns detailed information for a single ONU
func (d *Driver) GetONUDetail(ctx context.Context, boardID, ponID, onuID int) (*model.ONUDetail, error) {
	if !d.connected {
		if err := d.Connect(); err != nil {
			return nil, err
		}
	}

	cfg := GenerateBoardPonOID(boardID, ponID)
	onuIDStr := strconv.Itoa(onuID)

	detail := &model.ONUDetail{
		ONUInfo: model.ONUInfo{
			Board: boardID,
			PON:   ponID,
			ID:    onuID,
		},
	}

	// Get Name
	if val, err := d.snmpGet(BaseOID1 + cfg.OnuIDNameOID + "." + onuIDStr); err == nil {
		detail.Name = extractString(val)
	}

	// Get Type
	if val, err := d.snmpGet(BaseOID2 + cfg.OnuTypeOID + "." + onuIDStr); err == nil {
		detail.Type = extractString(val)
	}

	// Get Serial Number
	if val, err := d.snmpGet(BaseOID1 + cfg.OnuSerialNumberOID + "." + onuIDStr); err == nil {
		detail.SerialNumber = extractSerialNumber(val)
	}

	// Get RX Power
	if val, err := d.snmpGet(BaseOID1 + cfg.OnuRxPowerOID + "." + onuIDStr + ".1"); err == nil {
		detail.RXPower = convertPower(val)
	}

	// Get TX Power
	if val, err := d.snmpGet(BaseOID2 + cfg.OnuTxPowerOID + "." + onuIDStr + ".1"); err == nil {
		detail.TXPower = convertPower(val)
	}

	// Get Status
	if val, err := d.snmpGet(BaseOID1 + cfg.OnuStatusOID + "." + onuIDStr); err == nil {
		detail.Status = convertStatus(val)
	}

	// Get IP Address
	if val, err := d.snmpGet(BaseOID2 + cfg.OnuIPAddressOID + "." + onuIDStr + ".1"); err == nil {
		detail.IPAddress = extractString(val)
	}

	// Get Description
	if val, err := d.snmpGet(BaseOID1 + cfg.OnuDescriptionOID + "." + onuIDStr); err == nil {
		detail.Description = extractString(val)
	}

	// Get Last Online
	if val, err := d.snmpGet(BaseOID1 + cfg.OnuLastOnlineOID + "." + onuIDStr); err == nil {
		detail.LastOnline = convertDateTime(val)
	}

	// Get Last Offline
	if val, err := d.snmpGet(BaseOID1 + cfg.OnuLastOfflineOID + "." + onuIDStr); err == nil {
		detail.LastOffline = convertDateTime(val)
	}

	// Get Offline Reason
	if val, err := d.snmpGet(BaseOID1 + cfg.OnuLastOfflineReasonOID + "." + onuIDStr); err == nil {
		detail.OfflineReason = convertOfflineReason(val)
	}

	// Get Distance
	if val, err := d.snmpGet(BaseOID1 + cfg.OnuGponOpticalDistanceOID + "." + onuIDStr); err == nil {
		detail.Distance = fmt.Sprintf("%v", val)
	}

	// Calculate Uptime
	if detail.LastOnline != "" {
		detail.Uptime = calculateUptime(detail.LastOnline)
	}

	return detail, nil
}

// GetEmptySlots returns available ONU slots
func (d *Driver) GetEmptySlots(ctx context.Context, boardID, ponID int) ([]model.ONUSlot, error) {
	if !d.connected {
		if err := d.Connect(); err != nil {
			return nil, err
		}
	}

	cfg := GenerateBoardPonOID(boardID, ponID)
	
	// Track used ONU IDs
	usedIDs := make(map[int]bool)
	
	oid := BaseOID1 + cfg.OnuIDNameOID
	err := d.client.Walk(oid, func(pdu gosnmp.SnmpPDU) error {
		onuID := extractOnuID(pdu.Name)
		if onuID > 0 {
			usedIDs[onuID] = true
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("SNMP walk failed: %w", err)
	}

	// Find empty slots
	var emptySlots []model.ONUSlot
	for i := 1; i <= MaxOnuPerPon; i++ {
		if !usedIDs[i] {
			emptySlots = append(emptySlots, model.ONUSlot{
				Board: boardID,
				PON:   ponID,
				ONUID: i,
			})
		}
	}

	return emptySlots, nil
}

// GetSystemInfo returns OLT system information
func (d *Driver) GetSystemInfo(ctx context.Context) (*driver.SystemInfo, error) {
	if !d.connected {
		if err := d.Connect(); err != nil {
			return nil, err
		}
	}

	info := &driver.SystemInfo{}

	// System Description
	if val, err := d.snmpGet("1.3.6.1.2.1.1.1.0"); err == nil {
		info.Description = extractString(val)
	}

	// System Name
	if val, err := d.snmpGet("1.3.6.1.2.1.1.5.0"); err == nil {
		info.Name = extractString(val)
	}

	// System Uptime
	if val, err := d.snmpGet("1.3.6.1.2.1.1.3.0"); err == nil {
		info.Uptime = fmt.Sprintf("%v", val)
	}

	// System Contact
	if val, err := d.snmpGet("1.3.6.1.2.1.1.4.0"); err == nil {
		info.Contact = extractString(val)
	}

	// System Location
	if val, err := d.snmpGet("1.3.6.1.2.1.1.6.0"); err == nil {
		info.Location = extractString(val)
	}

	return info, nil
}

// GetBoardInfo returns board information
func (d *Driver) GetBoardInfo(ctx context.Context, boardID int) (*model.BoardInfo, error) {
	// TODO: Implement board info retrieval
	return &model.BoardInfo{BoardID: boardID}, nil
}

// GetPONInfo returns PON port information
func (d *Driver) GetPONInfo(ctx context.Context, boardID, ponID int) (*model.PONInfo, error) {
	// TODO: Implement PON info retrieval
	return &model.PONInfo{BoardID: boardID, PonID: ponID}, nil
}

// snmpGet performs an SNMP GET request
func (d *Driver) snmpGet(oid string) (interface{}, error) {
	result, err := d.client.Get([]string{oid})
	if err != nil {
		return nil, err
	}
	if len(result.Variables) == 0 {
		return nil, fmt.Errorf("no result for OID: %s", oid)
	}
	return result.Variables[0].Value, nil
}

// Helper functions

func extractOnuID(oid string) int {
	// Extract last component from OID
	// e.g., ".1.3.6.1.4.1.3902.1082.500.10.2.3.3.1.2.285278465.5" -> 5
	parts := splitOID(oid)
	if len(parts) < 2 {
		return 0
	}
	lastPart := parts[len(parts)-1]
	id, _ := strconv.Atoi(lastPart)
	return id
}

func splitOID(oid string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(oid); i++ {
		if oid[i] == '.' {
			if i > start {
				parts = append(parts, oid[start:i])
			}
			start = i + 1
		}
	}
	if start < len(oid) {
		parts = append(parts, oid[start:])
	}
	return parts
}

func extractString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func extractSerialNumber(val interface{}) string {
	s := extractString(val)
	// Remove "1," prefix if present
	if len(s) > 2 && s[:2] == "1," {
		return s[2:]
	}
	return s
}

func convertPower(val interface{}) string {
	intVal, ok := val.(int)
	if !ok {
		return "0.00"
	}
	result := float64(intVal)*0.002 - 30.0
	return fmt.Sprintf("%.2f", result)
}

func convertStatus(val interface{}) string {
	intVal, ok := val.(int)
	if !ok {
		return "Unknown"
	}
	return model.ONUStatus(intVal).String()
}

func convertOfflineReason(val interface{}) string {
	intVal, ok := val.(int)
	if !ok {
		return "Unknown"
	}
	return model.OfflineReason(intVal).String()
}

func convertDateTime(val interface{}) string {
	bytes, ok := val.([]byte)
	if !ok || len(bytes) != 8 {
		return ""
	}
	// Parse 8-byte datetime format
	// Year(2), Month(1), Day(1), Hour(1), Min(1), Sec(1), Reserved(2)
	year := int(bytes[0])<<8 | int(bytes[1])
	month := int(bytes[2])
	day := int(bytes[3])
	hour := int(bytes[4])
	minute := int(bytes[5])
	second := int(bytes[6])
	
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, minute, second)
}

func calculateUptime(lastOnline string) string {
	// Parse lastOnline and calculate duration
	// TODO: Implement proper uptime calculation
	return lastOnline
}

// Ensure Driver implements driver.Driver interface
var _ driver.Driver = (*Driver)(nil)
