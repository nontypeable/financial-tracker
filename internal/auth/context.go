package auth

import (
	"context"

	"github.com/google/uuid"
	contextKeys "github.com/nontypeable/financial-tracker/internal/context"
)

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(contextKeys.UserIDKey).(string)
	if !ok {
		return uuid.UUID{}, false
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return uuid.UUID{}, false
	}

	return id, true
}
