package transaction

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Service interface {
	Create(ctx context.Context, accountID uuid.UUID, amount decimal.Decimal, transactionType TransactionType, description string) (uuid.UUID, error)
}
