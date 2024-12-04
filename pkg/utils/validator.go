package utils

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func SetupValidator() {
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        v.RegisterValidation("passwordvalidator", PasswordValidator)
    }
}

func PasswordValidator(fl validator.FieldLevel) bool {
    password := fl.Field().String()

    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
    hasSpecial := regexp.MustCompile(`[!@#$%^&*()+-_=\[\]{}|;:,.<>?]`).MatchString(password)

    return hasUpper && hasNumber && hasSpecial && hasLower
}