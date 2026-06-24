package handler

import (
	"net/http"

	"github.com/bpka/tps-pkb/internal/adapter/http/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	response.WriteSuccess(w, http.StatusOK, "service is running", map[string]string{
		"status":  "up",
		"service": "tps-pkb",
	})
}
