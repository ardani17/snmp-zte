package cli

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ZTEC320Client khusus untuk ZTE C320 CLI commands
type ZTEC320Client struct {
	client *Client
}

// NewZTEC320Client membuat client ZTE C320
func NewZTEC320Client(cfg Config) *ZTEC320Client {
	return &ZTEC320Client{
		client: New(cfg),
	}
}

// Connect melakukan koneksi
func (z *ZTEC320Client) Connect() error {
	return z.client.Connect()
}

// Close menutup koneksi
func (z *ZTEC320Client) Close() error {
	return z.client.Close()
}

// ============================================================
// SHOW COMMANDS
// ============================================================

// ONUInfo informasi ONU
type ONUInfo struct {
	Index      string `json:"index"`
	SN         string `json:"sn"`
	State      string `json:"state"`
	AdminState string `json:"admin_state,omitempty"`
	OmccState  string `json:"omcc_state,omitempty"`
	O7State    string `json:"o7_state,omitempty"`
	PhaseState string `json:"phase_state,omitempty"`
	Type       string `json:"type,omitempty"`
	Name       string `json:"name,omitempty"`
}

// UncfgONU ONU yang belum terdaftar
type UncfgONU struct {
	Index string `json:"index"`
	SN    string `json:"sn"`
	State string `json:"state"`
}

// ShowGPONONUUnCfg menampilkan ONU yang belum terautentikasi
// Command: show gpon onu uncfg gpon-olt_{rack}/{shelf}/{slot}
func (z *ZTEC320Client) ShowGPONONUUnCfg(ctx context.Context, rack, shelf, slot int) ([]UncfgONU, error) {
	cmd := fmt.Sprintf("show gpon onu uncfg gpon-olt_%d/%d/%d", rack, shelf, slot)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return z.parseUncfgONU(output), nil
}

// ShowGPONONUState menampilkan state ONU
// Command: show gpon onu state gpon-olt_{rack}/{shelf}/{slot}
func (z *ZTEC320Client) ShowGPONONUState(ctx context.Context, rack, shelf, slot int) ([]ONUInfo, error) {
	cmd := fmt.Sprintf("show gpon onu state gpon-olt_%d/%d/%d", rack, shelf, slot)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return z.parseONUState(output), nil
}

// CardInfo informasi card
type CardInfo struct {
	Rack      int    `json:"rack"`
	Shelf     int    `json:"shelf"`
	Slot      int    `json:"slot"`
	CfgType   string `json:"cfg_type"`
	RealType  string `json:"real_type"`
	Port      int    `json:"port"`
	HardVer   string `json:"hard_ver"`
	SoftVer   string `json:"soft_ver"`
	Status    string `json:"status"`
	CPUUsage  string `json:"cpu_usage,omitempty"`
	MemUsage  string `json:"mem_usage,omitempty"`
	Uptime    string `json:"uptime,omitempty"`
	SerialNum string `json:"serial_number,omitempty"`
}

// ShowCard menampilkan semua card
// Command: show card
func (z *ZTEC320Client) ShowCard(ctx context.Context) ([]CardInfo, error) {
	output, err := z.client.Execute(ctx, "show card")
	if err != nil {
		return nil, err
	}

	return z.parseCardInfo(output), nil
}

// ShowCardSlot menampilkan card tertentu
// Command: show card slotno {slot}
func (z *ZTEC320Client) ShowCardSlot(ctx context.Context, slot int) (*CardInfo, error) {
	cmd := fmt.Sprintf("show card slotno %d", slot)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}

	cards := z.parseCardDetail(output)
	if len(cards) > 0 {
		return &cards[0], nil
	}
	return nil, fmt.Errorf("card not found at slot %d", slot)
}

// TCONTProfile T-CONT bandwidth profile
type TCONTProfile struct {
	Name    string `json:"name"`
	Type    int    `json:"type"`
	FBW     int    `json:"fbw_kbps"`
	ABW     int    `json:"abw_kbps"`
	MBW     int    `json:"mbw_kbps"`
}

// ShowGPONProfileTcont menampilkan T-CONT profiles
// Command: show gpon profile tcont
func (z *ZTEC320Client) ShowGPONProfileTcont(ctx context.Context) ([]TCONTProfile, error) {
	output, err := z.client.Execute(ctx, "show gpon profile tcont")
	if err != nil {
		return nil, err
	}

	return z.parseTCONTProfile(output), nil
}

// ONUTypeInfo informasi tipe ONU
type ONUTypeInfo struct {
	Name              string `json:"name"`
	PonType           string `json:"pon_type"`
	Description       string `json:"description"`
	MaxTcont          int    `json:"max_tcont"`
	MaxGemport        int    `json:"max_gemport"`
	MaxSwitchPerSlot  int    `json:"max_switch_per_slot"`
	MaxFlowPerSwitch  int    `json:"max_flow_per_switch"`
}

// ShowONUType menampilkan tipe ONU
// Command: show onu-type gpon {type}
func (z *ZTEC320Client) ShowONUType(ctx context.Context, onuType string) (*ONUTypeInfo, error) {
	cmd := fmt.Sprintf("show onu-type gpon %s", onuType)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return z.parseONUType(output), nil
}

// ShowONUTypeList menampilkan list tipe ONU
// Command: show onu-type gpon
func (z *ZTEC320Client) ShowONUTypeList(ctx context.Context) ([]string, error) {
	output, err := z.client.Execute(ctx, "show onu-type gpon")
	if err != nil {
		return nil, err
	}

	return z.parseONUTypeList(output), nil
}

// FanInfo informasi fan
type FanInfo struct {
	ControlType            string `json:"control_type"`
	TempThreshold          string `json:"temp_threshold"`
	FanSpeedPercent        string `json:"fan_speed_percent"`
	HighTempThreshold      string `json:"high_temp_threshold"`
	EnvironmentTemperature string `json:"environment_temperature"`
}

// ShowFan menampilkan info fan
// Command: show fan
func (z *ZTEC320Client) ShowFan(ctx context.Context) (*FanInfo, error) {
	output, err := z.client.Execute(ctx, "show fan")
	if err != nil {
		return nil, err
	}

	return z.parseFanInfo(output), nil
}

// SystemInfo informasi sistem
type SystemInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Uptime      string `json:"uptime"`
	Time        string `json:"time"`
}

// ShowVersion menampilkan versi sistem
// Command: show version
func (z *ZTEC320Client) ShowVersion(ctx context.Context) (*SystemInfo, error) {
	output, err := z.client.Execute(ctx, "show version")
	if err != nil {
		return nil, err
	}

	return z.parseSystemInfo(output), nil
}

// ShowClock menampilkan waktu sistem
// Command: show clock
func (z *ZTEC320Client) ShowClock(ctx context.Context) (string, error) {
	output, err := z.client.Execute(ctx, "show clock")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}

// ShowRunningConfig menampilkan running config
// Command: show running-config
func (z *ZTEC320Client) ShowRunningConfig(ctx context.Context) (string, error) {
	output, err := z.client.Execute(ctx, "show running-config")
	if err != nil {
		return "", err
	}

	return output, nil
}

// ============================================================
// PARSING FUNCTIONS
// ============================================================

func (z *ZTEC320Client) parseUncfgONU(output string) []UncfgONU {
	var onus []UncfgONU
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "OnuIndex") || strings.HasPrefix(line, "---") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 3 {
			onus = append(onus, UncfgONU{
				Index: fields[0],
				SN:    fields[1],
				State: fields[2],
			})
		}
	}

	return onus
}

func (z *ZTEC320Client) parseONUState(output string) []ONUInfo {
	var onus []ONUInfo
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "OnuIndex") || strings.HasPrefix(line, "---") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 5 {
			onus = append(onus, ONUInfo{
				Index:      fields[0],
				AdminState: fields[1],
				OmccState:  fields[2],
				O7State:    fields[3],
				PhaseState: fields[4],
			})
		}
	}

	return onus
}

func (z *ZTEC320Client) parseCardInfo(output string) []CardInfo {
	var cards []CardInfo
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Rack") || strings.HasPrefix(line, "---") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 9 {
			rack, _ := strconv.Atoi(fields[0])
			shelf, _ := strconv.Atoi(fields[1])
			slot, _ := strconv.Atoi(fields[2])
			port, _ := strconv.Atoi(fields[5])

			cards = append(cards, CardInfo{
				Rack:     rack,
				Shelf:    shelf,
				Slot:     slot,
				CfgType:  fields[3],
				RealType: fields[4],
				Port:     port,
				HardVer:  fields[6],
				SoftVer:  fields[7],
				Status:   fields[8],
			})
		}
	}

	return cards
}

func (z *ZTEC320Client) parseCardDetail(output string) []CardInfo {
	var cards []CardInfo
	card := CardInfo{}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Config-Type") {
			card.CfgType = z.extractValue(line)
		} else if strings.Contains(line, "Real-Type") {
			card.RealType = z.extractValue(line)
		} else if strings.Contains(line, "Status") {
			card.Status = z.extractValue(line)
		} else if strings.Contains(line, "Software-VER") {
			card.SoftVer = z.extractValue(line)
		} else if strings.Contains(line, "PCB-VER") {
			card.HardVer = z.extractValue(line)
		} else if strings.Contains(line, "Cpu-Usage") {
			card.CPUUsage = z.extractValue(line)
		} else if strings.Contains(line, "Mem-Usage") {
			card.MemUsage = z.extractValue(line)
		} else if strings.Contains(line, "Uptime") {
			card.Uptime = z.extractValue(line)
		} else if strings.Contains(line, "Serial-Number") {
			card.SerialNum = z.extractValue(line)
		}
	}

	if card.CfgType != "" {
		cards = append(cards, card)
	}
	return cards
}

func (z *ZTEC320Client) parseTCONTProfile(output string) []TCONTProfile {
	var profiles []TCONTProfile
	var current *TCONTProfile

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Name") {
			if current != nil {
				profiles = append(profiles, *current)
			}
			current = &TCONTProfile{
				Name: z.extractValue(line),
			}
		} else if current != nil && strings.HasPrefix(line, "Type") {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				current.Type, _ = strconv.Atoi(fields[1])
				current.FBW, _ = strconv.Atoi(fields[2])
				current.ABW, _ = strconv.Atoi(fields[3])
				if len(fields) >= 5 {
					current.MBW, _ = strconv.Atoi(fields[4])
				}
			}
		}
	}

	if current != nil {
		profiles = append(profiles, *current)
	}
	return profiles
}

func (z *ZTEC320Client) parseONUType(output string) *ONUTypeInfo {
	info := &ONUTypeInfo{}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Onu type name") {
			info.Name = z.extractValue(line)
		} else if strings.Contains(line, "Pon type") {
			info.PonType = z.extractValue(line)
		} else if strings.Contains(line, "Description") {
			info.Description = z.extractValue(line)
		} else if strings.Contains(line, "Max tcont") {
			info.MaxTcont, _ = strconv.Atoi(z.extractValue(line))
		} else if strings.Contains(line, "Max gemport") {
			info.MaxGemport, _ = strconv.Atoi(z.extractValue(line))
		} else if strings.Contains(line, "Max switch per slot") {
			info.MaxSwitchPerSlot, _ = strconv.Atoi(z.extractValue(line))
		} else if strings.Contains(line, "Max flow per switch") {
			info.MaxFlowPerSwitch, _ = strconv.Atoi(z.extractValue(line))
		}
	}

	return info
}

func (z *ZTEC320Client) parseONUTypeList(output string) []string {
	var types []string
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Onu type") || strings.HasPrefix(line, "---") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 0 && strings.HasPrefix(fields[0], "ZTE") {
			types = append(types, fields[0])
		}
	}

	return types
}

func (z *ZTEC320Client) parseFanInfo(output string) *FanInfo {
	info := &FanInfo{}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "FanControlType") {
			info.ControlType = z.extractValue(line)
		} else if strings.Contains(line, "TemperatureThreshold") {
			info.TempThreshold = z.extractValue(line)
		} else if strings.Contains(line, "FanSpeedLevelPercent") {
			info.FanSpeedPercent = z.extractValue(line)
		} else if strings.Contains(line, "HighTemperatureThreshold") {
			info.HighTempThreshold = z.extractValue(line)
		} else if strings.Contains(line, "Environment Temperature") {
			info.EnvironmentTemperature = z.extractValue(line)
		}
	}

	return info
}

func (z *ZTEC320Client) parseSystemInfo(output string) *SystemInfo {
	info := &SystemInfo{}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Product name") {
			info.Name = z.extractValue(line)
		} else if strings.Contains(line, "uptime") {
			info.Uptime = z.extractValue(line)
		}
	}

	return info
}

func (z *ZTEC320Client) extractValue(line string) string {
	// Extract value after colon
	parts := strings.SplitN(line, ":", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

// ============================================================
// CONFIGURATION COMMANDS
// ============================================================

// AuthenticateONU mendaftarkan ONU baru
// Command: onu {id} type {type} sn {sn}
func (z *ZTEC320Client) AuthenticateONU(ctx context.Context, rack, shelf, slot, onuID int, onuType, sn string) error {
	commands := []string{
		"configure terminal",
		fmt.Sprintf("interface gpon-olt_%d/%d/%d", rack, shelf, slot),
		fmt.Sprintf("onu %d type %s sn %s", onuID, onuType, sn),
		"exit",
		"exit",
	}

	for _, cmd := range commands {
		_, err := z.client.Execute(ctx, cmd)
		if err != nil {
			return fmt.Errorf("command '%s' failed: %w", cmd, err)
		}
		time.Sleep(50 * time.Millisecond)
	}

	return nil
}

// DeleteONU menghapus ONU
// Command: no onu {id}
func (z *ZTEC320Client) DeleteONU(ctx context.Context, rack, shelf, slot, onuID int) error {
	commands := []string{
		"configure terminal",
		fmt.Sprintf("interface gpon-olt_%d/%d/%d", rack, shelf, slot),
		fmt.Sprintf("no onu %d", onuID),
		"exit",
		"exit",
	}

	for _, cmd := range commands {
		_, err := z.client.Execute(ctx, cmd)
		if err != nil {
			return fmt.Errorf("command '%s' failed: %w", cmd, err)
		}
		time.Sleep(50 * time.Millisecond)
	}

	return nil
}

// SaveConfig menyimpan konfigurasi
// Command: write
func (z *ZTEC320Client) SaveConfig(ctx context.Context) error {
	_, err := z.client.Execute(ctx, "write")
	return err
}

// RawCommand menjalankan command mentah
func (z *ZTEC320Client) RawCommand(ctx context.Context, cmd string) (string, error) {
	return z.client.Execute(ctx, cmd)
}

// ValidateSN memvalidasi format Serial Number ONU
func ValidateSN(sn string) bool {
	// ZTE SN format: ZTEG00000002 (4 letters + 8 hex)
	matched, _ := regexp.MatchString(`^[A-Z]{4}[0-9A-Fa-f]{8}$`, sn)
	return matched
}
