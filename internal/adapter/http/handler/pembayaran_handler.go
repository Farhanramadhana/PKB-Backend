package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bpka/tps-pkb/internal/adapter/http/middleware"
	"github.com/bpka/tps-pkb/internal/adapter/http/response"
	"github.com/bpka/tps-pkb/internal/domain/dto"
	"github.com/bpka/tps-pkb/internal/domain/port"
	"github.com/go-playground/validator/v10"
)

type PembayaranHandler struct {
	svc      port.PembayaranService
	validate *validator.Validate
}

func NewPembayaranHandler(svc port.PembayaranService, validate *validator.Validate) *PembayaranHandler {
	return &PembayaranHandler{svc: svc, validate: validate}
}

func (h *PembayaranHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())

	var req dto.CreatePembayaranRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "request tidak valid")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validationErrors(err)...)
		return
	}

	pembayaran, dendaList, err := h.svc.Create(r.Context(), req, claims.UserID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusCreated, "pembayaran berhasil dicatat", map[string]any{
		"pembayaran": pembayaran,
		"denda":      dendaList,
	})
}
