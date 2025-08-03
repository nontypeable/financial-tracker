package dto

import "github.com/nontypeable/financial-tracker/internal/validator"

type SignInRequest struct {
	Email    string `json:"username" validate:"required,min=6,max=254"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

func (r *SignInRequest) Validate() error {
	return validator.GetValidator().ValidateStruct(r)
}

type SignInResponse struct {
	AccessToken string `json:"access_token"`
}
