package dto

import (
	"github.com/nontypeable/financial-tracker/internal/validator"
	"github.com/shopspring/decimal"
)

type CreateRequest struct {
	Name    string          `json:"name"`
	Balance decimal.Decimal `json:"balance"`
}

func (r *CreateRequest) Validate() error {
	return validator.GetValidator().ValidateStruct(r)
}
