package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ardani/snmp-zte/internal/cli"
	"github.com/ardani/snmp-zte/pkg/response"
)

// CLIHandler menangani CLI commands via SSH
type CLIHandler struct{}

// NewCLIHandler membuat handler CLI baru
func NewCLIHandler() *CLIHandler {
	return &CLIHandler{}
}

// CLIRequest permintaan CLI
type CLIRequest struct {
	// Koneksi
	Host     string `json:"host"`
	Port     int    `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`

	// Command
	Command string `json:"command,omitempty"`
	Query   string `json:"query,omitempty"`

	// Parameter untuk command spesifik
	Rack    int    `json:"rack,omitempty"`
	Shelf   int    `json:"shelf,omitempty"`
	Slot    int    `json:"slot,omitempty"`
	OnuID   int    `json:"onu_id,omitempty"`
	OnuType string `json:"onu_type,omitempty"`
	SN      string `json:"sn,omitempty"`
}

// CLIResponse response CLI
type CLIResponse struct {
	Query     string      `json:"query"`
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp"`
	Duration  string      `json:"duration"`
	Source    string      `json:"source"`
}

// getClient membuat client ZTE C320
func (h *CLIHandler) getClient(req CLIRequest) *cli.ZTEC320Client {
	cfg := cli.Config{
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
	}
	return cli.NewZTEC320Client(cfg)
}

// Execute godoc
// @Summary Execute CLI Command via SSH
// @Description Menjalankan command CLI ke OLT ZTE via SSH
// @Tags CLI
// @Accept json
// @Produce json
// @Param request body CLIRequest true "Detail koneksi dan command"
// @Success 200 {object} response.Response{data=CLIResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/cli [post]
func (h *CLIHandler) Execute(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Validasi
	if req.Host == "" {
		response.BadRequest(w, "host is required")
		return
	}
	if req.Command == "" && req.Query == "" {
		response.BadRequest(w, "command or query is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed: "+err.Error())
		return
	}
	defer client.Close()

	var result interface{}
	var err error
	var queryType string

	if req.Command != "" {
		// Raw command
		queryType = "raw"
		result, err = client.RawCommand(ctx, req.Command)
	} else {
		// Structured query
		queryType = req.Query
		switch req.Query {
		case "show_card":
			result, err = client.ShowCard(ctx)
		case "show_card_slot":
			if req.Slot == 0 {
				response.BadRequest(w, "slot is required for show_card_slot")
				return
			}
			result, err = client.ShowCardSlot(ctx, req.Slot)
		case "show_gpon_onu_uncfg":
			rack := req.Rack
			shelf := req.Shelf
			slot := req.Slot
			if rack == 0 {
				rack = 1
			}
			if shelf == 0 {
				shelf = 1
			}
			if slot == 0 {
				response.BadRequest(w, "slot is required for show_gpon_onu_uncfg")
				return
			}
			result, err = client.ShowGPONONUUnCfg(ctx, rack, shelf, slot)
		case "show_gpon_onu_state":
			rack := req.Rack
			shelf := req.Shelf
			slot := req.Slot
			if rack == 0 {
				rack = 1
			}
			if shelf == 0 {
				shelf = 1
			}
			if slot == 0 {
				response.BadRequest(w, "slot is required for show_gpon_onu_state")
				return
			}
			result, err = client.ShowGPONONUState(ctx, rack, shelf, slot)
		case "show_gpon_profile_tcont":
			result, err = client.ShowGPONProfileTcont(ctx)
		case "show_onu_type":
			if req.OnuType == "" {
				// List all
				result, err = client.ShowONUTypeList(ctx)
			} else {
				result, err = client.ShowONUType(ctx, req.OnuType)
			}
		case "show_fan":
			result, err = client.ShowFan(ctx)
		case "show_version":
			result, err = client.ShowVersion(ctx)
		case "show_clock":
			result, err = client.ShowClock(ctx)
		case "show_running_config":
			result, err = client.ShowRunningConfig(ctx)
		// Configuration commands
		case "onu_authenticate":
			if req.OnuID == 0 || req.OnuType == "" || req.SN == "" {
				response.BadRequest(w, "onu_id, onu_type, and sn are required for onu_authenticate")
				return
			}
			rack := req.Rack
			shelf := req.Shelf
			slot := req.Slot
			if rack == 0 {
				rack = 1
			}
			if shelf == 0 {
				shelf = 1
			}
			if slot == 0 {
				response.BadRequest(w, "slot is required for onu_authenticate")
				return
			}
			err = client.AuthenticateONU(ctx, rack, shelf, slot, req.OnuID, req.OnuType, req.SN)
			if err == nil {
				result = map[string]interface{}{
					"success": true,
					"message": fmt.Sprintf("ONU %d authenticated on slot %d", req.OnuID, slot),
					"onu_id":  req.OnuID,
					"type":    req.OnuType,
					"sn":      req.SN,
				}
			}
		case "onu_delete":
			if req.OnuID == 0 {
				response.BadRequest(w, "onu_id is required for onu_delete")
				return
			}
			rack := req.Rack
			shelf := req.Shelf
			slot := req.Slot
			if rack == 0 {
				rack = 1
			}
			if shelf == 0 {
				shelf = 1
			}
			if slot == 0 {
				response.BadRequest(w, "slot is required for onu_delete")
				return
			}
			err = client.DeleteONU(ctx, rack, shelf, slot, req.OnuID)
			if err == nil {
				result = map[string]interface{}{
					"success": true,
					"message": fmt.Sprintf("ONU %d deleted from slot %d", req.OnuID, slot),
					"onu_id":  req.OnuID,
				}
			}
		case "save_config":
			err = client.SaveConfig(ctx)
			if err == nil {
				result = map[string]interface{}{
					"success": true,
					"message": "Configuration saved",
				}
			}
		default:
			response.BadRequest(w, "Unknown query: "+req.Query)
			return
		}
	}

	if err != nil {
		response.Error(w, http.StatusInternalServerError, "CLI command failed: "+err.Error())
		return
	}

	resp := CLIResponse{
		Query:     queryType,
		Data:      result,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Duration:  time.Since(start).String(),
		Source:    "cli-ssh",
	}

	response.JSON(w, http.StatusOK, resp)
}

// ShowCard godoc
// @Summary Show Card Status
// @Description Menampilkan status semua card di OLT
// @Tags CLI
// @Accept json
// @Produce json
// @Param request body CLIRequest true "Detail koneksi"
// @Success 200 {object} response.Response{data=CLIResponse}
// @Router /api/v1/cli/card [post]
func (h *CLIHandler) ShowCard(w http.ResponseWriter, r *http.Request) {
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if req.Host == "" {
		response.BadRequest(w, "host is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	start := time.Now()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed: "+err.Error())
		return
	}
	defer client.Close()

	result, err := client.ShowCard(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Command failed: "+err.Error())
		return
	}

	response.JSON(w, http.StatusOK, CLIResponse{
		Query:     "show_card",
		Data:      result,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Duration:  time.Since(start).String(),
		Source:    "cli-ssh",
	})
}

// ShowONUState godoc
// @Summary Show ONU State
// @Description Menampilkan state ONU di PON port tertentu
// @Tags CLI
// @Accept json
// @Produce json
// @Param request body CLIRequest true "Detail koneksi dan slot"
// @Success 200 {object} response.Response{data=CLIResponse}
// @Router /api/v1/cli/onu/state [post]
func (h *CLIHandler) ShowONUState(w http.ResponseWriter, r *http.Request) {
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if req.Host == "" {
		response.BadRequest(w, "host is required")
		return
	}
	if req.Slot == 0 {
		response.BadRequest(w, "slot is required")
		return
	}

	rack := req.Rack
	shelf := req.Shelf
	if rack == 0 {
		rack = 1
	}
	if shelf == 0 {
		shelf = 1
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	start := time.Now()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed: "+err.Error())
		return
	}
	defer client.Close()

	result, err := client.ShowGPONONUState(ctx, rack, shelf, req.Slot)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Command failed: "+err.Error())
		return
	}

	response.JSON(w, http.StatusOK, CLIResponse{
		Query:     "show_gpon_onu_state",
		Data:      result,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Duration:  time.Since(start).String(),
		Source:    "cli-ssh",
	})
}

// ShowONUUncfg godoc
// @Summary Show Unconfigured ONUs
// @Description Menampilkan ONU yang belum terdaftar
// @Tags CLI
// @Accept json
// @Produce json
// @Param request body CLIRequest true "Detail koneksi dan slot"
// @Success 200 {object} response.Response{data=CLIResponse}
// @Router /api/v1/cli/onu/uncfg [post]
func (h *CLIHandler) ShowONUUncfg(w http.ResponseWriter, r *http.Request) {
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if req.Host == "" {
		response.BadRequest(w, "host is required")
		return
	}
	if req.Slot == 0 {
		response.BadRequest(w, "slot is required")
		return
	}

	rack := req.Rack
	shelf := req.Shelf
	if rack == 0 {
		rack = 1
	}
	if shelf == 0 {
		shelf = 1
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	start := time.Now()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed: "+err.Error())
		return
	}
	defer client.Close()

	result, err := client.ShowGPONONUUnCfg(ctx, rack, shelf, req.Slot)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Command failed: "+err.Error())
		return
	}

	response.JSON(w, http.StatusOK, CLIResponse{
		Query:     "show_gpon_onu_uncfg",
		Data:      result,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Duration:  time.Since(start).String(),
		Source:    "cli-ssh",
	})
}

// AuthenticateONU godoc
// @Summary Authenticate/Register ONU
// @Description Mendaftarkan ONU baru ke OLT
// @Tags CLI
// @Accept json
// @Produce json
// @Param request body CLIRequest true "Detail koneksi, slot, ONU ID, type, dan SN"
// @Success 200 {object} response.Response{data=CLIResponse}
// @Router /api/v1/cli/onu/auth [post]
func (h *CLIHandler) AuthenticateONU(w http.ResponseWriter, r *http.Request) {
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if req.Host == "" {
		response.BadRequest(w, "host is required")
		return
	}
	if req.Slot == 0 {
		response.BadRequest(w, "slot is required")
		return
	}
	if req.OnuID == 0 {
		response.BadRequest(w, "onu_id is required")
		return
	}
	if req.OnuType == "" {
		response.BadRequest(w, "onu_type is required (e.g., ZTEG-F620)")
		return
	}
	if req.SN == "" {
		response.BadRequest(w, "sn is required")
		return
	}

	// Validate SN format
	if !cli.ValidateSN(req.SN) {
		response.BadRequest(w, "invalid SN format (expected: ZTEG00000002)")
		return
	}

	rack := req.Rack
	shelf := req.Shelf
	if rack == 0 {
		rack = 1
	}
	if shelf == 0 {
		shelf = 1
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	start := time.Now()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed: "+err.Error())
		return
	}
	defer client.Close()

	err := client.AuthenticateONU(ctx, rack, shelf, req.Slot, req.OnuID, req.OnuType, req.SN)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Authentication failed: "+err.Error())
		return
	}

	response.JSON(w, http.StatusOK, CLIResponse{
		Query: "onu_authenticate",
		Data: map[string]interface{}{
			"success":  true,
			"message":  fmt.Sprintf("ONU %d authenticated on gpon-olt_%d/%d/%d", req.OnuID, rack, shelf, req.Slot),
			"onu_id":   req.OnuID,
			"onu_type": req.OnuType,
			"sn":       req.SN,
			"port":     fmt.Sprintf("gpon-olt_%d/%d/%d", rack, shelf, req.Slot),
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Duration:  time.Since(start).String(),
		Source:    "cli-ssh",
	})
}

// DeleteONU godoc
// @Summary Delete ONU
// @Description Menghapus ONU dari OLT
// @Tags CLI
// @Accept json
// @Produce json
// @Param request body CLIRequest true "Detail koneksi, slot, dan ONU ID"
// @Success 200 {object} response.Response{data=CLIResponse}
// @Router /api/v1/cli/onu/delete [post]
func (h *CLIHandler) DeleteONU(w http.ResponseWriter, r *http.Request) {
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if req.Host == "" {
		response.BadRequest(w, "host is required")
		return
	}
	if req.Slot == 0 {
		response.BadRequest(w, "slot is required")
		return
	}
	if req.OnuID == 0 {
		response.BadRequest(w, "onu_id is required")
		return
	}

	rack := req.Rack
	shelf := req.Shelf
	if rack == 0 {
		rack = 1
	}
	if shelf == 0 {
		shelf = 1
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	start := time.Now()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed: "+err.Error())
		return
	}
	defer client.Close()

	err := client.DeleteONU(ctx, rack, shelf, req.Slot, req.OnuID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Delete failed: "+err.Error())
		return
	}

	response.JSON(w, http.StatusOK, CLIResponse{
		Query: "onu_delete",
		Data: map[string]interface{}{
			"success": true,
			"message": fmt.Sprintf("ONU %d deleted from gpon-olt_%d/%d/%d", req.OnuID, rack, shelf, req.Slot),
			"onu_id":  req.OnuID,
			"port":    fmt.Sprintf("gpon-olt_%d/%d/%d", rack, shelf, req.Slot),
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Duration:  time.Since(start).String(),
		Source:    "cli-ssh",
	})
}


