package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bpka/tps-pkb/internal/adapter/http/middleware"
	"github.com/bpka/tps-pkb/internal/adapter/http/response"
	"github.com/bpka/tps-pkb/internal/domain/dto"
	"github.com/bpka/tps-pkb/internal/domain/port"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type KewajibanHandler struct {
	svc      port.KewajibanService
	validate *validator.Validate
}

func NewKewajibanHandler(svc port.KewajibanService, validate *validator.Validate) *KewajibanHandler {
	return &KewajibanHandler{svc: svc, validate: validate}
}

func (h *KewajibanHandler) List(w http.ResponseWriter, r *http.Request) {
	page, limit := pageParams(r)
	filter := dto.KewajibanFilter{
		Status: r.URL.Query().Get("status"),
		Page:   page,
		Limit:  limit,
	}

	// WAJIB_PAJAK role: filter by own wajib_pajak_id
	if ownerID := middleware.OwnerFilterFromContext(r.Context()); ownerID != nil {
		filter.WajibPajakID = ownerID
	} else if wpIDStr := r.URL.Query().Get("wajib_pajak_id"); wpIDStr != "" {
		wpID, err := uuid.Parse(wpIDStr)
		if err == nil {
			filter.WajibPajakID = &wpID
		}
	}

	items, total, err := h.svc.List(r.Context(), filter)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.WriteSuccessPaginated(w, http.StatusOK, "data kewajiban pajak", items,
		response.Meta{Page: page, Limit: limit, Total: total})
}

func (h *KewajibanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateKewajibanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "request tidak valid")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validationErrors(err)...)
		return
	}

	kewajiban, err := h.svc.Create(r.Context(), req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusCreated, "kewajiban pajak berhasil dibuat", kewajiban)
}
