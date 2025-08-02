package transaction

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID          uuid.UUID
	AccountID   uuid.UUID
	Amount      decimal.Decimal
	Type        TransactionType
	Description string
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

func NewTransaction(accountID uuid.UUID, amount decimal.Decimal, transactionType TransactionType, description string) *Transaction {
	return &Transaction{
		AccountID:   accountID,
		Amount:      amount,
		Type:        transactionType,
		Description: description,
	}
}

func (t *Transaction) Delete() {
	now := time.Now()
	t.DeletedAt = &now
	t.UpdatedAt = now
}

func (t *Transaction) Restore() {
	t.DeletedAt = nil
	t.UpdatedAt = time.Now()
}
