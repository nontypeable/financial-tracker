package http

import (
	"errors"
	"log"
	"net/http"

	apperror "github.com/nontypeable/financial-tracker/internal/errors"
)

func MapAppErrorToHTTP(err error) (status int, message string) {
	switch {
	// User-related errors
	case errors.Is(err, apperror.ErrUserAlreadyExists):
		return http.StatusConflict, "user already exists"
	case errors.Is(err, apperror.ErrInvalidUserID):
		return http.StatusBadRequest, "invalid user ID"
	case errors.Is(err, apperror.ErrUserNotFound):
		return http.StatusNotFound, "user not found"
	case errors.Is(err, apperror.ErrInvalidCredentials):
		return http.StatusUnauthorized, "invalid credentials"

	// Validation
	case errors.Is(err, apperror.ErrInvalidInput), errors.Is(err, apperror.ErrValidationFailed):
		return http.StatusBadRequest, "invalid input"

	// Token
	case errors.Is(err, apperror.ErrTokenIsEmpty),
		errors.Is(err, apperror.ErrInvalidToken),
		errors.Is(err, apperror.ErrInvalidTokenClaims):
		return http.StatusUnauthorized, "invalid or missing token"

	// Account
	case errors.Is(err, apperror.ErrAccountNotFound):
		return http.StatusNotFound, "account not found"

	// Transactions
	case errors.Is(err, apperror.ErrTransactionNotFound):
		return http.StatusNotFound, "transaction not found"

	// Request or Technical
	case errors.Is(err, apperror.ErrNilRequest),
		errors.Is(err, apperror.ErrNilResponseWriter),
		errors.Is(err, apperror.ErrNilDestination),
		errors.Is(err, apperror.ErrUnsupportedMethod),
		errors.Is(err, apperror.ErrUnsupportedContentType):
		return http.StatusBadRequest, "invalid request"

	default:
		log.Printf("internal server error: %v", err)
		return http.StatusInternalServerError, "internal server error"
	}
}
