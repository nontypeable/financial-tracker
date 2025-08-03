package auth

import (
	"context"

	"github.com/google/uuid"
)

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return uuid.UUID{}, false
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return uuid.UUID{}, false
	}

	return id, true
}
