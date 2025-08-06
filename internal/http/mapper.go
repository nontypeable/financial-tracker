package http

import (
	"errors"
	"net/http"

	apperror "github.com/nontypeable/financial-tracker/internal/errors"
)

func MapAppErrorToHTTP(err error) (status int, message string) {
	switch {
	case errors.Is(err, apperror.ErrUserAlreadyExists):
		return http.StatusConflict, "user already exists"
	case errors.Is(err, apperror.ErrInvalidInput):
		return http.StatusBadRequest, "invalid input"
	case errors.Is(err, apperror.ErrUserNotFound):
		return http.StatusNotFound, "user not found"
	case errors.Is(err, apperror.ErrInvalidCredentials):
		return http.StatusUnauthorized, "invalid credentials"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
