package transaction

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nontypeable/financial-tracker/internal/domain/transaction"
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
		return uuid.Nil, fmt.Errorf("create transaction: %w", err)
	}

	return id, nil
}
