package model

// OLT represents an OLT device
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

// OLTSummary represents summary information about an OLT
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

// BoardInfo represents board/slot information
type BoardInfo struct {
	BoardID  int    `json:"board_id"`
	Type     string `json:"type"`
	Status   string `json:"status"`
	CpuLoad  int    `json:"cpu_load"`
	MemUsage int    `json:"mem_usage"`
}

// PONInfo represents PON port information
type PONInfo struct {
	BoardID  int     `json:"board_id"`
	PonID    int     `json:"pon_id"`
	Status   string  `json:"status"`
	TxPower  float64 `json:"tx_power"`
	RxPower  float64 `json:"rx_power"`
	ONUCount int     `json:"onu_count"`
}

// ONUTraffic represents ONU traffic statistics
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

// InterfaceStats represents interface statistics
type InterfaceStats struct {
	Index       int    `json:"index"`
	Description string `json:"description"`
	Status      string `json:"status"`
	RxBytes     int64  `json:"rx_bytes"`
	TxBytes     int64  `json:"tx_bytes"`
}

// CardStatus represents card status codes
type CardStatus int

const (
	CardStatusEmpty CardStatus = 1
	CardStatusReset CardStatus = 2
	CardStatusInit  CardStatus = 3
	CardStatusReady CardStatus = 4
	CardStatusFault CardStatus = 5
)

func (s CardStatus) String() string {
	switch s {
	case CardStatusEmpty:
		return "Empty"
	case CardStatusReset:
		return "Reset"
	case CardStatusInit:
		return "Initializing"
	case CardStatusReady:
		return "Ready"
	case CardStatusFault:
		return "Fault"
	default:
		return "Unknown"
	}
}
