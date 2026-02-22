package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ardani/snmp-zte/internal/driver"
	"github.com/ardani/snmp-zte/internal/driver/c320"
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
	Query string `json:"query" example:"onu_list"` // onu_list, onu_detail, dsb.
	Board int    `json:"board" example:"1"`
	Pon   int    `json:"pon" example:"1"`
	OnuID int    `json:"onu_id,omitempty" example:"1"`
}

// QueryResponse represents a query response
type QueryResponse struct {
	Query     string      `json:"query"`
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp"`
	Duration  string      `json:"duration"`
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
// @Description - pon_info: Statistik port PON (TX/RX Power)
// @Description - interface_stats: Statistik lalu lintas interface (semua port)
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

	// Validate required fields
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
	case "pon_info":
		result, err = drv.GetPONInfo(ctx, req.Board, req.Pon)
	case "interface_stats":
		result, err = drv.GetInterfaceStats(ctx)
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

	response.JSON(w, http.StatusOK, resp)
}

// OLTInfoRequest represents OLT info request
type OLTInfoRequest struct {
	IP        string `json:"ip" example:"192.168.1.1"`
	Port      int    `json:"port" example:"161"`
	Community string `json:"community" example:"public"`
	Model     string `json:"model" example:"C320"`
}

// OLTInfo godoc
// @Summary Get OLT Info
// @Description Get OLT system information
// @Tags Query
// @Accept json
// @Produce json
// @Param request body OLTInfoRequest true "OLT Info Request"
// @Success 200 {object} response.Response
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

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"system":   info,
		"model":    modelInfo,
		"duration": time.Since(start).String(),
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
