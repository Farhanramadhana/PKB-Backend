package handler

import (
	"net/http"

	"github.com/bpka/tps-pkb/internal/adapter/http/response"
	"github.com/bpka/tps-pkb/internal/domain/port"
	"github.com/google/uuid"
)

type DendaHandler struct {
	svc port.DendaService
}

func NewDendaHandler(svc port.DendaService) *DendaHandler {
	return &DendaHandler{svc: svc}
}

// Get returns calculated fines for a kewajiban. Read-only — does not persist.
func (h *DendaHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "id tidak valid")
		return
	}

	dendaList, err := h.svc.GetByKewajibanID(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, "data denda", dendaList)
}
