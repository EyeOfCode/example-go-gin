package handlers

import (
	"github.com/gin-gonic/gin"
)

// @Description Health check response
type HealthHandler struct {
    Status string `json:"status" example:"ok"`
}

// @Summary Health check endpoint
// @Description Get the API's health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthHandler
// @Router /api/v1/health [get]
func (s *HealthHandler) HealthCheck(c *gin.Context) {
    c.JSON(200, HealthHandler{Status: "ok"})
}