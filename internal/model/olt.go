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
	ID             string `json:"id"`
	Name           string `json:"name"`
	Model          string `json:"model"`
	IPAddress      string `json:"ip_address"`
	BoardCount     int    `json:"board_count"`
	PonPerBoard    int    `json:"pon_per_board"`
	TotalONU       int    `json:"total_onu"`
	OnlineONU      int    `json:"online_onu"`
	OfflineONU     int    `json:"offline_onu"`
}

// BoardInfo represents board/slot information
type BoardInfo struct {
	BoardID   int    `json:"board_id"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	CpuLoad   int    `json:"cpu_load"`
	MemUsage  int    `json:"mem_usage"`
}

// PONInfo represents PON port information
type PONInfo struct {
	BoardID     int     `json:"board_id"`
	PonID       int     `json:"pon_id"`
	Status      string  `json:"status"`
	Temperature float64 `json:"temperature"`
	TxPower     float64 `json:"tx_power"`
	RxPower     float64 `json:"rx_power"`
	ONUCount    int     `json:"onu_count"`
}
