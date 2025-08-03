package dto

import "github.com/nontypeable/financial-tracker/internal/validator"

type UpdateRequest struct {
	FirstName string `json:"first_name" validate:"omitempty,min=1,max=50"`
	LastName  string `json:"last_name" validate:"omitempty,min=1,max=50"`
}

func (r *UpdateRequest) Validate() error {
	return validator.GetValidator().ValidateStruct(r)
}
