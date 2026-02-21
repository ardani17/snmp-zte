package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ardani/snmp-zte/internal/model"
	"github.com/ardani/snmp-zte/internal/service"
	"github.com/ardani/snmp-zte/pkg/response"
	"github.com/go-chi/chi/v5"
)

// OLTHandler handles OLT management requests
type OLTHandler struct {
	service *service.OLTService
}

// NewOLTHandler creates a new OLT handler
func NewOLTHandler(service *service.OLTService) *OLTHandler {
	return &OLTHandler{service: service}
}

// List returns all OLTs
func (h *OLTHandler) List(w http.ResponseWriter, r *http.Request) {
	olts := h.service.List()
	response.JSON(w, http.StatusOK, olts)
}

// Get returns an OLT by ID
func (h *OLTHandler) Get(w http.ResponseWriter, r *http.Request) {
	oltID := chi.URLParam(r, "olt_id")
	if oltID == "" {
		response.BadRequest(w, "OLT ID is required")
		return
	}

	olt, err := h.service.Get(oltID)
	if err != nil {
		response.NotFound(w, "OLT not found: "+oltID)
		return
	}

	response.JSON(w, http.StatusOK, olt)
}

// Create creates a new OLT
func (h *OLTHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Model       string `json:"model"`
		IPAddress   string `json:"ip_address"`
		Port        int    `json:"port"`
		Community   string `json:"community"`
		BoardCount  int    `json:"board_count"`
		PonPerBoard int    `json:"pon_per_board"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Validate required fields
	if req.ID == "" {
		response.BadRequest(w, "ID is required")
		return
	}
	if req.Name == "" {
		response.BadRequest(w, "Name is required")
		return
	}
	if req.Model == "" {
		response.BadRequest(w, "Model is required")
		return
	}
	if req.IPAddress == "" {
		response.BadRequest(w, "IP Address is required")
		return
	}
	if req.Community == "" {
		response.BadRequest(w, "Community is required")
		return
	}

	// Set defaults
	if req.Port == 0 {
		req.Port = 161
	}
	if req.BoardCount == 0 {
		req.BoardCount = 2
	}
	if req.PonPerBoard == 0 {
		req.PonPerBoard = 16
	}

	olt := model.OLT{
		ID:          req.ID,
		Name:        req.Name,
		Model:       req.Model,
		IPAddress:   req.IPAddress,
		Port:        req.Port,
		Community:   req.Community,
		BoardCount:  req.BoardCount,
		PonPerBoard: req.PonPerBoard,
	}

	if err := h.service.Create(olt); err != nil {
		response.Error(w, http.StatusConflict, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, olt)
}

// Update updates an existing OLT
func (h *OLTHandler) Update(w http.ResponseWriter, r *http.Request) {
	oltID := chi.URLParam(r, "olt_id")
	if oltID == "" {
		response.BadRequest(w, "OLT ID is required")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Model       string `json:"model"`
		IPAddress   string `json:"ip_address"`
		Port        int    `json:"port"`
		Community   string `json:"community"`
		BoardCount  int    `json:"board_count"`
		PonPerBoard int    `json:"pon_per_board"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	olt := model.OLT{
		ID:          oltID,
		Name:        req.Name,
		Model:       req.Model,
		IPAddress:   req.IPAddress,
		Port:        req.Port,
		Community:   req.Community,
		BoardCount:  req.BoardCount,
		PonPerBoard: req.PonPerBoard,
	}

	if err := h.service.Update(oltID, olt); err != nil {
		response.NotFound(w, "OLT not found: "+oltID)
		return
	}

	response.JSON(w, http.StatusOK, olt)
}

// Delete removes an OLT
func (h *OLTHandler) Delete(w http.ResponseWriter, r *http.Request) {
	oltID := chi.URLParam(r, "olt_id")
	if oltID == "" {
		response.BadRequest(w, "OLT ID is required")
		return
	}

	if err := h.service.Delete(oltID); err != nil {
		response.NotFound(w, "OLT not found: "+oltID)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "OLT deleted successfully",
		"olt_id":  oltID,
	})
}

func parseIntParam(r *http.Request, param string) (int, error) {
	val := chi.URLParam(r, param)
	return strconv.Atoi(val)
}
