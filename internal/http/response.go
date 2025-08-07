package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	apperror "github.com/nontypeable/financial-tracker/internal/errors"
)

type SuccessResponse struct {
	StatusCode int `json:"status_code"`
	Data       any `json:"data,omitempty"`
}

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

func JSON(w http.ResponseWriter, statusCode int, data any) error {
	return writeJSON(w, statusCode, SuccessResponse{
		StatusCode: statusCode,
		Data:       data,
	})
}

func Error(w http.ResponseWriter, statusCode int, msg string) error {
	return writeJSON(w, statusCode, ErrorResponse{
		StatusCode: statusCode,
		Error:      msg,
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) error {
	if w == nil {
		return apperror.ErrNilResponseWriter
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		return fmt.Errorf("%w: failed to encode json: %v", apperror.ErrInvalidInput, err)
	}
	return nil
}
