package account

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, account *Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*Account, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Account, error)
	Update(ctx context.Context, account *Account) error
	Delete(ctx context.Context, userID, id uuid.UUID) error
}
