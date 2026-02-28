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

// CLIHandler menangani CLI commands via Telnet
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

	// Parameter
	Rack    int    `json:"rack,omitempty"`
	Shelf   int    `json:"shelf,omitempty"`
	Slot    int    `json:"slot,omitempty"`
	OnuID   int    `json:"onu_id,omitempty"`
	VlanID  int    `json:"vlan_id,omitempty"`
	OnuType string `json:"onu_type,omitempty"`
	SN      string `json:"sn,omitempty"`
	Name    string `json:"name,omitempty"`
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

// respond helper
func (h *CLIHandler) respond(w http.ResponseWriter, query string, data interface{}, start time.Time) {
	response.JSON(w, http.StatusOK, CLIResponse{
		Query:     query,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Duration:  time.Since(start).String(),
		Source:    "cli-telnet",
	})
}

// ============================================================
// SYSTEM ENDPOINTS
// ============================================================

// ShowClock godoc
// @Summary Show System Clock
// @Tags CLI-System
// @Accept json
// @Produce json
// @Param request body CLIRequest true "Connection"
// @Success 200 {object} response.Response
// @Router /api/v1/cli/system/clock [post]
func (h *CLIHandler) ShowClock(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	output, err := client.Execute(ctx, "show clock")
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_clock", output, start)
}

// ============================================================
// HARDWARE ENDPOINTS
// ============================================================

// ShowCard godoc
// @Summary Show Card Status
// @Tags CLI-Hardware
// @Accept json
// @Produce json
// @Param request body CLIRequest true "Connection"
// @Success 200 {object} response.Response
// @Router /api/v1/cli/card [post]
func (h *CLIHandler) ShowCard(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	result, err := client.ShowCard(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_card", result, start)
}

// ShowRack godoc
// @Summary Show Rack Info
// @Tags CLI-Hardware
// @Router /api/v1/cli/rack [post]
func (h *CLIHandler) ShowRack(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	output, err := client.Execute(ctx, "show rack")
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_rack", output, start)
}

// ShowShelf godoc
// @Summary Show Shelf Info
// @Tags CLI-Hardware
// @Router /api/v1/cli/shelf [post]
func (h *CLIHandler) ShowShelf(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	output, err := client.Execute(ctx, "show shelf")
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_shelf", output, start)
}

// ShowFan godoc
// @Summary Show Fan Status
// @Tags CLI-Hardware
// @Router /api/v1/cli/fan [post]
func (h *CLIHandler) ShowFan(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	output, err := client.Execute(ctx, "show fan")
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_fan", output, start)
}

// ============================================================
// GPON PROFILE ENDPOINTS
// ============================================================

// ShowTcontProfile godoc
// @Summary Show T-CONT Profiles
// @Tags CLI-GPON
// @Router /api/v1/cli/gpon/tcont [post]
func (h *CLIHandler) ShowTcontProfile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	output, err := client.Execute(ctx, "show gpon profile tcont")
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_gpon_profile_tcont", output, start)
}

// ShowOnuType godoc
// @Summary Show ONU Types
// @Tags CLI-GPON
// @Router /api/v1/cli/gpon/onu-type [post]
func (h *CLIHandler) ShowOnuType(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	output, err := client.Execute(ctx, "show onu-type gpon")
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_onu_type_gpon", output, start)
}

// ShowVlanProfile godoc
// @Summary Show VLAN Profiles
// @Tags CLI-GPON
// @Router /api/v1/cli/gpon/vlan-profile [post]
func (h *CLIHandler) ShowVlanProfile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	output, err := client.Execute(ctx, "show gpon onu profile vlan")
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_gpon_onu_profile_vlan", output, start)
}

// ============================================================
// GPON ONU ENDPOINTS
// ============================================================

// ShowONUState godoc
// @Summary Show ONU State
// @Tags CLI-ONU
// @Router /api/v1/cli/onu/state [post]
func (h *CLIHandler) ShowONUState(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
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

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	result, err := client.ShowGPONONUState(ctx, rack, shelf, req.Slot)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_gpon_onu_state", result, start)
}

// ShowONUUncfg godoc
// @Summary Show Unconfigured ONUs
// @Tags CLI-ONU
// @Router /api/v1/cli/onu/uncfg [post]
func (h *CLIHandler) ShowONUUncfg(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
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

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	result, err := client.ShowGPONONUUnCfg(ctx, rack, shelf, req.Slot)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_gpon_onu_uncfg", result, start)
}

// ShowONUConfig godoc
// @Summary Show ONU Config
// @Tags CLI-ONU
// @Router /api/v1/cli/onu/config [post]
func (h *CLIHandler) ShowONUConfig(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Slot == 0 || req.OnuID == 0 {
		response.BadRequest(w, "slot and onu_id are required")
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

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	cmd := fmt.Sprintf("show running-config interface gpon-onu_%d/%d/%d:%d", rack, shelf, req.Slot, req.OnuID)
	output, err := client.Execute(ctx, cmd)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_running-config_interface", output, start)
}

// ShowONURunning godoc
// @Summary Show ONU Running Config
// @Tags CLI-ONU
// @Router /api/v1/cli/onu/running [post]
func (h *CLIHandler) ShowONURunning(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Slot == 0 || req.OnuID == 0 {
		response.BadRequest(w, "slot and onu_id are required")
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

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	cmd := fmt.Sprintf("show onu running config gpon-onu_%d/%d/%d:%d", rack, shelf, req.Slot, req.OnuID)
	output, err := client.Execute(ctx, cmd)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_onu_running_config", output, start)
}

// ============================================================
// INTERFACE ENDPOINTS
// ============================================================

// ShowInterface godoc
// @Summary Show Interface
// @Tags CLI-Interface
// @Router /api/v1/cli/interface [post]
func (h *CLIHandler) ShowInterface(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Name == "" {
		response.BadRequest(w, "name is required (e.g., gpon-olt_1/1/1, gei_1/4/1)")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	cmd := fmt.Sprintf("show interface %s", req.Name)
	output, err := client.Execute(ctx, cmd)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_interface", output, start)
}

// ShowMgmtInterface godoc
// @Summary Show Management Interface
// @Tags CLI-Interface
// @Router /api/v1/cli/interface/mng [post]
func (h *CLIHandler) ShowMgmtInterface(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	output, err := client.Execute(ctx, "show interface mng1")
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_interface_mng1", output, start)
}

// ============================================================
// SERVICE PORT ENDPOINTS
// ============================================================

// ShowServicePort godoc
// @Summary Show Service Port
// @Tags CLI-Service
// @Router /api/v1/cli/service-port [post]
func (h *CLIHandler) ShowServicePort(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Slot == 0 || req.OnuID == 0 {
		response.BadRequest(w, "slot and onu_id are required")
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

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	cmd := fmt.Sprintf("show service-port interface gpon-onu_%d/%d/%d:%d", rack, shelf, req.Slot, req.OnuID)
	output, err := client.Execute(ctx, cmd)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_service-port", output, start)
}

// ============================================================
// IGMP ENDPOINTS
// ============================================================

// ShowIGMP godoc
// @Summary Show IGMP Status
// @Tags CLI-IGMP
// @Router /api/v1/cli/igmp [post]
func (h *CLIHandler) ShowIGMP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	output, err := client.Execute(ctx, "show igmp")
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_igmp", output, start)
}

// ============================================================
// USER MANAGEMENT ENDPOINTS
// ============================================================

// ShowUsers godoc
// @Summary Show Users
// @Tags CLI-User
// @Router /api/v1/cli/user/list [post]
func (h *CLIHandler) ShowUsers(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	output, err := client.Execute(ctx, "show username")
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_username", output, start)
}

// ============================================================
// PRIORITY 1: ONU DETAIL ENDPOINTS
// ============================================================

// ShowONUDetail godoc
// @Summary Show ONU Detail Info
// @Tags CLI-ONU
// @Router /api/v1/cli/onu/detail [post]
func (h *CLIHandler) ShowONUDetail(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Slot == 0 || req.OnuID == 0 {
		response.BadRequest(w, "slot and onu_id are required")
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

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowONUDetail(ctx, rack, shelf, req.Slot, req.OnuID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_onu_detail_info", data, start)
}

// ShowONUBaseInfo godoc
// @Summary Show ONU Base Info (replaces distance)
// @Tags CLI-ONU
// @Router /api/v1/cli/onu/baseinfo [post]
func (h *CLIHandler) ShowONUBaseInfo(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
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

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowGPONONUBaseInfo(ctx, rack, shelf, req.Slot)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_gpon_onu_baseinfo", data, start)
}

// ShowONUTraffic godoc
// @Summary Show ONU Traffic
// @Tags CLI-ONU
// @Router /api/v1/cli/onu/traffic [post]
func (h *CLIHandler) ShowONUTraffic(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Slot == 0 || req.OnuID == 0 {
		response.BadRequest(w, "slot and onu_id are required")
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

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowONUTraffic(ctx, rack, shelf, req.Slot, req.OnuID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_onu_traffic", data, start)
}

// ShowONUOptical godoc
// @Summary Show ONU Optical Info
// @Tags CLI-ONU
// @Router /api/v1/cli/onu/optical [post]
func (h *CLIHandler) ShowONUOptical(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Slot == 0 || req.OnuID == 0 {
		response.BadRequest(w, "slot and onu_id are required")
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

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowONUOptical(ctx, rack, shelf, req.Slot, req.OnuID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_onu_optical_info", data, start)
}

// ============================================================
// PRIORITY 2: HARDWARE DETAIL
// ============================================================

// ShowCardBySlot godoc
// @Summary Show Card by Slot
// @Tags CLI-Hardware
// @Router /api/v1/cli/card/{slot} [post]
func (h *CLIHandler) ShowCardBySlot(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Slot == 0 {
		response.BadRequest(w, "slot is required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowCardBySlot(ctx, req.Slot)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_card_slotno", data, start)
}

// ShowSubCard godoc
// @Summary Show SubCard
// @Tags CLI-Hardware
// @Router /api/v1/cli/subcard [post]
func (h *CLIHandler) ShowSubCard(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowSubCard(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_subcard", data, start)
}

// ============================================================
// PRIORITY 3: GPON PROFILES
// ============================================================

// ShowIPProfile godoc
// @Summary Show IP Profile
// @Tags CLI-GPON
// @Router /api/v1/cli/gpon/ip-profile [post]
func (h *CLIHandler) ShowIPProfile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Name == "" {
		response.BadRequest(w, "name is required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowIPProfile(ctx, req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_gpon_onu_profile_ip", data, start)
}

// ShowSIPProfile godoc
// @Summary Show SIP Profile
// @Tags CLI-GPON
// @Router /api/v1/cli/gpon/sip-profile [post]
func (h *CLIHandler) ShowSIPProfile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Name == "" {
		response.BadRequest(w, "name is required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowSIPProfile(ctx, req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_gpon_onu_profile_sip", data, start)
}

// ShowMGCProfile godoc
// @Summary Show MGC Profile
// @Tags CLI-GPON
// @Router /api/v1/cli/gpon/mgc-profile [post]
func (h *CLIHandler) ShowMGCProfile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Name == "" {
		response.BadRequest(w, "name is required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowMGCProfile(ctx, req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_gpon_onu_profile_mgc", data, start)
}

// ============================================================
// PRIORITY 4: LINE & REMOTE PROFILES
// ============================================================

// ShowLineProfileList godoc
// @Summary Show Line Profile List
// @Tags CLI-Profile
// @Router /api/v1/cli/profile/line/list [post]
func (h *CLIHandler) ShowLineProfileList(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowLineProfileList(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_pon_onu-profile_gpon_line", data, start)
}

// ShowLineProfile godoc
// @Summary Show Line Profile Detail
// @Tags CLI-Profile
// @Router /api/v1/cli/profile/line [post]
func (h *CLIHandler) ShowLineProfile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Name == "" {
		response.BadRequest(w, "name is required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowLineProfile(ctx, req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_pon_onu-profile_gpon_line_name", data, start)
}

// ShowRemoteProfileList godoc
// @Summary Show Remote Profile List
// @Tags CLI-Profile
// @Router /api/v1/cli/profile/remote/list [post]
func (h *CLIHandler) ShowRemoteProfileList(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowRemoteProfileList(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_pon_onu-profile_gpon_remote", data, start)
}

// ShowRemoteProfile godoc
// @Summary Show Remote Profile Detail
// @Tags CLI-Profile
// @Router /api/v1/cli/profile/remote [post]
func (h *CLIHandler) ShowRemoteProfile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Name == "" {
		response.BadRequest(w, "name is required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowRemoteProfile(ctx, req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_pon_onu-profile_gpon_remote_name", data, start)
}

// ============================================================
// PRIORITY 5: VLAN
// ============================================================

// ShowVLANList godoc
// @Summary Show VLAN List
// @Tags CLI-VLAN
// @Router /api/v1/cli/vlan/list [post]
func (h *CLIHandler) ShowVLANList(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowVLANList(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_vlan", data, start)
}

// ShowVLANByID godoc
// @Summary Show VLAN by ID
// @Tags CLI-VLAN
// @Router /api/v1/cli/vlan/{id} [post]
func (h *CLIHandler) ShowVLANByID(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.VlanID == 0 {
		response.BadRequest(w, "vlan_id is required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowVLANByID(ctx, req.VlanID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_vlan_id", data, start)
}

// ============================================================
// PRIORITY 6: IGMP/MULTICAST
// ============================================================

// ShowIGMPMVlan godoc
// @Summary Show IGMP MVLAN List
// @Tags CLI-IGMP
// @Router /api/v1/cli/igmp/mvlan [post]
func (h *CLIHandler) ShowIGMPMVlan(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowIGMPMVlan(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_igmp_mvlan", data, start)
}

// ShowIGMPMVlanByID godoc
// @Summary Show IGMP MVLAN by ID
// @Tags CLI-IGMP
// @Router /api/v1/cli/igmp/mvlan/{id} [post]
func (h *CLIHandler) ShowIGMPMVlanByID(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.VlanID == 0 {
		response.BadRequest(w, "vlan_id is required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowIGMPMVlanByID(ctx, req.VlanID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_igmp_mvlan_id", data, start)
}

// ShowIGMPDynamicMember godoc
// @Summary Show IGMP Dynamic Member
// @Tags CLI-IGMP
// @Router /api/v1/cli/igmp/dynamic-member [post]
func (h *CLIHandler) ShowIGMPDynamicMember(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	output, err := client.ShowIGMPDynamicMember(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_igmp_dynamic_member", output, start)
}

// ShowIGMPForwardingTable godoc
// @Summary Show IGMP Forwarding Table
// @Tags CLI-IGMP
// @Router /api/v1/cli/igmp/forwarding-table [post]
func (h *CLIHandler) ShowIGMPForwardingTable(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	output, err := client.ShowIGMPForwardingTable(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_igmp_forwarding_table", output, start)
}

// ShowIGMPInterface godoc
// @Summary Show IGMP Interface
// @Tags CLI-IGMP
// @Router /api/v1/cli/igmp/interface [post]
func (h *CLIHandler) ShowIGMPInterface(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	output, err := client.ShowIGMPInterface(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_igmp_interface", output, start)
}

// ============================================================
// PRIORITY 7: INTERFACE DETAIL
// ============================================================

// ShowInterfaceByType godoc
// @Summary Show Interface by Type
// @Tags CLI-Interface
// @Router /api/v1/cli/interface/detail [post]
func (h *CLIHandler) ShowInterfaceByType(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Name == "" {
		response.BadRequest(w, "name is required (e.g., gpon-olt_1/1/1, gei_1/4/1)")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowInterfaceByType(ctx, req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_interface_detail", data, start)
}

// ============================================================
// PRIORITY 8: USER MANAGEMENT
// ============================================================

// ShowOnlineUsers godoc
// @Summary Show Online Users
// @Tags CLI-User
// @Router /api/v1/cli/user/online [post]
func (h *CLIHandler) ShowOnlineUsers(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowOnlineUsers(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_users", data, start)
}

// ============================================================
// REMAINING READ ENDPOINTS
// ============================================================

// ShowDialPlanProfile godoc
// @Summary Show Dial Plan Profile
// @Tags CLI-GPON
// @Router /api/v1/cli/gpon/dial-plan [post]
func (h *CLIHandler) ShowDialPlanProfile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Name == "" {
		response.BadRequest(w, "name is required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowDialPlanProfile(ctx, req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_gpon_onu_profile_dial-plan", data, start)
}

// ShowVoipAccesscodeProfile godoc
// @Summary Show VoIP Accesscode Profile
// @Tags CLI-GPON
// @Router /api/v1/cli/gpon/voip-accesscode [post]
func (h *CLIHandler) ShowVoipAccesscodeProfile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Name == "" {
		response.BadRequest(w, "name is required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowVoipAccesscodeProfile(ctx, req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_gpon_onu_profile_voip-accesscode", data, start)
}

// ShowVoipAppsrvProfile godoc
// @Summary Show VoIP Appsrv Profile
// @Tags CLI-GPON
// @Router /api/v1/cli/gpon/voip-appsrv [post]
func (h *CLIHandler) ShowVoipAppsrvProfile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Name == "" {
		response.BadRequest(w, "name is required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowVoipAppsrvProfile(ctx, req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_gpon_onu_profile_voip-appsrv", data, start)
}

// ShowSNMPCommunity godoc
// @Summary Show SNMP Community
// @Tags CLI-SNMP
// @Router /api/v1/cli/snmp/community [post]
func (h *CLIHandler) ShowSNMPCommunity(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowSNMPCommunity(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_snmp_community", data, start)
}

// ShowSNMPHost godoc
// @Summary Show SNMP Host
// @Tags CLI-SNMP
// @Router /api/v1/cli/snmp/host [post]
func (h *CLIHandler) ShowSNMPHost(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowSNMPHost(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_snmp_host", data, start)
}

// ShowRunningConfig godoc
// @Summary Show Running Config
// @Tags CLI-Config
// @Router /api/v1/cli/config/running [post]
func (h *CLIHandler) ShowRunningConfig(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	output, err := client.ShowRunningConfig(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_running-config", output, start)
}

// SaveConfig godoc
// @Summary Save Config
// @Tags CLI-Config
// @Router /api/v1/cli/config/save [post]
func (h *CLIHandler) SaveConfig(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	err := client.SaveConfig(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "write", map[string]interface{}{
		"success": true,
		"message": "Configuration saved successfully",
	}, start)
}

// BackupConfig godoc
// @Summary Backup Config to TFTP
// @Tags CLI-Config
// @Router /api/v1/cli/config/backup [post]
func (h *CLIHandler) BackupConfig(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Name == "" || req.SN == "" {
		response.BadRequest(w, "name (tftp_ip) and sn (filename) are required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	output, err := client.BackupConfig(ctx, req.Name, req.SN)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "copy_running-config_tftp", map[string]interface{}{
		"success":  true,
		"message":  output,
		"tftp_ip":  req.Name,
		"filename": req.SN,
	}, start)
}

// RestoreConfig godoc
// @Summary Restore Config from TFTP
// @Tags CLI-Config
// @Router /api/v1/cli/config/restore [post]
func (h *CLIHandler) RestoreConfig(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Name == "" || req.SN == "" {
		response.BadRequest(w, "name (tftp_ip) and sn (filename) are required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	output, err := client.RestoreConfig(ctx, req.Name, req.SN)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "copy_tftp_running-config", map[string]interface{}{
		"success":  true,
		"message":  output,
		"tftp_ip":  req.Name,
		"filename": req.SN,
	}, start)
}

// ShowInterfaceVLAN godoc
// @Summary Show Interface VLAN
// @Tags CLI-Interface
// @Router /api/v1/cli/interface/vlan [post]
func (h *CLIHandler) ShowInterfaceVLAN(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.VlanID == 0 {
		response.BadRequest(w, "vlan_id is required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowInterfaceVLAN(ctx, req.VlanID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_interface_vlan", data, start)
}

// ShowPowerSupply godoc
// @Summary Show Power Supply
// @Tags CLI-Hardware
// @Router /api/v1/cli/power [post]
func (h *CLIHandler) ShowPowerSupply(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowPowerSupply(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_power", data, start)
}

// ShowTemperature godoc
// @Summary Show Temperature
// @Tags CLI-Hardware
// @Router /api/v1/cli/temperature [post]
func (h *CLIHandler) ShowTemperature(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	data, err := client.ShowTemperature(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "show_temperature", data, start)
}

// ============================================================
// WRITE ENDPOINTS (Provisioning)
// ============================================================

// AuthenticateONU godoc
// @Summary Authenticate ONU
// @Tags CLI-ONU
// @Router /api/v1/cli/onu/auth [post]
func (h *CLIHandler) AuthenticateONU(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Slot == 0 || req.OnuID == 0 || req.OnuType == "" || req.SN == "" {
		response.BadRequest(w, "slot, onu_id, onu_type, and sn are required")
		return
	}

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

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	err := client.AuthenticateONU(ctx, rack, shelf, req.Slot, req.OnuID, req.OnuType, req.SN)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "onu_authenticate", map[string]interface{}{
		"success":  true,
		"message":  fmt.Sprintf("ONU %d authenticated", req.OnuID),
		"onu_id":   req.OnuID,
		"onu_type": req.OnuType,
		"sn":       req.SN,
	}, start)
}

// DeleteONU godoc
// @Summary Delete ONU
// @Tags CLI-ONU
// @Router /api/v1/cli/onu/delete [post]
func (h *CLIHandler) DeleteONU(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Slot == 0 || req.OnuID == 0 {
		response.BadRequest(w, "slot and onu_id are required")
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

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	err := client.DeleteONU(ctx, rack, shelf, req.Slot, req.OnuID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "onu_delete", map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("ONU %d deleted", req.OnuID),
		"onu_id":  req.OnuID,
	}, start)
}

// ============================================================
// WRITE OPERATIONS - ONU PROVISIONING
// ============================================================

// RenameONU godoc
// @Summary Rename ONU
// @Tags CLI-ONU-WRITE
// @Router /api/v1/cli/onu/rename [post]
func (h *CLIHandler) RenameONU(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Slot == 0 || req.OnuID == 0 || req.Name == "" {
		response.BadRequest(w, "slot, onu_id, and name are required")
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

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	err := client.RenameONU(ctx, rack, shelf, req.Slot, req.OnuID, req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "onu_rename", map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("ONU %d renamed to %s", req.OnuID, req.Name),
		"onu_id":  req.OnuID,
		"name":    req.Name,
	}, start)
}

// ResetONU godoc
// @Summary Reset/Reboot ONU
// @Tags CLI-ONU-WRITE
// @Router /api/v1/cli/onu/reset [post]
func (h *CLIHandler) ResetONU(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Slot == 0 || req.OnuID == 0 {
		response.BadRequest(w, "slot and onu_id are required")
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

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	err := client.ResetONU(ctx, rack, shelf, req.Slot, req.OnuID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "onu_reset", map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("ONU %d reset/rebooted", req.OnuID),
		"onu_id":  req.OnuID,
	}, start)
}

// ============================================================
// WRITE OPERATIONS - T-CONT & GEM PORT
// ============================================================

// CreateTCONT godoc
// @Summary Create T-CONT
// @Tags CLI-TCONT-WRITE
// @Router /api/v1/cli/tcont/create [post]
func (h *CLIHandler) CreateTCONT(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Slot == 0 || req.OnuID == 0 || req.Name == "" || req.SN == "" {
		response.BadRequest(w, "slot, onu_id, name (tcont name), and sn (profile name) are required")
		return
	}

	rack := req.Rack
	shelf := req.Shelf
	tcontID := req.VlanID
	if tcontID == 0 {
		tcontID = 1
	}
	if rack == 0 {
		rack = 1
	}
	if shelf == 0 {
		shelf = 1
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	err := client.CreateTCONT(ctx, rack, shelf, req.Slot, req.OnuID, tcontID, req.Name, req.SN)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "tcont_create", map[string]interface{}{
		"success":    true,
		"message":    fmt.Sprintf("T-CONT %d created", tcontID),
		"tcont_id":   tcontID,
		"name":       req.Name,
		"profile":    req.SN,
		"interface":  fmt.Sprintf("gpon-onu_%d/%d/%d:%d", rack, shelf, req.Slot, req.OnuID),
	}, start)
}

// CreateGEMPort godoc
// @Summary Create GEM Port
// @Tags CLI-GEMPORT-WRITE
// @Router /api/v1/cli/gemport/create [post]
func (h *CLIHandler) CreateGEMPort(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Slot == 0 || req.OnuID == 0 || req.Name == "" {
		response.BadRequest(w, "slot, onu_id, and name are required")
		return
	}

	rack := req.Rack
	shelf := req.Shelf
	gemportID := req.VlanID
	tcontID := req.Port
	if gemportID == 0 {
		gemportID = 1
	}
	if tcontID == 0 {
		tcontID = 1
	}
	if rack == 0 {
		rack = 1
	}
	if shelf == 0 {
		shelf = 1
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	err := client.CreateGEMPort(ctx, rack, shelf, req.Slot, req.OnuID, gemportID, tcontID, req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "gemport_create", map[string]interface{}{
		"success":    true,
		"message":    fmt.Sprintf("GEM Port %d created", gemportID),
		"gemport_id": gemportID,
		"tcont_id":   tcontID,
		"name":       req.Name,
		"interface":  fmt.Sprintf("gpon-onu_%d/%d/%d:%d", rack, shelf, req.Slot, req.OnuID),
	}, start)
}

// ============================================================
// WRITE OPERATIONS - SERVICE PORT
// ============================================================

// CreateServicePort godoc
// @Summary Create Service Port
// @Tags CLI-SERVICEPORT-WRITE
// @Router /api/v1/cli/service-port/create [post]
func (h *CLIHandler) CreateServicePort(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Slot == 0 || req.OnuID == 0 || req.VlanID == 0 {
		response.BadRequest(w, "slot, onu_id, and vlan_id are required")
		return
	}

	rack := req.Rack
	shelf := req.Shelf
	servicePortID := req.Port
	vport := req.VlanID
	if servicePortID == 0 {
		servicePortID = 1
	}
	if vport == 0 {
		vport = 1
	}
	if rack == 0 {
		rack = 1
	}
	if shelf == 0 {
		shelf = 1
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	err := client.CreateServicePort(ctx, rack, shelf, req.Slot, req.OnuID, servicePortID, vport, req.VlanID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "service-port_create", map[string]interface{}{
		"success":        true,
		"message":        fmt.Sprintf("Service Port %d created", servicePortID),
		"service_port_id": servicePortID,
		"vport":          vport,
		"vlan_id":        req.VlanID,
		"interface":      fmt.Sprintf("gpon-onu_%d/%d/%d:%d", rack, shelf, req.Slot, req.OnuID),
	}, start)
}

// DeleteServicePort godoc
// @Summary Delete Service Port
// @Tags CLI-SERVICEPORT-WRITE
// @Router /api/v1/cli/service-port/delete [post]
func (h *CLIHandler) DeleteServicePort(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Slot == 0 || req.OnuID == 0 || req.Port == 0 {
		response.BadRequest(w, "slot, onu_id, and port (service_port_id) are required")
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

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	err := client.DeleteServicePort(ctx, rack, shelf, req.Slot, req.OnuID, req.Port)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "service-port_delete", map[string]interface{}{
		"success":         true,
		"message":         fmt.Sprintf("Service Port %d deleted", req.Port),
		"service_port_id": req.Port,
	}, start)
}

// ============================================================
// WRITE OPERATIONS - VLAN
// ============================================================

// CreateVLAN godoc
// @Summary Create VLAN
// @Tags CLI-VLAN-WRITE
// @Router /api/v1/cli/vlan/create [post]
func (h *CLIHandler) CreateVLAN(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.VlanID == 0 {
		response.BadRequest(w, "vlan_id is required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	err := client.CreateVLAN(ctx, req.VlanID, req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "vlan_create", map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("VLAN %d created", req.VlanID),
		"vlan_id": req.VlanID,
		"name":    req.Name,
	}, start)
}

// DeleteVLAN godoc
// @Summary Delete VLAN
// @Tags CLI-VLAN-WRITE
// @Router /api/v1/cli/vlan/delete [post]
func (h *CLIHandler) DeleteVLAN(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.VlanID == 0 {
		response.BadRequest(w, "vlan_id is required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	err := client.DeleteVLAN(ctx, req.VlanID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "vlan_delete", map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("VLAN %d deleted", req.VlanID),
		"vlan_id": req.VlanID,
	}, start)
}

// AddPortToVLAN godoc
// @Summary Add Port to VLAN
// @Tags CLI-VLAN-WRITE
// @Router /api/v1/cli/vlan/port/add [post]
func (h *CLIHandler) AddPortToVLAN(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Name == "" || req.VlanID == 0 {
		response.BadRequest(w, "name (interface name) and vlan_id are required")
		return
	}

	mode := req.OnuType
	if mode == "" {
		mode = "tag"
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	err := client.AddPortToVLAN(ctx, req.Name, req.VlanID, mode)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "vlan_port_add", map[string]interface{}{
		"success":     true,
		"message":     fmt.Sprintf("Port %s added to VLAN %d", req.Name, req.VlanID),
		"interface":   req.Name,
		"vlan_id":     req.VlanID,
		"mode":        mode,
	}, start)
}

// ============================================================
// WRITE OPERATIONS - PROFILE CREATION
// ============================================================

// CreateLineProfile godoc
// @Summary Create Line Profile
// @Tags CLI-PROFILE-WRITE
// @Router /api/v1/cli/profile/line/create [post]
func (h *CLIHandler) CreateLineProfile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Name == "" {
		response.BadRequest(w, "name is required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	err := client.CreateLineProfile(ctx, req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "profile_line_create", map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Line Profile %s created", req.Name),
		"name":    req.Name,
	}, start)
}

// CreateRemoteProfile godoc
// @Summary Create Remote Profile
// @Tags CLI-PROFILE-WRITE
// @Router /api/v1/cli/profile/remote/create [post]
func (h *CLIHandler) CreateRemoteProfile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Name == "" {
		response.BadRequest(w, "name is required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	err := client.CreateRemoteProfile(ctx, req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "profile_remote_create", map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Remote Profile %s created", req.Name),
		"name":    req.Name,
	}, start)
}

// CreateVLANProfile godoc
// @Summary Create VLAN Profile
// @Tags CLI-PROFILE-WRITE
// @Router /api/v1/cli/profile/vlan/create [post]
func (h *CLIHandler) CreateVLANProfile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Name == "" || req.VlanID == 0 {
		response.BadRequest(w, "name and vlan_id are required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	err := client.CreateVLANProfile(ctx, req.Name, req.VlanID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "profile_vlan_create", map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("VLAN Profile %s created", req.Name),
		"name":    req.Name,
		"vlan_id": req.VlanID,
	}, start)
}

// CreateTCONTProfile godoc
// @Summary Create T-CONT Profile
// @Tags CLI-PROFILE-WRITE
// @Router /api/v1/cli/profile/tcont/create [post]
func (h *CLIHandler) CreateTCONTProfile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.Name == "" || req.OnuType == "" {
		response.BadRequest(w, "name and onu_type (profile type) are required")
		return
	}

	bandwidth := req.VlanID
	if bandwidth == 0 {
		bandwidth = 10000
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	err := client.CreateTCONTProfile(ctx, req.Name, req.OnuType, bandwidth)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "profile_tcont_create", map[string]interface{}{
		"success":    true,
		"message":    fmt.Sprintf("T-CONT Profile %s created", req.Name),
		"name":       req.Name,
		"type":       req.OnuType,
		"bandwidth":  bandwidth,
	}, start)
}

// ============================================================
// WRITE OPERATIONS - IGMP/MULTICAST
// ============================================================

// EnableIGMP godoc
// @Summary Enable IGMP
// @Tags CLI-IGMP-WRITE
// @Router /api/v1/cli/igmp/enable [post]
func (h *CLIHandler) EnableIGMP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	err := client.EnableIGMP(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "igmp_enable", map[string]interface{}{
		"success": true,
		"message": "IGMP enabled successfully",
	}, start)
}

// CreateMVLAN godoc
// @Summary Create MVLAN
// @Tags CLI-IGMP-WRITE
// @Router /api/v1/cli/mvlan/create [post]
func (h *CLIHandler) CreateMVLAN(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.VlanID == 0 {
		response.BadRequest(w, "vlan_id (mvlan_id) is required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	err := client.CreateMVLAN(ctx, req.VlanID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "mvlan_create", map[string]interface{}{
		"success":  true,
		"message":  fmt.Sprintf("MVLAN %d created", req.VlanID),
		"mvlan_id": req.VlanID,
	}, start)
}

// AddMVLANGroup godoc
// @Summary Add MVLAN Group
// @Tags CLI-IGMP-WRITE
// @Router /api/v1/cli/mvlan/group/add [post]
func (h *CLIHandler) AddMVLANGroup(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req CLIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request")
		return
	}

	if req.VlanID == 0 || req.Name == "" {
		response.BadRequest(w, "vlan_id (mvlan_id) and name (group_ip) are required")
		return
	}

	ctx := context.Background()
	client := h.getClient(req)
	if err := client.Connect(); err != nil {
		response.Error(w, http.StatusGatewayTimeout, "Connection failed")
		return
	}
	defer client.Close()

	err := client.AddMVLANGroup(ctx, req.VlanID, req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respond(w, "mvlan_group_add", map[string]interface{}{
		"success":   true,
		"message":   fmt.Sprintf("Group %s added to MVLAN %d", req.Name, req.VlanID),
		"mvlan_id":  req.VlanID,
		"group_ip":  req.Name,
	}, start)
}
