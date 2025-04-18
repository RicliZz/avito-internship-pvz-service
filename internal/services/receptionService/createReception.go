package receptionService

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (s *ReceptionService) CreateReception(c *gin.Context) {
	logger.Logger.Info("CreateReception service was started")
	var reception models.CreateReceptionRequest
	if err := c.ShouldBindJSON(&reception); err != nil {
		logger.Logger.Debug("Validation failed")
		c.JSON(400, models.Error{Message: "Invalid request"})
		return
	}
	err, newReception := s.ReceptionRepo.CreateReception(reception)
	if err != nil {
		logger.Logger.Info("Error creating reception")
		c.JSON(400, models.Error{Message: "Invalid request"})
		return
	}
	c.JSON(201, newReception)
}
