package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

type Meta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type Error struct {
	Code    string       `json:"code"`
	Details []FieldError `json:"details,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func WriteSuccess(w http.ResponseWriter, statusCode int, message string, data any) {
	write(w, statusCode, Response{Success: true, Message: message, Data: data})
}

func WriteSuccessPaginated(w http.ResponseWriter, statusCode int, message string, data any, meta Meta) {
	write(w, statusCode, Response{Success: true, Message: message, Data: data, Meta: &meta})
}

func WriteError(w http.ResponseWriter, statusCode int, code, message string, details ...FieldError) {
	write(w, statusCode, Response{
		Success: false,
		Message: message,
		Error:   &Error{Code: code, Details: details},
	})
}

func write(w http.ResponseWriter, statusCode int, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}
