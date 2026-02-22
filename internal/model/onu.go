package model

// ONUInfo represents basic ONU information (for list view)
type ONUInfo struct {
	OLTID        string `json:"olt_id"`
	Board        int    `json:"board"`
	PON          int    `json:"pon"`
	ID           int    `json:"onu_id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	SerialNumber string `json:"serial_number"`
	RXPower      string `json:"rx_power"`
	TXPower      string `json:"tx_power"`
	Distance     string `json:"distance"`
	Status       string `json:"status"`
}

// ONUDetail represents detailed ONU information
type ONUDetail struct {
	ONUInfo
	TXPower              string `json:"tx_power"`
	IPAddress            string `json:"ip_address"`
	Description          string `json:"description"`
	LastOnline           string `json:"last_online"`
	LastOffline          string `json:"last_offline"`
	Uptime               string `json:"uptime"`
	LastDownTimeDuration string `json:"last_down_duration"`
	OfflineReason        string `json:"offline_reason"`
	Distance             string `json:"distance"`
}

// ONUSlot represents an available ONU slot
type ONUSlot struct {
	Board  int `json:"board"`
	PON    int `json:"pon"`
	ONUID  int `json:"onu_id"`
}

// ONUStatus represents ONU status codes
type ONUStatus int

const (
	StatusLogging      ONUStatus = 1
	StatusLOS          ONUStatus = 2
	StatusSynchronization ONUStatus = 3
	StatusOnline       ONUStatus = 4
	StatusDyingGasp    ONUStatus = 5
	StatusAuthFailed   ONUStatus = 6
	StatusOffline      ONUStatus = 7
)

func (s ONUStatus) String() string {
	switch s {
	case StatusLogging:
		return "Logging"
	case StatusLOS:
		return "LOS"
	case StatusSynchronization:
		return "Synchronization"
	case StatusOnline:
		return "Online"
	case StatusDyingGasp:
		return "Dying Gasp"
	case StatusAuthFailed:
		return "Auth Failed"
	case StatusOffline:
		return "Offline"
	default:
		return "Unknown"
	}
}

// OfflineReason represents offline reason codes
type OfflineReason int

const (
	ReasonUnknown      OfflineReason = 1
	ReasonLOS          OfflineReason = 2
	ReasonLOSi         OfflineReason = 3
	ReasonLOFi         OfflineReason = 4
	ReasonSfi          OfflineReason = 5
	ReasonLoai         OfflineReason = 6
	ReasonLoami        OfflineReason = 7
	ReasonAuthFail     OfflineReason = 8
	ReasonPowerOff     OfflineReason = 9
	ReasonDeactiveSucc OfflineReason = 10
	ReasonDeactiveFail OfflineReason = 11
	ReasonReboot       OfflineReason = 12
	ReasonShutdown     OfflineReason = 13
)

func (r OfflineReason) String() string {
	switch r {
	case ReasonUnknown:
		return "Unknown"
	case ReasonLOS:
		return "LOS"
	case ReasonLOSi:
		return "LOSi"
	case ReasonLOFi:
		return "LOFi"
	case ReasonSfi:
		return "sfi"
	case ReasonLoai:
		return "loai"
	case ReasonLoami:
		return "loami"
	case ReasonAuthFail:
		return "AuthFail"
	case ReasonPowerOff:
		return "PowerOff"
	case ReasonDeactiveSucc:
		return "DeactiveSucc"
	case ReasonDeactiveFail:
		return "DeactiveFail"
	case ReasonReboot:
		return "Reboot"
	case ReasonShutdown:
		return "Shutdown"
	default:
		return "Unknown"
	}
}
