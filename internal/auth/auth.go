package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
		panic("secrets cannot be empty")
	}

	if accessTTL <= 0 || refreshTTL <= 0 {
		panic("ttl values must be positive")
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
		return "", fmt.Errorf("invalid user ID")
	}

	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.accessTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        generateRandomString(32),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.accessSecret))
}

func (tm *tokenManager) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	if userID == uuid.Nil {
		return "", fmt.Errorf("invalid user ID")
	}

	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.refreshTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        generateRandomString(32),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.refreshSecret))
}

func (tm *tokenManager) ValidateAccessToken(token string) (*jwt.RegisteredClaims, error) {
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	return tm.parseToken(token, tm.accessSecret)
}

func (tm *tokenManager) ValidateRefreshToken(token string) (*jwt.RegisteredClaims, error) {
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	return tm.parseToken(token, tm.refreshSecret)
}

func (tm *tokenManager) parseToken(tokenStr, secret string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func generateRandomString(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
