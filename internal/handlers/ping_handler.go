package handlers

import (
	"context"
	"example-go-project/internal/dto"
	serviceRepository "example-go-project/internal/repository/service"
	"example-go-project/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PingHandler struct {
	serviceRepo serviceRepository.ServiceRepository
}

func NewPingHandler(serviceRepo serviceRepository.ServiceRepository) *PingHandler {
	return &PingHandler{
		serviceRepo: serviceRepo,
	}
}

// @Summary Ping endpoint
// @Description Post the API's ping
// @Tags ping
// @Accept json
// @Produce json
// @Param request body dto.PingRequest true "Ping details"
// @Router /ping [post]
func (p *PingHandler) Ping(c *gin.Context) {
	var req dto.PingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errors,
			})
			return
		}
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.serviceRepo.Get(c, req.Url)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.SendSuccess(c, http.StatusOK, nil, "pong")
}
