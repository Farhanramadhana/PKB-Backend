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

type WajibPajakHandler struct {
	svc      port.WajibPajakService
	validate *validator.Validate
}

func NewWajibPajakHandler(svc port.WajibPajakService, validate *validator.Validate) *WajibPajakHandler {
	return &WajibPajakHandler{svc: svc, validate: validate}
}

func (h *WajibPajakHandler) List(w http.ResponseWriter, r *http.Request) {
	page, limit := pageParams(r)
	filter := dto.WajibPajakFilter{
		Nama:     r.URL.Query().Get("nama"),
		Jenis:    r.URL.Query().Get("jenis"),
		IsActive: parseBoolPtr(r.URL.Query().Get("is_active")),
		Page:     page,
		Limit:    limit,
		OwnerID:  middleware.OwnerFilterFromContext(r.Context()),
	}

	items, total, err := h.svc.List(r.Context(), filter)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.WriteSuccessPaginated(w, http.StatusOK, "data wajib pajak", items,
		response.Meta{Page: page, Limit: limit, Total: total})
}

func (h *WajibPajakHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateWajibPajakRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "request tidak valid")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validationErrors(err)...)
		return
	}

	wp, err := h.svc.Create(r.Context(), req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusCreated, "wajib pajak berhasil dibuat", wp)
}

func (h *WajibPajakHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "id tidak valid")
		return
	}

	wp, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, "data wajib pajak", wp)
}

func (h *WajibPajakHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "id tidak valid")
		return
	}

	var req dto.UpdateWajibPajakRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "request tidak valid")
		return
	}

	wp, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, "wajib pajak berhasil diperbarui", wp)
}
