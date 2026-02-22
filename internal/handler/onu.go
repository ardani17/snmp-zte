package handler

import (
	"net/http"
	"strconv"

	"github.com/ardani/snmp-zte/internal/service"
	"github.com/ardani/snmp-zte/pkg/response"
	"github.com/go-chi/chi/v5"
)

// ONUHandler menangani permintaan query ONU (modem pelanggan).
type ONUHandler struct {
	service *service.ONUService
}

// NewONUHandler membuat instance ONU handler baru.
func NewONUHandler(service *service.ONUService) *ONUHandler {
	return &ONUHandler{service: service}
}

// List mengembalikan daftar ONU untuk Board dan port PON tertentu.
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

// Detail mengambil informasi rinci untuk satu ONU tunggal.
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

// EmptySlots mengambil slot ONU yang tersedia (belum terpakai) di port PON.
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

// ClearCache menghapus data cache ONU untuk port PON tertentu.
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
