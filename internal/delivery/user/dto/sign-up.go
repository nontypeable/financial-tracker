package dto

import "github.com/nontypeable/financial-tracker/internal/validator"

type SignUpRequest struct {
	Email     string `json:"email" validate:"required,min=6,max=254"`
	Password  string `json:"password" validate:"required,min=8,max=72,password"`
	FirstName string `json:"first_name" validate:"required,min=1,max=50"`
	LastName  string `json:"last_name" validate:"required,min=1,max=50"`
}

func (r *SignUpRequest) Validate() error {
	return validator.GetValidator().ValidateStruct(r)
}

type SignUpResponse struct {
	AccessToken string `json:"access_token"`
}
