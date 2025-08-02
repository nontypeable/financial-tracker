package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Validable interface {
	Validate() error
}

func DecodeAndValidate[T any](r *http.Request, dest *T) error {
	if r == nil {
		return fmt.Errorf("request is nil")
	}

	if dest == nil {
		return fmt.Errorf("destination is nil")
	}

	if r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodPatch {
		return fmt.Errorf("unsupported method for JSON decoding: %s", r.Method)
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "" && contentType != "application/json" {
		return fmt.Errorf("unsupported content type: %s", contentType)
	}

	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if v, ok := any(dest).(Validable); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	}

	return nil
}
