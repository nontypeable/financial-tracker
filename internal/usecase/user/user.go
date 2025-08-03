package user

import (
	"context"
	"fmt"

	"github.com/nontypeable/financial-tracker/internal/auth"
	"github.com/nontypeable/financial-tracker/internal/domain/user"
)

type service struct {
	repository   user.Repository
	tokenManager auth.TokenManager
}

func NewService(repository user.Repository, tokenManager auth.TokenManager) user.Service {
	return &service{
		repository:   repository,
		tokenManager: tokenManager,
	}
}

func (s *service) SignUp(ctx context.Context, email, password, firstName, lastName string) (string, string, error) {
	user, err := user.NewUser(email, password, firstName, lastName)
	if err != nil {
		return "", "", err
	}

	id, err := s.repository.Create(ctx, user)
	if err != nil {
		return "", "", err
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(id)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(id)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *service) SignIn(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("failed to get user by email: %w", err)
	}

	if !user.CheckPassword(password) {
		return "", "", fmt.Errorf("incorrect password")
	}

	userID := user.ID

	accessToken, err := s.tokenManager.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
