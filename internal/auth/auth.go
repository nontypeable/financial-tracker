package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	apperror "github.com/nontypeable/financial-tracker/internal/errors"
)

type TokenManager interface {
	GenerateAccessToken(userID uuid.UUID) (string, error)
	GenerateRefreshToken(userID uuid.UUID) (string, error)
	ValidateAccessToken(token string) (*jwt.RegisteredClaims, error)
	ValidateRefreshToken(token string) (*jwt.RegisteredClaims, error)
}

type tokenManager struct {
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewTokenManager(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) TokenManager {
	if accessSecret == "" || refreshSecret == "" {
		panic(apperror.ErrEmptyTokenSecret.Error())
	}

	if accessTTL <= 0 || refreshTTL <= 0 {
		panic(apperror.ErrInvalidTokenLifetime.Error())
	}

	return &tokenManager{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (tm *tokenManager) GenerateAccessToken(userID uuid.UUID) (string, error) {
	if userID == uuid.Nil {
		return "", apperror.ErrInvalidUserID
	}

	return tm.generateToken(userID, tm.accessSecret, tm.accessTTL)
}

func (tm *tokenManager) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	if userID == uuid.Nil {
		return "", apperror.ErrInvalidUserID
	}

	return tm.generateToken(userID, tm.refreshSecret, tm.refreshTTL)
}

func (tm *tokenManager) ValidateAccessToken(token string) (*jwt.RegisteredClaims, error) {
	if token == "" {
		return nil, apperror.ErrTokenIsEmpty
	}

	return tm.parseToken(token, tm.accessSecret)
}

func (tm *tokenManager) ValidateRefreshToken(token string) (*jwt.RegisteredClaims, error) {
	if token == "" {
		return nil, apperror.ErrTokenIsEmpty
	}

	return tm.parseToken(token, tm.refreshSecret)
}

func (tm *tokenManager) generateToken(userID uuid.UUID, secret string, ttl time.Duration) (string, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        generateRandomString(32),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signedToken, nil
}

func (tm *tokenManager) parseToken(tokenStr, secret string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, apperror.ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, apperror.ErrInvalidTokenClaims
	}

	return claims, nil
}

func generateRandomString(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
