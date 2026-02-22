package handler

import (
	"net/http"
	"strconv"

	"github.com/ardani/snmp-zte/internal/service"
	"github.com/ardani/snmp-zte/pkg/response"
	"github.com/go-chi/chi/v5"
)

// ONUHandler handles ONU query requests
type ONUHandler struct {
	service *service.ONUService
}

// NewONUHandler creates a new ONU handler
func NewONUHandler(service *service.ONUService) *ONUHandler {
	return &ONUHandler{service: service}
}

// List godoc
// @Summary List ONUs
// @Description Get list of ONUs for a specific Board and PON port
// @Tags ONU
// @Produce json
// @Param olt_id path string true "OLT ID"
// @Param board_id path int true "Board ID"
// @Param pon_id path int true "PON ID"
// @Success 200 {array} model.ONU
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id} [get]
func (h *ONUHandler) List(w http.ResponseWriter, r *http.Request) {
	oltID := chi.URLParam(r, "olt_id")
	boardID, err := strconv.Atoi(chi.URLParam(r, "board_id"))
	if err != nil {
		response.BadRequest(w, "Invalid board ID")
		return
	}

	ponID, err := strconv.Atoi(chi.URLParam(r, "pon_id"))
	if err != nil {
		response.BadRequest(w, "Invalid PON ID")
		return
	}

	onuList, err := h.service.GetONUList(r.Context(), oltID, boardID, ponID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, onuList)
}

// Detail godoc
// @Summary Get ONU Detail
// @Description Get detailed information for a single ONU
// @Tags ONU
// @Produce json
// @Param olt_id path string true "OLT ID"
// @Param board_id path int true "Board ID"
// @Param pon_id path int true "PON ID"
// @Param onu_id path int true "ONU ID"
// @Success 200 {object} model.ONUDetail
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id}/onu/{onu_id} [get]
func (h *ONUHandler) Detail(w http.ResponseWriter, r *http.Request) {
	oltID := chi.URLParam(r, "olt_id")
	boardID, err := strconv.Atoi(chi.URLParam(r, "board_id"))
	if err != nil {
		response.BadRequest(w, "Invalid board ID")
		return
	}

	ponID, err := strconv.Atoi(chi.URLParam(r, "pon_id"))
	if err != nil {
		response.BadRequest(w, "Invalid PON ID")
		return
	}

	onuID, err := strconv.Atoi(chi.URLParam(r, "onu_id"))
	if err != nil {
		response.BadRequest(w, "Invalid ONU ID")
		return
	}

	detail, err := h.service.GetONUDetail(r.Context(), oltID, boardID, ponID, onuID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, detail)
}

// EmptySlots godoc
// @Summary Get Empty ONU Slots
// @Description Get available ONU slots for a specific Board and PON port
// @Tags ONU
// @Produce json
// @Param olt_id path string true "OLT ID"
// @Param board_id path int true "Board ID"
// @Param pon_id path int true "PON ID"
// @Success 200 {array} int
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id}/empty [get]
func (h *ONUHandler) EmptySlots(w http.ResponseWriter, r *http.Request) {
	oltID := chi.URLParam(r, "olt_id")
	boardID, err := strconv.Atoi(chi.URLParam(r, "board_id"))
	if err != nil {
		response.BadRequest(w, "Invalid board ID")
		return
	}

	ponID, err := strconv.Atoi(chi.URLParam(r, "pon_id"))
	if err != nil {
		response.BadRequest(w, "Invalid PON ID")
		return
	}

	slots, err := h.service.GetEmptySlots(r.Context(), oltID, boardID, ponID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, slots)
}

// ClearCache godoc
// @Summary Clear ONU Cache
// @Description Clear cached ONU data for a specific Board and PON port
// @Tags ONU
// @Produce json
// @Param olt_id path string true "OLT ID"
// @Param board_id path int true "Board ID"
// @Param pon_id path int true "PON ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id}/cache [delete]
func (h *ONUHandler) ClearCache(w http.ResponseWriter, r *http.Request) {
	oltID := chi.URLParam(r, "olt_id")
	boardID, err := strconv.Atoi(chi.URLParam(r, "board_id"))
	if err != nil {
		response.BadRequest(w, "Invalid board ID")
		return
	}

	ponID, err := strconv.Atoi(chi.URLParam(r, "pon_id"))
	if err != nil {
		response.BadRequest(w, "Invalid PON ID")
		return
	}

	if err := h.service.ClearCache(r.Context(), oltID, boardID, ponID); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"message":  "Cache cleared successfully",
		"olt_id":   oltID,
		"board_id": boardID,
		"pon_id":   ponID,
	})
}
