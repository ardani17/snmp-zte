package cli

import (
	"context"
	"fmt"
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

// Execute menjalankan command mentah (wrapper for client.Execute)
func (z *ZTEC320Client) Execute(ctx context.Context, cmd string) (string, error) {
	return z.client.Execute(ctx, cmd)
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

// ============================================================
// WRITE OPERATIONS - ONU PROVISIONING
// ============================================================

// RenameONU mengubah nama ONU
// Command: onu {id} name {name}
func (z *ZTEC320Client) RenameONU(ctx context.Context, rack, shelf, slot, onuID int, newName string) error {
	commands := []string{
		"configure terminal",
		fmt.Sprintf("interface gpon-olt_%d/%d/%d", rack, shelf, slot),
		fmt.Sprintf("onu %d name %s", onuID, newName),
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

// ResetONU mereset/reboot ONU
// Command: reset onu {id}
func (z *ZTEC320Client) ResetONU(ctx context.Context, rack, shelf, slot, onuID int) error {
	commands := []string{
		"configure terminal",
		fmt.Sprintf("interface gpon-olt_%d/%d/%d", rack, shelf, slot),
		fmt.Sprintf("reset onu %d", onuID),
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

// ============================================================
// WRITE OPERATIONS - T-CONT & GEM PORT
// ============================================================

// CreateTCONT membuat T-CONT baru
// Command: tcont {id} name {name} profile {profile}
func (z *ZTEC320Client) CreateTCONT(ctx context.Context, rack, shelf, slot, onuID, tcontID int, name, profile string) error {
	commands := []string{
		"configure terminal",
		fmt.Sprintf("interface gpon-onu_%d/%d/%d:%d", rack, shelf, slot, onuID),
		fmt.Sprintf("tcont %d name %s profile %s", tcontID, name, profile),
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

// CreateGEMPort membuat GEM port baru
// Command: gemport {id} unicast tcont {tcont_id}
func (z *ZTEC320Client) CreateGEMPort(ctx context.Context, rack, shelf, slot, onuID, gemportID, tcontID int, name string) error {
	commands := []string{
		"configure terminal",
		fmt.Sprintf("interface gpon-onu_%d/%d/%d:%d", rack, shelf, slot, onuID),
		fmt.Sprintf("gemport %d name %s unicast tcont %d", gemportID, name, tcontID),
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

// ============================================================
// WRITE OPERATIONS - SERVICE PORT
// ============================================================

// CreateServicePort membuat service port baru
// Command: service-port {id} vport {vport} user-vlan {vlan} vlan {vlan}
func (z *ZTEC320Client) CreateServicePort(ctx context.Context, rack, shelf, slot, onuID, servicePortID, vport, vlan int) error {
	commands := []string{
		"configure terminal",
		fmt.Sprintf("interface gpon-onu_%d/%d/%d:%d", rack, shelf, slot, onuID),
		fmt.Sprintf("service-port %d vport %d user-vlan %d vlan %d", servicePortID, vport, vlan, vlan),
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

// DeleteServicePort menghapus service port
// Command: no service-port {id}
func (z *ZTEC320Client) DeleteServicePort(ctx context.Context, rack, shelf, slot, onuID, servicePortID int) error {
	commands := []string{
		"configure terminal",
		fmt.Sprintf("interface gpon-onu_%d/%d/%d:%d", rack, shelf, slot, onuID),
		fmt.Sprintf("no service-port %d", servicePortID),
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

// ============================================================
// WRITE OPERATIONS - VLAN
// ============================================================

// CreateVLAN membuat VLAN baru
// Command: vlan {id}
func (z *ZTEC320Client) CreateVLAN(ctx context.Context, vlanID int, name string) error {
	commands := []string{
		"configure terminal",
		fmt.Sprintf("vlan %d", vlanID),
	}

	if name != "" {
		commands = append(commands, fmt.Sprintf("name %s", name))
	}

	commands = append(commands, "exit", "exit")

	for _, cmd := range commands {
		_, err := z.client.Execute(ctx, cmd)
		if err != nil {
			return fmt.Errorf("command '%s' failed: %w", cmd, err)
		}
		time.Sleep(50 * time.Millisecond)
	}

	return nil
}

// DeleteVLAN menghapus VLAN
// Command: no vlan {id}
func (z *ZTEC320Client) DeleteVLAN(ctx context.Context, vlanID int) error {
	commands := []string{
		"configure terminal",
		fmt.Sprintf("no vlan %d", vlanID),
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

// AddPortToVLAN menambahkan port ke VLAN
// Command: switchport vlan {vlan_id} tag/untag
func (z *ZTEC320Client) AddPortToVLAN(ctx context.Context, interfaceName string, vlanID int, mode string) error {
	commands := []string{
		"configure terminal",
		fmt.Sprintf("interface %s", interfaceName),
		fmt.Sprintf("switchport vlan %d %s", vlanID, mode),
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

// ============================================================
// WRITE OPERATIONS - PROFILE CREATION
// ============================================================

// CreateLineProfile membuat line profile baru
// Command: onu-profile gpon line {name}
func (z *ZTEC320Client) CreateLineProfile(ctx context.Context, name string) error {
	commands := []string{
		"configure terminal",
		fmt.Sprintf("onu-profile gpon line %s", name),
		"commit",
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

// CreateRemoteProfile membuat remote profile baru
// Command: onu-profile gpon remote {name}
func (z *ZTEC320Client) CreateRemoteProfile(ctx context.Context, name string) error {
	commands := []string{
		"configure terminal",
		fmt.Sprintf("onu-profile gpon remote %s", name),
		"commit",
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

// CreateVLANProfile membuat VLAN profile baru
// Command: onu profile vlan {name}
func (z *ZTEC320Client) CreateVLANProfile(ctx context.Context, name string, vlanID int) error {
	commands := []string{
		"configure terminal",
		fmt.Sprintf("onu profile vlan %s", name),
		fmt.Sprintf("vlan %d", vlanID),
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

// CreateTCONTProfile membuat T-CONT profile baru
// Command: profile tcont {name} type {type} bandwidth {bandwidth}
func (z *ZTEC320Client) CreateTCONTProfile(ctx context.Context, name string, profileType string, bandwidth int) error {
	commands := []string{
		"configure terminal",
		fmt.Sprintf("profile tcont %s type %s bandwidth %d", name, profileType, bandwidth),
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

// ============================================================
// WRITE OPERATIONS - IGMP/MULTICAST
// ============================================================

// EnableIGMP mengaktifkan IGMP
// Command: igmp enable
func (z *ZTEC320Client) EnableIGMP(ctx context.Context) error {
	commands := []string{
		"configure terminal",
		"igmp enable",
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

// CreateMVLAN membuat MVLAN baru
// Command: igmp mvlan {id}
func (z *ZTEC320Client) CreateMVLAN(ctx context.Context, mvlanID int) error {
	commands := []string{
		"configure terminal",
		fmt.Sprintf("igmp mvlan %d", mvlanID),
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

// AddMVLANGroup menambahkan group ke MVLAN
// Command: igmp mvlan {id} group {ip}
func (z *ZTEC320Client) AddMVLANGroup(ctx context.Context, mvlanID int, groupIP string) error {
	commands := []string{
		"configure terminal",
		fmt.Sprintf("igmp mvlan %d group %s", mvlanID, groupIP),
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

// ============================================================
// PRIORITY 1: ONU DETAIL ENDPOINTS
// ============================================================

// ONUDetail informasi detail ONU
type ONUDetail struct {
	ONUName          string `json:"onu_name"`
	ONUType          string `json:"onu_type"`
	ONUSN            string `json:"onu_sn"`
	AdminState       string `json:"admin_state"`
	PhaseState       string `json:"phase_state"`
	ChannelState     string `json:"channel_state"`
	Authentication   string `json:"authentication"`
	Omcc             string `json:"omcc"`
	VoipState        string `json:"voip_state,omitempty"`
	VoipPortNum      string `json:"voip_port_num,omitempty"`
	FecUp            string `json:"fec_up"`
	FecDown          string `json:"fec_down"`
	LineProfile      string `json:"line_profile"`
	RemoteProfile    string `json:"remote_profile"`
	Description      string `json:"description,omitempty"`
	LastDownReason   string `json:"last_down_reason,omitempty"`
	LastDownTime     string `json:"last_down_time,omitempty"`
	DyingGasp        string `json:"dying_gasp,omitempty"`
	BatteryBackup    string `json:"battery_backup,omitempty"`
	ONUVersion       string `json:"onu_version,omitempty"`
	EquipmentID      string `json:"equipment_id,omitempty"`
	TrafficScheduling string `json:"traffic_scheduling,omitempty"`
}

// ShowONUDetail menampilkan detail ONU
// Command: show onu detail-info gpon-onu_{rack}/{shelf}/{slot}:{onu_id}
func (z *ZTEC320Client) ShowONUDetail(ctx context.Context, rack, shelf, slot, onuID int) (*ONUDetail, error) {
	cmd := fmt.Sprintf("show onu detail-info gpon-onu_%d/%d/%d:%d", rack, shelf, slot, onuID)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return z.parseONUDetail(output), nil
}

// ONUDistance informasi distance ONU
type ONUDistance struct {
	ONUID     string `json:"onu_id"`
	ONUType   string `json:"onu_type"`
	Distance  string `json:"distance"`
	Equalizer string `json:"equalizer_delay"`
}

// ShowGPONONUDistance menampilkan distance ONU
// Command: show gpon onu baseinfo gpon-olt_{rack}/{shelf}/{slot}
// Note: 'distance' command doesn't exist on some OLT models
func (z *ZTEC320Client) ShowGPONONUBaseInfo(ctx context.Context, rack, shelf, slot int) ([]ONUInfo, error) {
	cmd := fmt.Sprintf("show gpon onu baseinfo gpon-olt_%d/%d/%d", rack, shelf, slot)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return z.parseONUState(output), nil
}

// ONUTraffic informasi traffic ONU
type ONUTraffic struct {
	Interface string `json:"interface"`
	TxRate    string `json:"tx_rate"`
	RxRate    string `json:"rx_rate"`
	TxPkts    string `json:"tx_pkts"`
	RxPkts    string `json:"rx_pkts"`
	TxBytes   string `json:"tx_bytes"`
	RxBytes   string `json:"rx_bytes"`
}

// ShowONUTraffic menampilkan traffic ONU
// Command: show onu traffic gpon-onu_{rack}/{shelf}/{slot}:{onu_id}
func (z *ZTEC320Client) ShowONUTraffic(ctx context.Context, rack, shelf, slot, onuID int) (*ONUTraffic, error) {
	cmd := fmt.Sprintf("show onu traffic gpon-onu_%d/%d/%d:%d", rack, shelf, slot, onuID)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return z.parseONUTraffic(output), nil
}

// ONUOptical informasi optical ONU
type ONUOptical struct {
	ONUName      string `json:"onu_name"`
	ONUType      string `json:"onu_type"`
	ONUSN        string `json:"onu_sn"`
	OpticalPower string `json:"optical_power"`
	TxPower      string `json:"tx_power"`
	RxPower      string `json:"rx_power"`
	ONUTemp      string `json:"onu_temp"`
	Voltage      string `json:"voltage"`
	BiasCurrent  string `json:"bias_current"`
}

// ShowONUOptical menampilkan optical info ONU
// Command: show onu optical-info gpon-onu_{rack}/{shelf}/{slot}:{onu_id}
func (z *ZTEC320Client) ShowONUOptical(ctx context.Context, rack, shelf, slot, onuID int) (*ONUOptical, error) {
	cmd := fmt.Sprintf("show onu optical-info gpon-onu_%d/%d/%d:%d", rack, shelf, slot, onuID)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return z.parseONUOptical(output), nil
}

// ============================================================
// PRIORITY 2: HARDWARE DETAIL
// ============================================================

// ShowCardBySlot menampilkan card by slot
// Command: show card slotno {slot}
func (z *ZTEC320Client) ShowCardBySlot(ctx context.Context, slot int) (*CardInfo, error) {
	cmd := fmt.Sprintf("show card slotno %d", slot)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	cards := z.parseCardDetail(output)
	if len(cards) > 0 {
		return &cards[0], nil
	}
	return nil, fmt.Errorf("card not found")
}

// SubCardInfo informasi subcard
type SubCardInfo struct {
	Rack    int    `json:"rack"`
	Shelf   int    `json:"shelf"`
	Slot    int    `json:"slot"`
	SubSlot int    `json:"sub_slot"`
	Type    string `json:"type"`
	Status  string `json:"status"`
}

// ShowSubCard menampilkan subcard
// Command: show subcard
func (z *ZTEC320Client) ShowSubCard(ctx context.Context) ([]SubCardInfo, error) {
	output, err := z.client.Execute(ctx, "show subcard")
	if err != nil {
		return nil, err
	}
	return z.parseSubCard(output), nil
}

// ============================================================
// PRIORITY 3: GPON PROFILES
// ============================================================

// IPProfile informasi IP profile
type IPProfile struct {
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
	Mask      string `json:"mask"`
	Gateway   string `json:"gateway"`
}

// ShowIPProfile menampilkan IP profile
// Command: show gpon onu profile ip {name}
func (z *ZTEC320Client) ShowIPProfile(ctx context.Context, name string) (*IPProfile, error) {
	cmd := fmt.Sprintf("show gpon onu profile ip %s", name)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return z.parseIPProfile(output), nil
}

// SIPProfile informasi SIP profile
type SIPProfile struct {
	Name        string `json:"name"`
	ProxyServer string `json:"proxy_server"`
	Registrar   string `json:"registrar"`
	Outbound    string `json:"outbound"`
}

// ShowSIPProfile menampilkan SIP profile
// Command: show gpon onu profile sip {name}
func (z *ZTEC320Client) ShowSIPProfile(ctx context.Context, name string) (*SIPProfile, error) {
	cmd := fmt.Sprintf("show gpon onu profile sip %s", name)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return z.parseSIPProfile(output), nil
}

// MGCProfile informasi MGC profile
type MGCProfile struct {
	Name     string `json:"name"`
	MGC1IP   string `json:"mgc1_ip"`
	MGC1Port string `json:"mgc1_port"`
	MGC2IP   string `json:"mgc2_ip"`
	MGC2Port string `json:"mgc2_port"`
}

// ShowMGCProfile menampilkan MGC profile
// Command: show gpon onu profile mgc {name}
func (z *ZTEC320Client) ShowMGCProfile(ctx context.Context, name string) (*MGCProfile, error) {
	cmd := fmt.Sprintf("show gpon onu profile mgc %s", name)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return z.parseMGCProfile(output), nil
}

// ============================================================
// PRIORITY 4: LINE & REMOTE PROFILES
// ============================================================

// LineProfile informasi line profile
type LineProfile struct {
	Name       string   `json:"name"`
	TCONTList  []string `json:"tcont_list,omitempty"`
	GEMPortList []string `json:"gemport_list,omitempty"`
}

// ShowLineProfileList menampilkan list line profile
// Command: show pon onu-profile gpon line
func (z *ZTEC320Client) ShowLineProfileList(ctx context.Context) ([]string, error) {
	output, err := z.client.Execute(ctx, "show pon onu-profile gpon line")
	if err != nil {
		return nil, err
	}
	return z.parseProfileList(output), nil
}

// ShowLineProfile menampilkan detail line profile
// Command: show pon onu-profile gpon line {name}
func (z *ZTEC320Client) ShowLineProfile(ctx context.Context, name string) (*LineProfile, error) {
	cmd := fmt.Sprintf("show pon onu-profile gpon line %s", name)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return z.parseLineProfile(output), nil
}

// RemoteProfile informasi remote profile
type RemoteProfile struct {
	Name        string   `json:"name"`
	VLANList    []string `json:"vlan_list,omitempty"`
	ServiceList []string `json:"service_list,omitempty"`
}

// ShowRemoteProfileList menampilkan list remote profile
// Command: show pon onu-profile gpon remote
func (z *ZTEC320Client) ShowRemoteProfileList(ctx context.Context) ([]string, error) {
	output, err := z.client.Execute(ctx, "show pon onu-profile gpon remote")
	if err != nil {
		return nil, err
	}
	return z.parseProfileList(output), nil
}

// ShowRemoteProfile menampilkan detail remote profile
// Command: show pon onu-profile gpon remote {name}
func (z *ZTEC320Client) ShowRemoteProfile(ctx context.Context, name string) (*RemoteProfile, error) {
	cmd := fmt.Sprintf("show pon onu-profile gpon remote %s", name)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return z.parseRemoteProfile(output), nil
}

// ============================================================
// PRIORITY 5: VLAN
// ============================================================

// VLANInfo informasi VLAN
type VLANInfo struct {
	ID          int      `json:"id"`
	Name        string   `json:"name,omitempty"`
	Type        string   `json:"type"`
	Ports       []string `json:"ports,omitempty"`
	Description string   `json:"description,omitempty"`
}

// ShowVLANList menampilkan list VLAN
// Command: show vlan
func (z *ZTEC320Client) ShowVLANList(ctx context.Context) ([]VLANInfo, error) {
	output, err := z.client.Execute(ctx, "show vlan")
	if err != nil {
		return nil, err
	}
	return z.parseVLANList(output), nil
}

// ShowVLANByID menampilkan detail VLAN by ID
// Command: show vlan id {id}
func (z *ZTEC320Client) ShowVLANByID(ctx context.Context, id int) (*VLANInfo, error) {
	cmd := fmt.Sprintf("show vlan id %d", id)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	vlans := z.parseVLANList(output)
	if len(vlans) > 0 {
		return &vlans[0], nil
	}
	return nil, fmt.Errorf("vlan not found")
}

// ============================================================
// PRIORITY 6: IGMP/MULTICAST
// ============================================================

// IGMPMVlanInfo informasi MVLAN
type IGMPMVlanInfo struct {
	MVLANID   int      `json:"mvlan_id"`
	SourceIP  string   `json:"source_ip,omitempty"`
	WorkMode  string   `json:"work_mode,omitempty"`
	GroupList []string `json:"group_list,omitempty"`
}

// ShowIGMPMVlan menampilkan list MVLAN
// Command: show igmp mvlan
func (z *ZTEC320Client) ShowIGMPMVlan(ctx context.Context) ([]IGMPMVlanInfo, error) {
	output, err := z.client.Execute(ctx, "show igmp mvlan")
	if err != nil {
		return nil, err
	}
	return z.parseIGMPMVlan(output), nil
}

// ShowIGMPMVlanByID menampilkan detail MVLAN by ID
// Command: show igmp mvlan {id}
func (z *ZTEC320Client) ShowIGMPMVlanByID(ctx context.Context, id int) (*IGMPMVlanInfo, error) {
	cmd := fmt.Sprintf("show igmp mvlan %d", id)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	mvlan := z.parseIGMPMVlanDetail(output)
	if mvlan != nil {
		return mvlan, nil
	}
	return nil, fmt.Errorf("mvlan not found")
}

// ShowIGMPDynamicMember menampilkan IGMP dynamic membership
// Command: show igmp dynamic-member
func (z *ZTEC320Client) ShowIGMPDynamicMember(ctx context.Context) (string, error) {
	return z.client.Execute(ctx, "show igmp dynamic-member")
}

// ShowIGMPForwardingTable menampilkan IGMP forwarding table
// Command: show igmp forwarding-table
func (z *ZTEC320Client) ShowIGMPForwardingTable(ctx context.Context) (string, error) {
	return z.client.Execute(ctx, "show igmp forwarding-table")
}

// ShowIGMPInterface menampilkan IGMP interface config
// Command: show igmp interface
func (z *ZTEC320Client) ShowIGMPInterface(ctx context.Context) (string, error) {
	return z.client.Execute(ctx, "show igmp interface")
}

// ============================================================
// PRIORITY 7: INTERFACE DETAIL
// ============================================================

// InterfaceStats statistik interface
type InterfaceStats struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	TxRate    string `json:"tx_rate"`
	RxRate    string `json:"rx_rate"`
	TxPackets string `json:"tx_packets"`
	RxPackets string `json:"rx_packets"`
	TxBytes   string `json:"tx_bytes"`
	RxBytes   string `json:"rx_bytes"`
	TxErrors  string `json:"tx_errors"`
	RxErrors  string `json:"rx_errors"`
}

// ShowInterfaceByType menampilkan interface by type
// Command: show interface {interface_name}
func (z *ZTEC320Client) ShowInterfaceByType(ctx context.Context, ifName string) (*InterfaceStats, error) {
	cmd := fmt.Sprintf("show interface %s", ifName)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return z.parseInterfaceStats(output), nil
}

// ============================================================
// PRIORITY 8: USER MANAGEMENT
// ============================================================

// OnlineUser informasi user online
type OnlineUser struct {
	Username  string `json:"username"`
	IP        string `json:"ip"`
	LoginTime string `json:"login_time"`
	From      string `json:"from"`
}

// ShowOnlineUsers menampilkan user yang sedang online
// Command: show users
func (z *ZTEC320Client) ShowOnlineUsers(ctx context.Context) ([]OnlineUser, error) {
	output, err := z.client.Execute(ctx, "show users")
	if err != nil {
		return nil, err
	}
	return z.parseOnlineUsers(output), nil
}

// ============================================================
// REMAINING READ ENDPOINTS
// ============================================================

// DialPlanProfile informasi dial plan profile
type DialPlanProfile struct {
	Name    string   `json:"name"`
	Digits  []string `json:"digits,omitempty"`
	Timeout string   `json:"timeout,omitempty"`
}

// ShowDialPlanProfile menampilkan dial plan profile
// Command: show gpon onu profile dial-plan {name}
func (z *ZTEC320Client) ShowDialPlanProfile(ctx context.Context, name string) (*DialPlanProfile, error) {
	cmd := fmt.Sprintf("show gpon onu profile dial-plan %s", name)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return z.parseDialPlanProfile(output), nil
}

// VoipAccesscodeProfile informasi voip accesscode profile
type VoipAccesscodeProfile struct {
	Name       string   `json:"name"`
	Accesscode []string `json:"accesscode,omitempty"`
}

// ShowVoipAccesscodeProfile menampilkan voip accesscode profile
// Command: show gpon onu profile voip-accesscode {name}
func (z *ZTEC320Client) ShowVoipAccesscodeProfile(ctx context.Context, name string) (*VoipAccesscodeProfile, error) {
	cmd := fmt.Sprintf("show gpon onu profile voip-accesscode %s", name)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return z.parseVoipAccesscodeProfile(output), nil
}

// VoipAppsrvProfile informasi voip appsrv profile
type VoipAppsrvProfile struct {
	Name  string   `json:"name"`
	Apps  []string `json:"apps,omitempty"`
}

// ShowVoipAppsrvProfile menampilkan voip appsrv profile
// Command: show gpon onu profile voip-appsrv {name}
func (z *ZTEC320Client) ShowVoipAppsrvProfile(ctx context.Context, name string) (*VoipAppsrvProfile, error) {
	cmd := fmt.Sprintf("show gpon onu profile voip-appsrv %s", name)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return z.parseVoipAppsrvProfile(output), nil
}

// SNMPCommunity informasi SNMP community
type SNMPCommunity struct {
	Community string `json:"community"`
	Access    string `json:"access"`
}

// ShowSNMPCommunity menampilkan SNMP community
// Command: show snmp community
func (z *ZTEC320Client) ShowSNMPCommunity(ctx context.Context) ([]SNMPCommunity, error) {
	output, err := z.client.Execute(ctx, "show snmp community")
	if err != nil {
		return nil, err
	}
	return z.parseSNMPCommunity(output), nil
}

// SNMPHost informasi SNMP host
type SNMPHost struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	Community string `json:"community"`
	Version   string `json:"version"`
}

// ShowSNMPHost menampilkan SNMP host
// Command: show snmp host
func (z *ZTEC320Client) ShowSNMPHost(ctx context.Context) ([]SNMPHost, error) {
	output, err := z.client.Execute(ctx, "show snmp host")
	if err != nil {
		return nil, err
	}
	return z.parseSNMPHost(output), nil
}

// ShowRunningConfig menampilkan running config (already declared)
// Use existing method at line 232

// SaveConfig menyimpan konfigurasi (already declared)
// Use existing method at line 535

// BackupConfig backup konfigurasi ke TFTP
// Command: copy running-config tftp://{ip}/{filename}
func (z *ZTEC320Client) BackupConfig(ctx context.Context, tftpIP, filename string) (string, error) {
	cmd := fmt.Sprintf("copy running-config tftp://%s/%s", tftpIP, filename)
	return z.client.Execute(ctx, cmd)
}

// RestoreConfig restore konfigurasi dari TFTP
// Command: copy tftp://{ip}/{filename} running-config
func (z *ZTEC320Client) RestoreConfig(ctx context.Context, tftpIP, filename string) (string, error) {
	cmd := fmt.Sprintf("copy tftp://%s/%s running-config", tftpIP, filename)
	return z.client.Execute(ctx, cmd)
}

// ShowInterfaceVLAN menampilkan interface VLAN
// Command: show interface vlan{id}
func (z *ZTEC320Client) ShowInterfaceVLAN(ctx context.Context, vlanID int) (*InterfaceStats, error) {
	cmd := fmt.Sprintf("show interface vlan%d", vlanID)
	output, err := z.client.Execute(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return z.parseInterfaceStats(output), nil
}

// PowerSupplyInfo informasi power supply
type PowerSupplyInfo struct {
	Rack    int    `json:"rack"`
	Shelf   int    `json:"shelf"`
	Slot    int    `json:"slot"`
	Status  string `json:"status"`
	Voltage string `json:"voltage,omitempty"`
	Current string `json:"current,omitempty"`
}

// ShowPowerSupply menampilkan power supply info
// Command: show power
func (z *ZTEC320Client) ShowPowerSupply(ctx context.Context) ([]PowerSupplyInfo, error) {
	output, err := z.client.Execute(ctx, "show power")
	if err != nil {
		return nil, err
	}
	return z.parsePowerSupply(output), nil
}

// TemperatureInfo informasi temperature
type TemperatureInfo struct {
	Rack        int    `json:"rack"`
	Shelf       int    `json:"shelf"`
	Slot        int    `json:"slot"`
	Temperature string `json:"temperature"`
	Status      string `json:"status"`
}

// ShowTemperature menampilkan temperature info
// Command: show temperature
func (z *ZTEC320Client) ShowTemperature(ctx context.Context) ([]TemperatureInfo, error) {
	output, err := z.client.Execute(ctx, "show temperature")
	if err != nil {
		// Alternative: get from show fan
		output, err = z.client.Execute(ctx, "show fan")
		if err != nil {
			return nil, err
		}
	}
	return z.parseTemperature(output), nil
}

// ============================================================
// PARSER FUNCTIONS - PRIORITY 1
// ============================================================

func (z *ZTEC320Client) parseONUDetail(output string) *ONUDetail {
	detail := &ONUDetail{}
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "ONU Name") {
			detail.ONUName = z.extractValue(line)
		} else if strings.Contains(line, "ONU Type") {
			detail.ONUType = z.extractValue(line)
		} else if strings.Contains(line, "SN") && !strings.Contains(line, "Serial") {
			detail.ONUSN = z.extractValue(line)
		} else if strings.Contains(line, "Admin State") {
			detail.AdminState = z.extractValue(line)
		} else if strings.Contains(line, "Phase State") {
			detail.PhaseState = z.extractValue(line)
		} else if strings.Contains(line, "Channel") {
			detail.ChannelState = z.extractValue(line)
		} else if strings.Contains(line, "Authentication") {
			detail.Authentication = z.extractValue(line)
		} else if strings.Contains(line, "OMCC") {
			detail.Omcc = z.extractValue(line)
		} else if strings.Contains(line, "FEC Up") {
			detail.FecUp = z.extractValue(line)
		} else if strings.Contains(line, "FEC Down") {
			detail.FecDown = z.extractValue(line)
		} else if strings.Contains(line, "Line Profile") {
			detail.LineProfile = z.extractValue(line)
		} else if strings.Contains(line, "Remote Profile") {
			detail.RemoteProfile = z.extractValue(line)
		} else if strings.Contains(line, "Description") {
			detail.Description = z.extractValue(line)
		} else if strings.Contains(line, "Last Down Reason") {
			detail.LastDownReason = z.extractValue(line)
		} else if strings.Contains(line, "Last Down Time") {
			detail.LastDownTime = z.extractValue(line)
		}
	}
	
	return detail
}

func (z *ZTEC320Client) parseONUDistance(output string) []ONUDistance {
	var distances []ONUDistance
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "ONU") || strings.HasPrefix(line, "---") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			distances = append(distances, ONUDistance{
				ONUID:     fields[0],
				ONUType:   fields[1],
				Distance:  fields[2],
				Equalizer: fields[3],
			})
		}
	}
	
	return distances
}

func (z *ZTEC320Client) parseONUTraffic(output string) *ONUTraffic {
	traffic := &ONUTraffic{}
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "gpon-onu") {
			traffic.Interface = strings.Fields(line)[0]
		} else if strings.Contains(line, "Tx-rate") {
			traffic.TxRate = z.extractValue(line)
		} else if strings.Contains(line, "Rx-rate") {
			traffic.RxRate = z.extractValue(line)
		}
	}
	
	return traffic
}

func (z *ZTEC320Client) parseONUOptical(output string) *ONUOptical {
	optical := &ONUOptical{}
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "ONU Name") {
			optical.ONUName = z.extractValue(line)
		} else if strings.Contains(line, "ONU Type") {
			optical.ONUType = z.extractValue(line)
		} else if strings.Contains(line, "SN") && !strings.Contains(line, "Serial") {
			optical.ONUSN = z.extractValue(line)
		} else if strings.Contains(line, "Optical Power") || strings.Contains(line, "Tx Power") {
			optical.OpticalPower = z.extractValue(line)
		} else if strings.Contains(line, "Rx Power") {
			optical.RxPower = z.extractValue(line)
		} else if strings.Contains(line, "Temperature") {
			optical.ONUTemp = z.extractValue(line)
		} else if strings.Contains(line, "Voltage") {
			optical.Voltage = z.extractValue(line)
		} else if strings.Contains(line, "Bias Current") {
			optical.BiasCurrent = z.extractValue(line)
		}
	}
	
	return optical
}

// ============================================================
// PARSER FUNCTIONS - PRIORITY 2
// ============================================================

func (z *ZTEC320Client) parseSubCard(output string) []SubCardInfo {
	var subcards []SubCardInfo
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Rack") || strings.HasPrefix(line, "---") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 5 {
			rack, _ := strconv.Atoi(fields[0])
			shelf, _ := strconv.Atoi(fields[1])
			slot, _ := strconv.Atoi(fields[2])
			subslot, _ := strconv.Atoi(fields[3])
			subcards = append(subcards, SubCardInfo{
				Rack:    rack,
				Shelf:   shelf,
				Slot:    slot,
				SubSlot: subslot,
				Type:    fields[4],
				Status:  fields[5],
			})
		}
	}
	
	return subcards
}

// ============================================================
// PARSER FUNCTIONS - PRIORITY 3
// ============================================================

func (z *ZTEC320Client) parseIPProfile(output string) *IPProfile {
	profile := &IPProfile{}
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Profile Name") {
			profile.Name = z.extractValue(line)
		} else if strings.Contains(line, "IP Address") {
			profile.IPAddress = z.extractValue(line)
		} else if strings.Contains(line, "Mask") {
			profile.Mask = z.extractValue(line)
		} else if strings.Contains(line, "Gateway") {
			profile.Gateway = z.extractValue(line)
		}
	}
	
	return profile
}

func (z *ZTEC320Client) parseSIPProfile(output string) *SIPProfile {
	profile := &SIPProfile{}
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Profile Name") {
			profile.Name = z.extractValue(line)
		} else if strings.Contains(line, "Proxy Server") {
			profile.ProxyServer = z.extractValue(line)
		} else if strings.Contains(line, "Registrar") {
			profile.Registrar = z.extractValue(line)
		} else if strings.Contains(line, "Outbound") {
			profile.Outbound = z.extractValue(line)
		}
	}
	
	return profile
}

func (z *ZTEC320Client) parseMGCProfile(output string) *MGCProfile {
	profile := &MGCProfile{}
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Profile Name") {
			profile.Name = z.extractValue(line)
		} else if strings.Contains(line, "MGC1 IP") {
			profile.MGC1IP = z.extractValue(line)
		} else if strings.Contains(line, "MGC1 Port") {
			profile.MGC1Port = z.extractValue(line)
		} else if strings.Contains(line, "MGC2 IP") {
			profile.MGC2IP = z.extractValue(line)
		} else if strings.Contains(line, "MGC2 Port") {
			profile.MGC2Port = z.extractValue(line)
		}
	}
	
	return profile
}

// ============================================================
// PARSER FUNCTIONS - PRIORITY 4
// ============================================================

func (z *ZTEC320Client) parseProfileList(output string) []string {
	var profiles []string
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Profile") || strings.HasPrefix(line, "---") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			profiles = append(profiles, fields[0])
		}
	}
	
	return profiles
}

func (z *ZTEC320Client) parseLineProfile(output string) *LineProfile {
	profile := &LineProfile{}
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Profile Name") {
			profile.Name = z.extractValue(line)
		}
	}
	
	return profile
}

func (z *ZTEC320Client) parseRemoteProfile(output string) *RemoteProfile {
	profile := &RemoteProfile{}
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Profile Name") {
			profile.Name = z.extractValue(line)
		}
	}
	
	return profile
}

// ============================================================
// PARSER FUNCTIONS - PRIORITY 5
// ============================================================

func (z *ZTEC320Client) parseVLANList(output string) []VLANInfo {
	var vlans []VLANInfo
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "VLAN") || strings.HasPrefix(line, "---") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			id, err := strconv.Atoi(fields[0])
			if err == nil {
				vlans = append(vlans, VLANInfo{
					ID:   id,
					Type: fields[1],
				})
			}
		}
	}
	
	return vlans
}

// ============================================================
// PARSER FUNCTIONS - PRIORITY 6
// ============================================================

func (z *ZTEC320Client) parseIGMPMVlan(output string) []IGMPMVlanInfo {
	var mvlan []IGMPMVlanInfo
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "MVLAN") || strings.HasPrefix(line, "---") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 1 {
			id, err := strconv.Atoi(fields[0])
			if err == nil {
				mvlan = append(mvlan, IGMPMVlanInfo{
					MVLANID: id,
				})
			}
		}
	}
	
	return mvlan
}

func (z *ZTEC320Client) parseIGMPMVlanDetail(output string) *IGMPMVlanInfo {
	mvlan := &IGMPMVlanInfo{}
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "MVLAN ID") {
			mvlan.MVLANID, _ = strconv.Atoi(z.extractValue(line))
		} else if strings.Contains(line, "Source IP") {
			mvlan.SourceIP = z.extractValue(line)
		} else if strings.Contains(line, "Work Mode") {
			mvlan.WorkMode = z.extractValue(line)
		}
	}
	
	if mvlan.MVLANID > 0 {
		return mvlan
	}
	return nil
}

// ============================================================
// PARSER FUNCTIONS - PRIORITY 7
// ============================================================

func (z *ZTEC320Client) parseInterfaceStats(output string) *InterfaceStats {
	stats := &InterfaceStats{}
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Interface") && !strings.Contains(line, "Input") {
			stats.Name = strings.Fields(line)[0]
		} else if strings.Contains(line, "Status") {
			stats.Status = z.extractValue(line)
		}
	}
	
	return stats
}

// ============================================================
// PARSER FUNCTIONS - PRIORITY 8
// ============================================================

func (z *ZTEC320Client) parseOnlineUsers(output string) []OnlineUser {
	var users []OnlineUser
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Username") || strings.HasPrefix(line, "---") || strings.HasPrefix(line, "Line") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			users = append(users, OnlineUser{
				Username:  fields[0],
				IP:        fields[1],
				LoginTime: fields[2],
			})
		}
	}
	
	return users
}

// ============================================================
// UTILITY FUNCTIONS
// ============================================================

// ValidateSN memvalidasi format Serial Number ONU
func ValidateSN(sn string) bool {
	// ZTE SN format: ZTEG00000002 (4 letters + 8 hex)
	if len(sn) != 12 {
		return false
	}
	
	// Check first 4 characters are uppercase letters
	for i := 0; i < 4; i++ {
		c := sn[i]
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	
	// Check last 8 characters are hex digits
	for i := 4; i < 12; i++ {
		c := sn[i]
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	
	return true
}

func (z *ZTEC320Client) parseDialPlanProfile(output string) *DialPlanProfile {
	profile := &DialPlanProfile{}
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Profile Name") {
			profile.Name = z.extractValue(line)
		}
	}
	
	return profile
}

func (z *ZTEC320Client) parseVoipAccesscodeProfile(output string) *VoipAccesscodeProfile {
	profile := &VoipAccesscodeProfile{}
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Profile Name") {
			profile.Name = z.extractValue(line)
		}
	}
	
	return profile
}

func (z *ZTEC320Client) parseVoipAppsrvProfile(output string) *VoipAppsrvProfile {
	profile := &VoipAppsrvProfile{}
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Profile Name") {
			profile.Name = z.extractValue(line)
		}
	}
	
	return profile
}

func (z *ZTEC320Client) parseSNMPCommunity(output string) []SNMPCommunity {
	var communities []SNMPCommunity
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Community") || strings.HasPrefix(line, "---") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			communities = append(communities, SNMPCommunity{
				Community: fields[0],
				Access:    fields[1],
			})
		}
	}
	
	return communities
}

func (z *ZTEC320Client) parseSNMPHost(output string) []SNMPHost {
	var hosts []SNMPHost
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Host") || strings.HasPrefix(line, "---") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 1 {
			hosts = append(hosts, SNMPHost{
				Host: fields[0],
			})
		}
	}
	
	return hosts
}

func (z *ZTEC320Client) parsePowerSupply(output string) []PowerSupplyInfo {
	var supplies []PowerSupplyInfo
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Rack") || strings.HasPrefix(line, "---") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			rack, _ := strconv.Atoi(fields[0])
			shelf, _ := strconv.Atoi(fields[1])
			slot, _ := strconv.Atoi(fields[2])
			supplies = append(supplies, PowerSupplyInfo{
				Rack:   rack,
				Shelf:  shelf,
				Slot:   slot,
				Status: fields[3],
			})
		}
	}
	
	return supplies
}

func (z *ZTEC320Client) parseTemperature(output string) []TemperatureInfo {
	var temps []TemperatureInfo
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Temperature") || strings.Contains(line, "Temp") {
			// Extract temperature value
			if strings.Contains(line, ":") {
				value := z.extractValue(line)
				temps = append(temps, TemperatureInfo{
					Temperature: value,
					Status:      "normal",
				})
			}
		}
	}
	
	return temps
}
