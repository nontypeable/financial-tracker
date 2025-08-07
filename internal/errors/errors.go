package apperror

import "errors"

var (
	// Validation errors
	ErrInvalidInput     = errors.New("invalid input")
	ErrValidationFailed = errors.New("validation failed")

	// User-related errors
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")

	// Token-related errors
	ErrTokenIsEmpty         = errors.New("token is empty")
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidTokenClaims   = errors.New("invalid token claims")
	ErrEmptyTokenSecret     = errors.New("token secret cannot be empty")
	ErrInvalidTokenLifetime = errors.New("token TTL must be positive")

	// Account-related errors
	ErrAccountNotFound = errors.New("account is not found")

	// Transaction-related errors
	ErrTransactionNotFound = errors.New("transaction is not found")

	// Request-related errors
	ErrNilResponseWriter      = errors.New("response writer is nil")
	ErrNilRequest             = errors.New("request is nil")
	ErrNilDestination         = errors.New("destination is nil")
	ErrUnsupportedMethod      = errors.New("unsupported HTTP method")
	ErrUnsupportedContentType = errors.New("unsupported content type")
)
