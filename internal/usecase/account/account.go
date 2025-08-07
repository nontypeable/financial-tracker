package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nontypeable/financial-tracker/internal/domain/account"
	apperror "github.com/nontypeable/financial-tracker/internal/errors"
	"github.com/shopspring/decimal"
)

type service struct {
	repository account.Repository
}

func NewService(repository account.Repository) account.Service {
	return &service{repository: repository}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, name string, balance decimal.Decimal) error {
	account := account.NewAccount(userID, name, balance)

	err := s.repository.Create(ctx, account)
	if err != nil {
		if errors.Is(err, apperror.ErrInvalidInput) {
			return apperror.ErrInvalidInput
		}
		return fmt.Errorf("create account: %w", err)
	}

	return nil
}
