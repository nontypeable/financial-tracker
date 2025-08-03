package user

import "context"

type Service interface {
	SignUp(ctx context.Context, email, password, firstName, lastName string) (string, string, error)
}
