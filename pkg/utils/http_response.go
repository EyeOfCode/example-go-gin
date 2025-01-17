package utils

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func SendSuccess(c *gin.Context, status int, data interface{}, message ...string) {
	response := gin.H{
		"success": true,
		"data":    data,
	}

	if len(message) > 0 {
		response["message"] = message[0]
	}

	c.JSON(status, response)
}
func SendError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"error":   message,
	})
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
			case "eqfield":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be equal to %s", e.Field(), e.Param()))
			case "password_validator":
				errorMessages = append(errorMessages, "Password must contain at least one uppercase letter, one number, and one special character")
			}
		}
	}

	return errorMessages
}