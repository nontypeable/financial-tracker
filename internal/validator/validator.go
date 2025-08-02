package validator

import (
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/nontypeable/financial-tracker/internal/validator/custom"
)

type Validator struct {
	validate *validator.Validate
}

var (
	instance *Validator
	once     sync.Once
)

func GetValidator() *Validator {
	once.Do(func() {
		instance = &Validator{
			validate: validator.New(),
		}
		instance.registerCustomValidators()
	})
	return instance
}

func (v *Validator) ValidateStruct(s any) error {
	return v.validate.Struct(s)
}

func (v *Validator) registerCustomValidators() {
	v.validate.RegisterValidation("password", custom.ValidatePassword)
}
