package handler

import (
	"net/http"

	"github.com/bpka/tps-pkb/internal/adapter/http/response"
	"github.com/bpka/tps-pkb/internal/domain/dto"
	"github.com/bpka/tps-pkb/internal/domain/port"
)

type LaporanHandler struct {
	svc port.LaporanService
}

func NewLaporanHandler(svc port.LaporanService) *LaporanHandler {
	return &LaporanHandler{svc: svc}
}

func (h *LaporanHandler) Get(w http.ResponseWriter, r *http.Request) {
	filter := dto.LaporanFilter{
		StartDate: r.URL.Query().Get("start_date"),
		EndDate:   r.URL.Query().Get("end_date"),
		Status:    r.URL.Query().Get("status"),
	}

	laporan, err := h.svc.Get(r.Context(), filter)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, "laporan pembayaran pajak", laporan)
}
