package transaction

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, transaction *Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]*Transaction, error)
	Update(ctx context.Context, transaction *Transaction) error
	Delete(ctx context.Context, accountID, id uuid.UUID) error
}
