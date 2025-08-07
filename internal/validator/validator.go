package validator

import (
	"fmt"
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
		v := validator.New()
		if err := v.RegisterValidation("password", custom.ValidatePassword); err != nil {
			panic(fmt.Sprintf("failed to register password validation: %v", err))
		}
		instance = &Validator{validate: v}
	})
	return instance
}

func (v *Validator) ValidateStruct(s any) error {
	return v.validate.Struct(s)
}
