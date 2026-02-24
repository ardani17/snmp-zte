package model

// OLT merepresentasikan perangkat OLT
type OLT struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Model       string `json:"model"`
	IPAddress   string `json:"ip_address"`
	Port        int    `json:"port"`
	Community   string `json:"community"`
	BoardCount  int    `json:"board_count"`
	PonPerBoard int    `json:"pon_per_board"`
}

// OLTSummary merepresentasikan informasi ringkasan tentang OLT
type OLTSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Model       string `json:"model"`
	IPAddress   string `json:"ip_address"`
	BoardCount  int    `json:"board_count"`
	PonPerBoard int    `json:"pon_per_board"`
	TotalONU    int    `json:"total_onu"`
	OnlineONU   int    `json:"online_onu"`
	OfflineONU  int    `json:"offline_onu"`
}

// BoardInfo merepresentasikan informasi kartu/slot
type BoardInfo struct {
	BoardID   int    `json:"board_id"`
	Type      string `json:"type"`
	RealType  string `json:"real_type"`
	Status    string `json:"status"`
	PortCount int    `json:"port_count"`
	CpuLoad   int    `json:"cpu_load"`
	MemUsage  int    `json:"mem_usage"`
	SoftVer   string `json:"soft_ver"`
}

// PONInfo merepresentasikan informasi port PON
type PONInfo struct {
	BoardID  int     `json:"board_id"`
	PonID    int     `json:"pon_id"`
	Status   string  `json:"status"`
	TxPower  float64 `json:"tx_power"`
	RxPower  float64 `json:"rx_power"`
	ONUCount int     `json:"onu_count"`
	RxBytes  int64   `json:"rx_bytes"`
	TxBytes  int64   `json:"tx_bytes"`
}

// ONUTraffic merepresentasikan statistik trafik ONU
type ONUTraffic struct {
	Board     int    `json:"board"`
	PON       int    `json:"pon"`
	ONUID     int    `json:"onu_id"`
	RxBytes   int64  `json:"rx_bytes"`
	TxBytes   int64  `json:"tx_bytes"`
	RxPackets int64  `json:"rx_packets"`
	TxPackets int64  `json:"tx_packets"`
	Timestamp string `json:"timestamp"`
}

// InterfaceStats merepresentasikan statistik interface
type InterfaceStats struct {
	Index       int    `json:"index"`
	Description string `json:"description"`
	Status      string `json:"status"`
	RxBytes     int64  `json:"rx_bytes"`
	TxBytes     int64  `json:"tx_bytes"`
}

// TemperatureInfo merepresentasikan informasi sensor suhu
type TemperatureInfo struct {
	System    int `json:"system"`     // Ambient/System temperature (°C)
	CPU       int `json:"cpu"`        // CPU/Board temperature (°C)
	Timestamp string `json:"timestamp"`
}

// PONPortStats merepresentasikan statistik traffic per PON port
type PONPortStats struct {
	Board     int    `json:"board"`
	PON       int    `json:"pon"`
	RxBytes   int64  `json:"rx_bytes"`
	TxBytes   int64  `json:"tx_bytes"`
	RxPackets int64  `json:"rx_packets"`
	TxPackets int64  `json:"tx_packets"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// ONUErrors merepresentasikan error counter per ONU
type ONUErrors struct {
	Board          int    `json:"board"`
	PON            int    `json:"pon"`
	ONUID          int    `json:"onu_id"`
	CrcErrors      int64  `json:"crc_errors"`
	FecErrors      int64  `json:"fec_errors"`
	DroppedFrames  int64  `json:"dropped_frames"`
	LostPackets    int64  `json:"lost_packets"`
	Timestamp      string `json:"timestamp"`
}

// VoltageInfo merepresentasikan informasi voltage/power supply
type VoltageInfo struct {
	SystemVoltage int    `json:"system_voltage"` // mV
	CpuVoltage    int    `json:"cpu_voltage"`    // mV
	Timestamp     string `json:"timestamp"`
}

// CardStatus merepresentasikan kode status kartu
type CardStatus int

const (
	CardStatusInService CardStatus = 1 // Active/InService
	CardStatusFault     CardStatus = 2 // Fault/Error
	CardStatusOffline   CardStatus = 3 // Offline/Not Present
	CardStatusInit      CardStatus = 4 // Initializing
	CardStatusEmpty     CardStatus = 5 // Empty/No Card
)

func (s CardStatus) String() string {
	switch s {
	case CardStatusInService:
		return "InService"
	case CardStatusFault:
		return "Fault"
	case CardStatusOffline:
		return "Offline"
	case CardStatusInit:
		return "Initializing"
	case CardStatusEmpty:
		return "Empty"
	default:
		return "Unknown"
	}
}

// ONUDistance represents ONU distance information
type ONUDistance struct {
	Board    int    `json:"board"`
	PON      int    `json:"pon"`
	ONUID    int    `json:"onu_id"`
	Distance int    `json:"distance"` // Distance in meters
	EQD      int    `json:"eqd"`      // Equalized Delay
}

// VLANList represents list of VLANs
type VLANList struct {
	Count int        `json:"count"`
	VLANs []VLANInfo `json:"vlans"`
}

// VLANInfo represents VLAN information
type VLANInfo struct {
	VLANID int    `json:"vlan_id"`
	Name   string `json:"name"`
}

// ProfileList represents list of bandwidth profiles
type ProfileList struct {
	Count    int          `json:"count"`
	Profiles []ProfileInfo `json:"profiles"`
}

// ProfileInfo represents bandwidth profile information
type ProfileInfo struct {
	Index       int    `json:"index"`
	Name        string `json:"name"`
	FixedBW     int    `json:"fixed_bw"`
	AssuredBW   int    `json:"assured_bw"`
	MaxBW       int    `json:"max_bw"`
}
