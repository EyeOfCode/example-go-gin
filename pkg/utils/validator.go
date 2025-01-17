package utils

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func SetupValidator() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("password_validator", PasswordValidator); err != nil {
			return fmt.Errorf("failed to register password validator: %w", err)
		}
	}
	return nil
}

func PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()+-_=\[\]{}|;:,.<>?]`).MatchString(password)

	return hasUpper && hasNumber && hasSpecial && hasLower
}
