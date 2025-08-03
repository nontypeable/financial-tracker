package custom

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	length := len(password)

	if length < 8 || length > 72 {
		return false
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char), unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}
