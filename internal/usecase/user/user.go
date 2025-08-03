package user

import (
	"context"

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
