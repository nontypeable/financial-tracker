package account

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Account struct {
	ID        uuid.UUID       `db:"id"`
	UserID    uuid.UUID       `db:"user_id"`
	Name      string          `db:"name"`
	Balance   decimal.Decimal `db:"balance"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
	DeletedAt *time.Time      `db:"deleted_at"`
}

func NewAccount(userID uuid.UUID, name string, balance decimal.Decimal) *Account {
	return &Account{
		UserID:  userID,
		Name:    name,
		Balance: balance,
	}
}

func (a *Account) BelongsUser(userID uuid.UUID) bool {
	return a.UserID == userID
}

func (a *Account) Delete() {
	now := time.Now()
	a.DeletedAt = &now
	a.UpdatedAt = now
}

func (a *Account) Restore() {
	a.DeletedAt = nil
	a.UpdatedAt = time.Now()
}
