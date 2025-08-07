package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	apperror "github.com/nontypeable/financial-tracker/internal/errors"
)

type Validable interface {
	Validate() error
}

func DecodeAndValidate[T any](r *http.Request, dest *T) error {
	if r == nil {
		return apperror.ErrNilRequest
	}

	if dest == nil {
		return apperror.ErrNilDestination
	}

	switch r.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions:
	default:
		return fmt.Errorf("%w: %s", apperror.ErrUnsupportedMethod, r.Method)
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "" && !strings.HasPrefix(contentType, "application/json") {
		return fmt.Errorf("%w: got '%s'", apperror.ErrUnsupportedContentType, contentType)
	}

	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("%w: failed to decode json: %v", apperror.ErrInvalidInput, err)
	}

	if v, ok := any(dest).(Validable); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("%w: %v", apperror.ErrValidationFailed, err)
		}
	}

	return nil
}
