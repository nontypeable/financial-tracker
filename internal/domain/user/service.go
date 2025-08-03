package user

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	SignUp(ctx context.Context, email, password, firstName, lastName string) (string, string, error)
	SignIn(ctx context.Context, email, password string) (string, string, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
	Update(ctx context.Context, id uuid.UUID, firstName, lastName string) error
}
