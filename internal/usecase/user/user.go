package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nontypeable/financial-tracker/internal/auth"
	"github.com/nontypeable/financial-tracker/internal/domain/user"
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

func (s *service) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	token, err := s.tokenManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(uuid.Must(uuid.Parse(token.Subject)))
	if err != nil {
		return "", "", err
	}

	refreshToken, err = s.tokenManager.GenerateRefreshToken(uuid.Must(uuid.Parse(token.Subject)))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *service) GetUserInfo(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	return s.repository.GetByID(ctx, userID)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, firstName, lastName string) error {
	user, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	var updated bool

	if firstName != "" && firstName != user.FirstName {
		user.FirstName = firstName
	}

	if lastName != "" && lastName != user.LastName {
		user.LastName = lastName
	}

	if !updated {
		return nil
	}

	if err := s.repository.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *service) ChangeEmail(ctx context.Context, id uuid.UUID, newEmail string, currentPassword string) error {
	user, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return fmt.Errorf("invalid credentials")
	}

	exists, err := s.repository.EmailExists(ctx, newEmail)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already in use")
	}

	user.Email = newEmail

	return s.repository.Update(ctx, user)
}

func (s *service) ChangePassword(ctx context.Context, id uuid.UUID, newPassword string, currentPassword string) error {
	user, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return fmt.Errorf("invalid credentials")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashed)

	return s.repository.Update(ctx, user)
}
