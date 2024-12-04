package utils

import (
	"errors"
	"fmt"
	"unicode"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func SetupValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom password validator
		v.RegisterValidation("passwordvalidator", func(fl validator.FieldLevel) bool {
			password := fl.Field().String()

			hasUpper := false
			hasNumber := false
			hasSpecial := false

			for _, c := range password {
				switch {
				case unicode.IsUpper(c):
					hasUpper = true
				case unicode.IsNumber(c):
					hasNumber = true
				case unicode.IsPunct(c) || unicode.IsSymbol(c):
					hasSpecial = true
				}
			}

			return hasUpper && hasNumber && hasSpecial
		})
	}
}

func FormatValidationError(err error) []string {
    var validationErrors validator.ValidationErrors
    errorMessages := make([]string, 0)
    
    if errors.As(err, &validationErrors) {
        for _, e := range validationErrors {
            switch e.Tag() {
            case "required":
                errorMessages = append(errorMessages, fmt.Sprintf("%s is required", e.Field()))
            case "email":
                errorMessages = append(errorMessages, "Invalid email format")
            case "min":
                errorMessages = append(errorMessages, fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param()))
            case "max":
                errorMessages = append(errorMessages, fmt.Sprintf("%s must not exceed %s characters", e.Field(), e.Param()))
            case "passwordvalidator":
                errorMessages = append(errorMessages, "Password must contain at least one uppercase letter, one number, and one special character")
            }
        }
    }
    
    return errorMessages
}