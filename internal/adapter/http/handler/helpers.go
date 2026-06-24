package handler

import (
	"errors"
	"net/http"

	"github.com/bpka/tps-pkb/internal/adapter/http/response"
	"github.com/bpka/tps-pkb/internal/domain/dto"
	"github.com/go-playground/validator/v10"
)

func validationErrors(err error) []response.FieldError {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return nil
	}
	out := make([]response.FieldError, 0, len(ve))
	for _, fe := range ve {
		out = append(out, response.FieldError{
			Field:   fe.Field(),
			Message: fe.Tag(),
		})
	}
	return out
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, dto.ErrNotFound):
		response.WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
	case errors.Is(err, dto.ErrConflict):
		response.WriteError(w, http.StatusConflict, "CONFLICT", err.Error())
	case errors.Is(err, dto.ErrWajibPajakInaktif),
		errors.Is(err, dto.ErrTanggalMasaDepan),
		errors.Is(err, dto.ErrJumlahHarusPositif):
		response.WriteError(w, http.StatusUnprocessableEntity, "UNPROCESSABLE", err.Error())
	case errors.Is(err, dto.ErrUnauthorized):
		response.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
	case errors.Is(err, dto.ErrForbidden):
		response.WriteError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
	default:
		response.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "terjadi kesalahan pada server")
	}
}

func pageParams(r *http.Request) (page, limit int) {
	page = 1
	limit = 20
	if p := r.URL.Query().Get("page"); p != "" {
		if v := parseInt(p); v > 0 {
			page = v
		}
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		if v := parseInt(l); v > 0 && v <= 100 {
			limit = v
		}
	}
	return
}

func parseInt(s string) int {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	return n
}

func parseBoolPtr(s string) *bool {
	if s == "true" {
		v := true
		return &v
	}
	if s == "false" {
		v := false
		return &v
	}
	return nil
}
