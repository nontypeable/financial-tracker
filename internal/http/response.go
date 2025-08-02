package http

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	if w == nil {
		return fmt.Errorf("response writer is nil")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(SuccessResponse{
		StatusCode: statusCode,
		Data:       data,
	})
	if err != nil {
		return fmt.Errorf("failed to encode JSON response: %w", err)
	}

	return nil
}

func Error(w http.ResponseWriter, statusCode int, msg string) error {
	if w == nil {
		return fmt.Errorf("response writer is nil")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(ErrorResponse{
		StatusCode: statusCode,
		Error:      msg,
	})
	if err != nil {
		return fmt.Errorf("failed to encode JSON error response: %w", err)
	}

	return nil
}
