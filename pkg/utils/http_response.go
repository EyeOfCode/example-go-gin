package utils

import (
	"github.com/gin-gonic/gin"
)

// SendSuccess sends a successful JSON response
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

// SendError sends an error JSON response
func SendError(c *gin.Context, status int, message string) {
	response := gin.H{
		"success": false,
		"error":   message,
	}
	c.JSON(status, response)
}