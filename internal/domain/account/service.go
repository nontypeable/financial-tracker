package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, name string, balance decimal.Decimal) (uuid.UUID, error)
}
