package apperror

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidUserID      = errors.New("invalid user id")
	ErrInvalidInput       = errors.New("invalid input")
	ErrUserNotFound       = errors.New("user is not found")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrTokenIsEmpty         = errors.New("token is empty")
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidTokenClaims   = errors.New("invalid token claims")
	ErrEmptyTokenSecret     = errors.New("token secrets cannot be empty")
	ErrInvalidTokenLifetime = errors.New("token TTL must be positive")

	ErrAccountNotFound = errors.New("account is not found")
)
