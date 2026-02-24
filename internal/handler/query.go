package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ardani/snmp-zte/internal/driver"
	"github.com/ardani/snmp-zte/internal/driver/c320"
	"github.com/ardani/snmp-zte/internal/model"
	"github.com/ardani/snmp-zte/internal/snmp"
	"github.com/ardani/snmp-zte/pkg/response"
)

// QueryHandler menangani query SNMP "stateless" (tanpa simpan data).
type QueryHandler struct {
	pool *snmp.Pool
}

// NewQueryHandler membuat handler query baru.
func NewQueryHandler() *QueryHandler {
	return &QueryHandler{
		pool: snmp.GetPool(),
	}
}

// QueryRequest merepresentasikan permintaan query stateless.
// Semua detail koneksi OLT harus dikirim di setiap request.
type QueryRequest struct {
	// Detail koneksi OLT (Data ini TIDAK disimpan oleh server)
	IP        string `json:"ip" example:"192.168.1.1"`
	Port      int    `json:"port" example:"161"`
	Community string `json:"community" example:"public"`
	Model     string `json:"model" example:"C320"` // C320, C300, C600

	// Parameter Query (Apa yang ingin ditanyakan ke OLT)
	// Enum: onu_list, onu_detail, empty_slots, system_info, board_info, all_boards, interface_stats, fan_info, temperature_info, onu_traffic
	Query string `json:"query" example:"onu_list" enums:"onu_list,onu_detail,empty_slots,system_info,board_info,all_boards,interface_stats,fan_info,temperature_info,onu_traffic"`
	Board int    `json:"board" example:"1"`
	Pon   int    `json:"pon" example:"1"`
	OnuID int    `json:"onu_id,omitempty" example:"1"`
	
	// Parameter Provisioning (untuk create/rename)
	Name string `json:"name,omitempty" example:"customer-john"`
}

// QueryResponse merepresentasikan respons query
type QueryResponse struct {
	Query     string      `json:"query"`
	Summary   string      `json:"summary,omitempty"`
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp"`
	Duration  string      `json:"duration"`
}

// OLTInfoResponse merepresentasikan detail lengkap informasi OLT.
type OLTInfoResponse struct {
	System   interface{} `json:"system"`   // Detail Sistem (Nama, Deskripsi, Uptime)
	Model    interface{} `json:"model"`    // Kapabilitas Model (ZTE C320, C300, dll)
	Duration string      `json:"duration"` // Waktu yang dibutuhkan untuk query
}

// Query godoc
// @Summary Stateless SNMP Query (Query Tanpa Kredensial)
// @Description Melakukan query SNMP ke OLT tanpa menyimpan data login.
// @Description List 'query' yang didukung:
// @Description - onu_list: Daftar semua ONU di Port PON tertentu
// @Description - onu_detail: Detail lengkap satu ONU (WAJIB isi onu_id)
// @Description - empty_slots: Cari ID ONU yang masih kosong/tersedia
// @Description - system_info: Informasi sistem OLT (Nama, Deskripsi, Uptime)
// @Description - board_info: Status kartu/board (CPU, Memori, Tipe)
// @Description - all_boards: Status semua kartu yang ada di OLT
// @Description - interface_stats: Statistik lalu lintas interface (semua port)
// @Description - fan_info: Informasi status fan/kipas
// @Description - temperature_info: Informasi suhu sistem dan CPU (°C)
// @Description - onu_traffic: Statistik traffic ONU (RX/TX bytes, WAJIB isi onu_id)
// @Description - onu_bandwidth: Bandwidth SLA per ONU (assured/max kbps, WAJIB isi onu_id)
// @Description - pon_port_stats: Statistik traffic per PON port
// @Description - onu_errors: Error counter per ONU (CRC, FEC, dropped, WAJIB isi onu_id)
// @Description - voltage_info: Informasi voltage/power supply OLT
// @Tags Query
// @Accept json
// @Produce json
// @Param request body QueryRequest true "Detail Query (IP, Community, Model, dan Jenis Query)"
// @Success 200 {object} response.Response{data=QueryResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure 504 {object} response.ErrorResponse
// @Router /api/v1/query [post]
func (h *QueryHandler) Query(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Validasi field yang wajib diisi
	if req.IP == "" {
		response.BadRequest(w, "IP is required")
		return
	}
	if req.Community == "" {
		response.BadRequest(w, "Community is required")
		return
	}
	if req.Query == "" {
		response.BadRequest(w, "Query is required")
		return
	}

	// Set defaults
	if req.Port == 0 {
		req.Port = 161
	}
	if req.Model == "" {
		req.Model = "C320"
	}

	// Berikan batas waktu query (timeout) agar sistem tidak gantung jika OLT lambat
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Ambil driver berdasarkan model yang diminta (misal: C320)
	drv, err := h.getDriver(req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Hubungkan ke OLT
	if err := drv.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Gagal terhubung ke OLT: "+err.Error())
		return
	}
	defer drv.Close()

	// Jalankan query sesuai permintaan
	var result interface{}
	switch req.Query {
	case "onu_list":
		result, err = drv.GetONUList(ctx, req.Board, req.Pon)
	case "onu_detail":
		if req.OnuID == 0 {
			response.BadRequest(w, "onu_id is required for onu_detail query")
			return
		}
		result, err = drv.GetONUDetail(ctx, req.Board, req.Pon, req.OnuID)
	case "empty_slots":
		result, err = drv.GetEmptySlots(ctx, req.Board, req.Pon)
	case "system_info":
		result, err = drv.GetSystemInfo(ctx)
	case "board_info":
		result, err = drv.GetBoardInfo(ctx, req.Board)
	case "all_boards":
		result, err = drv.GetAllBoards(ctx)
	case "onu_traffic":
		if req.OnuID == 0 {
			response.BadRequest(w, "onu_id required for onu_traffic query")
			return
		}
		result, err = drv.GetONUTraffic(ctx, req.Board, req.Pon, req.OnuID)
	case "interface_stats":
		result, err = drv.GetInterfaceStats(ctx)
	case "fan_info":
		result, err = drv.GetFanInfo(ctx)
	case "temperature_info":
		result, err = drv.GetTemperatureInfo(ctx)
	// Phase 2: Bandwidth & Performance
	case "onu_bandwidth":
		if req.OnuID == 0 {
			response.BadRequest(w, "onu_id is required for onu_bandwidth query")
			return
		}
		result, err = drv.GetONUBandwidth(ctx, req.Board, req.Pon, req.OnuID)
	case "pon_port_stats":
		result, err = drv.GetPonPortStats(ctx, req.Board, req.Pon)
	case "onu_errors":
		if req.OnuID == 0 {
			response.BadRequest(w, "onu_id is required for onu_errors query")
			return
		}
		result, err = drv.GetONUErrors(ctx, req.Board, req.Pon, req.OnuID)
	case "voltage_info":
		result, err = drv.GetVoltageInfo(ctx)
	// Phase 3: Provisioning (SNMP SET - requires write community)
	case "onu_create":
		if req.OnuID == 0 {
			response.BadRequest(w, "onu_id is required for onu_create")
			return
		}
		err = drv.CreateONU(ctx, req.Board, req.Pon, req.OnuID, req.Name)
		if err == nil {
			result = map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("ONU %d created on Board %d PON %d", req.OnuID, req.Board, req.Pon),
				"board":   req.Board,
				"pon":     req.Pon,
				"onu_id":  req.OnuID,
				"name":    req.Name,
			}
		}
	case "onu_delete":
		if req.OnuID == 0 {
			response.BadRequest(w, "onu_id is required for onu_delete")
			return
		}
		err = drv.DeleteONU(ctx, req.Board, req.Pon, req.OnuID)
		if err == nil {
			result = map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("ONU %d deleted from Board %d PON %d", req.OnuID, req.Board, req.Pon),
				"board":   req.Board,
				"pon":     req.Pon,
				"onu_id":  req.OnuID,
			}
		}
	case "onu_rename":
		if req.OnuID == 0 {
			response.BadRequest(w, "onu_id is required for onu_rename")
			return
		}
		if req.Name == "" {
			response.BadRequest(w, "name is required for onu_rename")
			return
		}
		err = drv.RenameONU(ctx, req.Board, req.Pon, req.OnuID, req.Name)
		if err == nil {
			result = map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("ONU %d renamed to '%s'", req.OnuID, req.Name),
				"board":   req.Board,
				"pon":     req.Pon,
				"onu_id":  req.OnuID,
				"name":    req.Name,
			}
		}
	case "onu_status":
		if req.OnuID == 0 {
			response.BadRequest(w, "onu_id is required for onu_status")
			return
		}
		var status int
		status, err = drv.GetONUStatus(ctx, req.Board, req.Pon, req.OnuID)
		if err == nil {
			statusStr := "offline"
			if status == 2 {
				statusStr = "online"
			}
			result = map[string]interface{}{
				"board":      req.Board,
				"pon":        req.Pon,
				"onu_id":     req.OnuID,
				"status":     status,
				"status_str": statusStr,
			}
		}
	// Phase 4: Statistics
	case "distance_info":
		if req.OnuID == 0 {
			response.BadRequest(w, "onu_id is required for distance_info")
			return
		}
		result, err = drv.GetDistance(ctx, req.Board, req.Pon, req.OnuID)
	// Phase 5: VLAN
	case "vlan_list":
		result, err = drv.GetVLANList(ctx)
	case "vlan_info":
		if req.OnuID == 0 {
			response.BadRequest(w, "onu_id (vlan_id) is required for vlan_info")
			return
		}
		result, err = drv.GetVLANInfo(ctx, req.OnuID) // OnuID used as VLAN ID
	default:
		response.BadRequest(w, "Unknown query: "+req.Query)
		return
	}

	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Query failed: "+err.Error())
		return
	}

	// Susun jawaban (response)
	resp := QueryResponse{
		Query:     req.Query,
		Data:      result,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Duration:  time.Since(start).String(),
	}

	// Tambahkan summary khusus untuk empty_slots (misal: "110/128")
	if req.Query == "empty_slots" {
		if slots, ok := result.([]model.ONUSlot); ok {
			total := 128 // Default ZTE C320
			if drv.GetModelInfo().MaxOnuPerPon > 0 {
				total = drv.GetModelInfo().MaxOnuPerPon
			}
			resp.Summary = fmt.Sprintf("%d/%d", len(slots), total)
		}
	}

	response.JSON(w, http.StatusOK, resp)
}

// OLTInfoRequest merepresentasikan permintaan info OLT
type OLTInfoRequest struct {
	IP        string `json:"ip" example:"192.168.1.1"`
	Port      int    `json:"port" example:"161"`
	Community string `json:"community" example:"public"`
	Model     string `json:"model" example:"C320"`
}

// OLTInfo godoc
// @Summary Ambil Detail Sistem & Model OLT
// @Description Mengambil informasi lengkap sistem OLT (Nama, Deskripsi, Uptime) serta informasi kapabilitas model perangkat (contoh: Maksimal ONU per PON).
// @Tags Query
// @Accept json
// @Produce json
// @Param request body OLTInfoRequest true "Detail koneksi OLT"
// @Success 200 {object} response.Response{data=OLTInfoResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/olt-info [post]
func (h *QueryHandler) OLTInfo(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req OLTInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if req.IP == "" {
		response.BadRequest(w, "IP is required")
		return
	}

	if req.Port == 0 {
		req.Port = 161
	}
	if req.Community == "" {
		req.Community = "public"
	}
	if req.Model == "" {
		req.Model = "C320"
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	drv, err := h.getDriverFromOLTInfo(req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := drv.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Failed to connect to OLT")
		return
	}
	defer drv.Close()

	info, err := drv.GetSystemInfo(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to get OLT info")
		return
	}

	modelInfo := drv.GetModelInfo()

	response.JSON(w, http.StatusOK, OLTInfoResponse{
		System:   info,
		Model:    modelInfo,
		Duration: time.Since(start).String(),
	})
}

// PoolStats godoc
// @Summary Get SNMP Pool Stats
// @Description Get connection pool statistics
// @Tags System
// @Produce json
// @Success 200 {object} response.Response
// @Router /stats [get]
func (h *QueryHandler) PoolStats(w http.ResponseWriter, r *http.Request) {
	stats := h.pool.Stats()
	response.JSON(w, http.StatusOK, stats)
}

func (h *QueryHandler) getDriver(req QueryRequest) (driver.Driver, error) {
	switch req.Model {
	case "C320", "c320":
		return c320.New(req.IP, uint16(req.Port), req.Community), nil
	default:
		return nil, fmt.Errorf("unsupported OLT model: %s (supported: C320)", req.Model)
	}
}

func (h *QueryHandler) getDriverFromOLTInfo(req OLTInfoRequest) (driver.Driver, error) {
	switch req.Model {
	case "C320", "c320":
		return c320.New(req.IP, uint16(req.Port), req.Community), nil
	default:
		return nil, fmt.Errorf("unsupported OLT model: %s", req.Model)
	}
}
