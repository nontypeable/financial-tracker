package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nontypeable/financial-tracker/internal/domain/transaction"
	apperror "github.com/nontypeable/financial-tracker/internal/errors"
	"github.com/shopspring/decimal"
)

type service struct {
	repository transaction.Repository
}

func NewService(repository transaction.Repository) transaction.Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Create(ctx context.Context, accountID uuid.UUID, amount decimal.Decimal, transactionType transaction.TransactionType, description string) (uuid.UUID, error) {
	transaction := transaction.NewTransaction(accountID, amount, transactionType, description)

	id, err := s.repository.Create(ctx, transaction)
	if err != nil {
		if errors.Is(err, apperror.ErrInvalidInput) {
			return uuid.Nil, apperror.ErrInvalidInput
		}
		return uuid.Nil, fmt.Errorf("create transaction: %w", err)
	}

	return id, nil
}
