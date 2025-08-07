package account

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nontypeable/financial-tracker/internal/domain/account"
	"github.com/shopspring/decimal"
)

type service struct {
	repository account.Repository
}

func NewService(repository account.Repository) account.Service {
	return &service{repository: repository}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, name string, balance decimal.Decimal) (uuid.UUID, error) {
	account := account.NewAccount(userID, name, balance)

	accountID, err := s.repository.Create(ctx, account)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create account: %w", err)
	}

	return accountID, nil
}
