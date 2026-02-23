package c320

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ardani/snmp-zte/internal/driver"
	"github.com/ardani/snmp-zte/internal/model"
	"github.com/gosnmp/gosnmp"
)

// Driver mengimplementasikan driver.Driver untuk ZTE C320.
type Driver struct {
	client    *gosnmp.GoSNMP
	snmpHost  string
	snmpPort  uint16
	community string
	connected bool
}

// New membuat instance driver C320 baru.
func New(host string, port uint16, community string) *Driver {
	return &Driver{
		snmpHost:  host,
		snmpPort:  port,
		community: community,
	}
}

// GetModelName mengembalikan nama model.
func (d *Driver) GetModelName() string {
	return "C320"
}

// GetModelInfo mengembalikan informasi model.
func (d *Driver) GetModelInfo() driver.ModelInfo {
	return ModelInfo()
}

// Connect membuka koneksi SNMP ke OLT.
func (d *Driver) Connect() error {
	d.client = &gosnmp.GoSNMP{
		Target:    d.snmpHost,
		Port:      d.snmpPort,
		Community: d.community,
		Version:   gosnmp.Version2c,
		Timeout:   5 * time.Second, // Tunggu respon OLT maksimal 5 detik
		Retries:   2,               // Coba lagi 2 kali jika gagal
		MaxOids:   60,
	}

	if err := d.client.Connect(); err != nil {
		return fmt.Errorf("SNMP connect failed: %w", err)
	}

	d.connected = true
	return nil
}

// Close menutup koneksi SNMP.
func (d *Driver) Close() error {
	if d.client != nil && d.client.Conn != nil {
		d.connected = false
		return d.client.Conn.Close()
	}
	return nil
}

// ValidateBoardID memvalidasi ID board.
func (d *Driver) ValidateBoardID(boardID int) bool {
	return boardID >= 1 && boardID <= MaxBoards
}

// ValidatePonID memvalidasi ID PON.
func (d *Driver) ValidatePonID(ponID int) bool {
	return ponID >= 1 && ponID <= MaxPonPerBoard
}

// ValidateOnuID memvalidasi ID ONU.
func (d *Driver) ValidateOnuID(onuID int) bool {
	return onuID >= 1 && onuID <= MaxOnuPerPon
}

// GetONUList mengambil daftar ONU untuk Board/PON tertentu.
func (d *Driver) GetONUList(ctx context.Context, boardID, ponID int) ([]model.ONUInfo, error) {
	if !d.connected {
		if err := d.Connect(); err != nil {
			return nil, err
		}
	}

	cfg := GenerateBoardPonOID(boardID, ponID)
	
	var onuList []model.ONUInfo
	onuMap := make(map[int]*model.ONUInfo)

	// 1. SNMP Walk untuk mendapatkan daftar Nama & ID ONU yang aktif di port tersebut.
	// OID di sini spesifik untuk ZTE, digabung dengan ID Board & PON.
	oid := BaseOID1 + cfg.OnuIDNameOID
	err := d.client.Walk(oid, func(pdu gosnmp.SnmpPDU) error {
		// Dari Nama OID yang didapat, kita ambil angka terakhirnya sebagai ID ONU.
		onuID := extractOnuID(pdu.Name)
		if onuID == 0 {
			return nil
		}

		info := &model.ONUInfo{
			Board: boardID,
			PON:   ponID,
			ID:    onuID,
			Name:  extractString(pdu.Value), // Nama ONU (biasanya diinput teknisi)
		}
		onuMap[onuID] = info
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("SNMP walk failed: %w", err)
	}

	// 2. Ambil informasi tambahan (Tipe, SN, Sinyal, Status) untuk setiap ONU yang ditemukan.
	for onuID, info := range onuMap {
		onuIDStr := strconv.Itoa(onuID)
		
		// Ambil Tipe ONU
		if val, err := d.snmpGet(BaseOID2 + cfg.OnuTypeOID + "." + onuIDStr); err == nil {
			info.Type = extractString(val)
		}

		// Ambil Serial Number
		if val, err := d.snmpGet(BaseOID1 + cfg.OnuSerialNumberOID + "." + onuIDStr); err == nil {
			info.SerialNumber = extractSerialNumber(val)
		}

		// Ambil Kekuatan Sinyal (RX Power)
		if val, err := d.snmpGet(BaseOID1 + cfg.OnuRxPowerOID + "." + onuIDStr + ".1"); err == nil {
			info.RXPower = convertPower(val)
		}

		// Ambil Kekuatan Sinyal (TX Power)
		if val, err := d.snmpGet(BaseOID2 + cfg.OnuTxPowerOID + "." + onuIDStr + ".1"); err == nil {
			info.TXPower = convertPower(val)
		}

		// Ambil Jarak (Distance)
		if val, err := d.snmpGet(BaseOID1 + cfg.OnuGponOpticalDistanceOID + "." + onuIDStr); err == nil {
			info.Distance = fmt.Sprintf("%v", val)
		}

		// Ambil Status (Online/Offline)
		if val, err := d.snmpGet(BaseOID1 + cfg.OnuStatusOID + "." + onuIDStr); err == nil {
			info.Status = convertStatus(val)
		}

		onuList = append(onuList, *info)
	}

	return onuList, nil
}

// GetONUDetail mengambil informasi detail untuk satu ONU tunggal.
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

	// Ambil Nama
	if val, err := d.snmpGet(BaseOID1 + cfg.OnuIDNameOID + "." + onuIDStr); err == nil {
		detail.Name = extractString(val)
	}

	// Ambil Tipe
	if val, err := d.snmpGet(BaseOID2 + cfg.OnuTypeOID + "." + onuIDStr); err == nil {
		detail.Type = extractString(val)
	}

	// Ambil Serial Number
	if val, err := d.snmpGet(BaseOID1 + cfg.OnuSerialNumberOID + "." + onuIDStr); err == nil {
		detail.SerialNumber = extractSerialNumber(val)
	}

	// Ambil Sinyal RX
	if val, err := d.snmpGet(BaseOID1 + cfg.OnuRxPowerOID + "." + onuIDStr + ".1"); err == nil {
		detail.RXPower = convertPower(val)
	}

	// Ambil Sinyal TX
	if val, err := d.snmpGet(BaseOID2 + cfg.OnuTxPowerOID + "." + onuIDStr + ".1"); err == nil {
		detail.TXPower = convertPower(val)
	}

	// Ambil Status
	if val, err := d.snmpGet(BaseOID1 + cfg.OnuStatusOID + "." + onuIDStr); err == nil {
		detail.Status = convertStatus(val)
	}

	// Ambil Alamat IP
	if val, err := d.snmpGet(BaseOID2 + cfg.OnuIPAddressOID + "." + onuIDStr + ".1"); err == nil {
		detail.IPAddress = extractString(val)
	}

	// Ambil Deskripsi
	if val, err := d.snmpGet(BaseOID1 + cfg.OnuDescriptionOID + "." + onuIDStr); err == nil {
		detail.Description = extractString(val)
	}

	// Ambil Waktu Terakhir Online
	if val, err := d.snmpGet(BaseOID1 + cfg.OnuLastOnlineOID + "." + onuIDStr); err == nil {
		detail.LastOnline = convertDateTime(val)
	}

	// Ambil Waktu Terakhir Offline
	if val, err := d.snmpGet(BaseOID1 + cfg.OnuLastOfflineOID + "." + onuIDStr); err == nil {
		detail.LastOffline = convertDateTime(val)
	}

	// Ambil Alasan Offline
	if val, err := d.snmpGet(BaseOID1 + cfg.OnuLastOfflineReasonOID + "." + onuIDStr); err == nil {
		detail.OfflineReason = convertOfflineReason(val)
	}

	// Ambil Jarak
	if val, err := d.snmpGet(BaseOID1 + cfg.OnuGponOpticalDistanceOID + "." + onuIDStr); err == nil {
		detail.Distance = fmt.Sprintf("%v", val)
	}

	// Hitung Uptime
	if detail.LastOnline != "" {
		detail.Uptime = calculateUptime(detail.LastOnline)
	}

	return detail, nil
}

// GetEmptySlots mengambil slot ONU yang masih kosong.
func (d *Driver) GetEmptySlots(ctx context.Context, boardID, ponID int) ([]model.ONUSlot, error) {
	if !d.connected {
		if err := d.Connect(); err != nil {
			return nil, err
		}
	}

	cfg := GenerateBoardPonOID(boardID, ponID)
	
	// Lacak ID ONU yang sudah terpakai
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

	// Cari slot yang kosong
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

// GetSystemInfo mengambil informasi sistem OLT.
func (d *Driver) GetSystemInfo(ctx context.Context) (*driver.SystemInfo, error) {
	if !d.connected {
		if err := d.Connect(); err != nil {
			return nil, err
		}
	}

	info := &driver.SystemInfo{}

	// Deskripsi Sistem
	if val, err := d.snmpGet("1.3.6.1.2.1.1.1.0"); err == nil {
		info.Description = extractString(val)
	}

	// Nama Sistem
	if val, err := d.snmpGet("1.3.6.1.2.1.1.5.0"); err == nil {
		info.Name = extractString(val)
	}

	// Uptime Sistem
	if val, err := d.snmpGet("1.3.6.1.2.1.1.3.0"); err == nil {
		info.Uptime = fmt.Sprintf("%v", val)
	}

	// Kontak Sistem
	if val, err := d.snmpGet("1.3.6.1.2.1.1.4.0"); err == nil {
		info.Contact = extractString(val)
	}

	// Lokasi Sistem
	if val, err := d.snmpGet("1.3.6.1.2.1.1.6.0"); err == nil {
		info.Location = extractString(val)
	}

	return info, nil
}

// GetBoardInfo mengambil informasi board/kartu.
func (d *Driver) GetBoardInfo(ctx context.Context, boardID int) (*model.BoardInfo, error) {
	if !d.connected {
		if err := d.Connect(); err != nil {
			return nil, err
		}
	}

	info := &model.BoardInfo{BoardID: boardID}

	// Format OID dari hasil snmpwalk:
	// .1.3.6.1.4.1.3902.1015.2.1.1.3.1.{attr}.1.1.{position}
	// position 1 = board 1 (GTGO), position 3 = board 2 (PRAM), position 4 = board 3 (SMXA)

	// Cari posisi berdasarkan board yang tersedia
	// Board 1 = position 1, Board 2 = position 3 (skip slot kosong)
	position := boardID
	if boardID == 2 {
		position = 3 // Slot 2 kosong, langsung ke slot 3
	}

	// Ambil Tipe Kartu (attribute 4)
	typeOID := fmt.Sprintf("%s.2.1.1.3.1.4.1.1.%d", BaseOID3, position)
	if val, err := d.snmpGet(typeOID); err == nil {
		info.Type = extractString(val)
	}

	// Ambil Status Kartu (attribute 5)
	statusOID := fmt.Sprintf("%s.2.1.1.3.1.5.1.1.%d", BaseOID3, position)
	if val, err := d.snmpGet(statusOID); err == nil {
		if intVal := extractInt(val); intVal > 0 {
			info.Status = model.CardStatus(intVal).String()
		}
	}

	// Ambil Beban CPU (attribute 9)
	cpuOID := fmt.Sprintf("%s.2.1.1.3.1.9.1.1.%d", BaseOID3, position)
	if val, err := d.snmpGet(cpuOID); err == nil {
		info.CpuLoad = extractInt(val)
	}

	// Ambil Penggunaan Memori (attribute 11)
	memOID := fmt.Sprintf("%s.2.1.1.3.1.11.1.1.%d", BaseOID3, position)
	if val, err := d.snmpGet(memOID); err == nil {
		info.MemUsage = extractInt(val)
	}

	return info, nil
}

// GetAllBoards mengambil semua informasi board.
func (d *Driver) GetAllBoards(ctx context.Context) ([]model.BoardInfo, error) {
	var boards []model.BoardInfo
	for i := 1; i <= MaxBoards; i++ {
		info, err := d.GetBoardInfo(ctx, i)
		if err != nil {
			continue
		}
		boards = append(boards, *info)
	}
	return boards, nil
}

// GetPONInfo mengambil informasi port PON.
func (d *Driver) GetPONInfo(ctx context.Context, boardID, ponID int) (*model.PONInfo, error) {
	if !d.connected {
		if err := d.Connect(); err != nil {
			return nil, err
		}
	}

	info := &model.PONInfo{
		BoardID: boardID,
		PonID:   ponID,
	}

	// Ambil ONUs pada port PON ini
	onus, err := d.GetONUList(ctx, boardID, ponID)
	if err == nil {
		info.ONUCount = len(onus)

		// Hitung average power jika ada ONU
		if len(onus) > 0 {
			var totalTX, totalRX float64
			for _, onu := range onus {
				if tx, err := strconv.ParseFloat(onu.TXPower, 64); err == nil {
					totalTX += tx
				}
				if rx, err := strconv.ParseFloat(onu.RXPower, 64); err == nil {
					totalRX += rx
				}
			}
			info.TxPower = totalTX / float64(len(onus))
			info.RxPower = totalRX / float64(len(onus))
		}
	}

	return info, nil
}

// GetONUTraffic mengambil statistik trafik ONU (Placeholder - butuh OID spesifik).
func (d *Driver) GetONUTraffic(ctx context.Context, boardID, ponID, onuID int) (*model.ONUTraffic, error) {
	if !d.connected {
		if err := d.Connect(); err != nil {
			return nil, err
		}
	}

	traffic := &model.ONUTraffic{
		Board:     boardID,
		PON:       ponID,
		ONUID:     onuID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	// Catatan: OID trafik ONU bervariasi tergantung konfigurasi OLT.
	// Ini adalah implementasi sementara (placeholder).
	return traffic, nil
}

// GetInterfaceStats mengambil statistik interface.
func (d *Driver) GetInterfaceStats(ctx context.Context) ([]model.InterfaceStats, error) {
	if !d.connected {
		if err := d.Connect(); err != nil {
			return nil, err
		}
	}

	var stats []model.InterfaceStats
	indexMap := make(map[int]*model.InterfaceStats)

	// Walk deskripsi interface
	d.client.Walk("1.3.6.1.2.1.2.2.1.2", func(pdu gosnmp.SnmpPDU) error {
		idx := extractLastOIDPart(pdu.Name)
		if idx > 0 {
			indexMap[idx] = &model.InterfaceStats{
				Index:       idx,
				Description: extractString(pdu.Value),
			}
		}
		return nil
	})

	// Walk status interface
	d.client.Walk("1.3.6.1.2.1.2.2.1.8", func(pdu gosnmp.SnmpPDU) error {
		idx := extractLastOIDPart(pdu.Name)
		if stat, ok := indexMap[idx]; ok {
			if intVal, ok := pdu.Value.(int); ok {
				if intVal == 1 {
					stat.Status = "Up"
				} else {
					stat.Status = "Down"
				}
			}
		}
		return nil
	})

	// Walk byte RX
	d.client.Walk("1.3.6.1.2.1.2.2.1.10", func(pdu gosnmp.SnmpPDU) error {
		idx := extractLastOIDPart(pdu.Name)
		if stat, ok := indexMap[idx]; ok {
			stat.RxBytes = extractCounter64(pdu.Value)
		}
		return nil
	})

	// Walk byte TX
	d.client.Walk("1.3.6.1.2.1.2.2.1.16", func(pdu gosnmp.SnmpPDU) error {
		idx := extractLastOIDPart(pdu.Name)
		if stat, ok := indexMap[idx]; ok {
			stat.TxBytes = extractCounter64(pdu.Value)
		}
		return nil
	})

	// Konversi ke slice
	for _, stat := range indexMap {
		stats = append(stats, *stat)
	}

	return stats, nil
}

// GetFanInfo mengambil informasi fan
func (d *Driver) GetFanInfo(ctx context.Context) ([]map[string]interface{}, error) {
	if !d.connected {
		if err := d.Connect(); err != nil {
			return nil, err
		}
	}

	var fans []map[string]interface{}
	fanMap := make(map[int]map[string]interface{})

	// Walk fan table to get indices
	err := d.client.Walk(FanTableOID, func(pdu gosnmp.SnmpPDU) error {
		// Extract fan index from OID
		parts := strings.Split(pdu.Name, ".")
		if len(parts) < 1 {
			return nil
		}
		
		idxStr := parts[len(parts)-1]
		idx, err := strconv.Atoi(idxStr)
		if err != nil {
			// Try second to last
			if len(parts) >= 2 {
				idxStr = parts[len(parts)-2]
				idx, _ = strconv.Atoi(idxStr)
			}
		}
		
		if idx == 0 {
			idx = len(fanMap) + 1
		}

		if _, exists := fanMap[idx]; !exists {
			fanMap[idx] = map[string]interface{}{
				"index": idx,
			}
		}
		
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Get details for each fan
	for idx := range fanMap {
		fan := map[string]interface{}{
			"index": idx,
		}

		// Get speed level
		speedOID := fmt.Sprintf("%s.%d", FanSpeedLevelOID, idx)
		if val, err := d.snmpGet(speedOID); err == nil {
			if intVal := extractInt(val); intVal > 0 {
				fan["speed_level"] = intVal
				// Convert to string
				switch intVal {
				case 1:
					fan["speed"] = "Low"
				case 2:
					fan["speed"] = "Standard"
				case 3:
					fan["speed"] = "High"
				case 4:
					fan["speed"] = "Super"
				default:
					fan["speed"] = "Unknown"
				}
			}
		}

		// Get status
		statusOID := fmt.Sprintf("%s.%d", FanStatusOID, idx)
		if val, err := d.snmpGet(statusOID); err == nil {
			if intVal := extractInt(val); intVal == 1 {
				fan["status"] = "Normal"
			} else {
				fan["status"] = "Abnormal"
			}
		}

		// Get present
		presentOID := fmt.Sprintf("%s.%d", FanPresentOID, idx)
		if val, err := d.snmpGet(presentOID); err == nil {
			if extractInt(val) == 1 {
				fan["present"] = true
			} else {
				fan["present"] = false
			}
		}

		fans = append(fans, fan)
	}

	return fans, nil
}

// snmpGet melakukan permintaan SNMP GET.
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

// Fungsi Pembantu (Helpers)

func extractOnuID(oid string) int {
	return extractLastOIDPart(oid)
}

func extractLastOIDPart(oid string) int {
	parts := splitOID(oid)
	if len(parts) < 1 {
		return 0
	}
	lastPart := parts[len(parts)-1]
	id, _ := strconv.Atoi(lastPart)
	return id
}

func extractOnuIDSuffix(boardID, ponID int) string {
	var baseOnuID int
	if boardID == 1 {
		baseOnuID = Board1OnuIDBase
	} else {
		baseOnuID = Board2OnuIDBase
	}
	suffix := baseOnuID + (ponID * OnuIDIncrement)
	return strconv.Itoa(suffix)
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

func extractInt(val interface{}) int {
	switch v := val.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case uint:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	default:
		return 0
	}
}

func extractSerialNumber(val interface{}) string {
	s := extractString(val)
	if len(s) > 2 && s[:2] == "1," {
		return s[2:]
	}
	return s
}

func extractCounter64(val interface{}) int64 {
	switch v := val.(type) {
	case uint:
		return int64(v)
	case uint64:
		return int64(v)
	case int:
		return int64(v)
	case int64:
		return v
	default:
		return 0
	}
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
	year := int(bytes[0])<<8 | int(bytes[1])
	month := int(bytes[2])
	day := int(bytes[3])
	hour := int(bytes[4])
	minute := int(bytes[5])
	second := int(bytes[6])
	
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, minute, second)
}

func calculateUptime(lastOnline string) string {
	return lastOnline
}

// Pastikan Driver mengimplementasikan interface driver.Driver
var _ driver.Driver = (*Driver)(nil)
