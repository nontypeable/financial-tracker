package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nontypeable/financial-tracker/internal/auth"
	"github.com/nontypeable/financial-tracker/internal/domain/user"
	apperror "github.com/nontypeable/financial-tracker/internal/errors"
	"golang.org/x/crypto/bcrypt"
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
		return "", "", fmt.Errorf("invalid user data: %w", err)
	}

	id, err := s.repository.Create(ctx, user)
	if err != nil {
		return "", "", fmt.Errorf("create user: %w", err)
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(id)
	if err != nil {
		return "", "", fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(id)
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *service) SignIn(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) {
			return "", "", apperror.ErrInvalidCredentials
		}
		return "", "", fmt.Errorf("get user by email: %w", err)
	}

	ok, err := user.CheckPassword(password)
	if err != nil {
		return "", "", fmt.Errorf("check password error: %w", err)
	}
	if !ok {
		return "", "", apperror.ErrInvalidCredentials
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *service) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	token, err := s.tokenManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("validate refresh token: %w", err)
	}

	userID, err := uuid.Parse(token.Subject)
	if err != nil {
		return "", "", fmt.Errorf("invalid token subject (user id): %w", err)
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("generate access token: %w", err)
	}

	newRefreshToken, err := s.tokenManager.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	return accessToken, newRefreshToken, nil
}

func (s *service) GetUserInfo(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user info: %w", err)
	}

	return user, nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, firstName, lastName string) error {
	user, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	var updated bool

	if firstName != "" && firstName != user.FirstName {
		user.FirstName = firstName
		updated = true
	}

	if lastName != "" && lastName != user.LastName {
		user.LastName = lastName
		updated = true
	}

	if !updated {
		return nil
	}

	if err := s.repository.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

func (s *service) ChangeEmail(ctx context.Context, id uuid.UUID, newEmail string, currentPassword string) error {
	user, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	ok, err := user.CheckPassword(currentPassword)
	if err != nil {
		return fmt.Errorf("check password error: %w", err)
	}
	if !ok {
		return apperror.ErrInvalidCredentials
	}

	exists, err := s.repository.EmailExists(ctx, newEmail)
	if err != nil {
		return fmt.Errorf("check email existence: %w", err)
	}
	if exists {
		return apperror.ErrUserAlreadyExists
	}

	user.Email = newEmail

	if err := s.repository.Update(ctx, user); err != nil {
		return fmt.Errorf("update email: %w", err)
	}

	return nil
}

func (s *service) ChangePassword(ctx context.Context, id uuid.UUID, newPassword string, currentPassword string) error {
	user, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	ok, err := user.CheckPassword(currentPassword)
	if err != nil {
		return fmt.Errorf("check password error: %w", err)
	}
	if !ok {
		return apperror.ErrInvalidCredentials
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	user.PasswordHash = string(hashed)

	if err := s.repository.Update(ctx, user); err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	return nil
}
