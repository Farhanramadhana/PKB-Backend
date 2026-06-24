package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bpka/tps-pkb/internal/adapter/http/response"
	"github.com/bpka/tps-pkb/internal/domain/dto"
	"github.com/bpka/tps-pkb/internal/domain/port"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	svc      port.AuthService
	validate *validator.Validate
}

func NewAuthHandler(svc port.AuthService, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{svc: svc, validate: validate}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "request tidak valid")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "request tidak valid", validationErrors(err)...)
		return
	}

	token, err := h.svc.Login(r.Context(), req)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "login berhasil", token)
}
