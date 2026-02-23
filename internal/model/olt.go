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
