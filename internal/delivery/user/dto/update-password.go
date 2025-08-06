package dto

import "github.com/nontypeable/financial-tracker/internal/validator"

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8,max=72,password"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=72,password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

func (r *UpdatePasswordRequest) Validate() error {
	return validator.GetValidator().ValidateStruct(r)
}
