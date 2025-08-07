package dto

import (
	"github.com/google/uuid"
	"github.com/nontypeable/financial-tracker/internal/domain/transaction"
	"github.com/nontypeable/financial-tracker/internal/validator"
	"github.com/shopspring/decimal"
)

type CreateRequest struct {
	AccountID   uuid.UUID                   `json:"account_id"`
	Amount      decimal.Decimal             `json:"amount"`
	Type        transaction.TransactionType `json:"type"`
	Description string                      `json:"description"`
}

func (r *CreateRequest) Validate() error {
	return validator.GetValidator().ValidateStruct(r)
}

type CreateResponse struct {
	ID uuid.UUID `json:"id"`
}
