package dto

import "github.com/nontypeable/financial-tracker/internal/validator"

type UpdateEmailRequest struct {
	Email    string `json:"email" validate:"required,min=6,max=254"`
	Password string `json:"password" validate:"required,min=8,max=72,password"`
}

func (r *UpdateEmailRequest) Validate() error {
	return validator.GetValidator().ValidateStruct(r)
}
