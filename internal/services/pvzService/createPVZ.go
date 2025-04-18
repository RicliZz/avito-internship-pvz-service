package pvzService

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (s *PVZService) CreatePVZ(c *gin.Context) {
	logger.Logger.Info("CreatePVZ service was started")
	var PVZ models.CreatePVZRequest
	if err := c.ShouldBindJSON(&PVZ); err != nil {
		logger.Logger.Debugw("Validation failed",
			"city", PVZ.City)
		c.JSON(400, models.Error{Message: "Invalid request"})
		return
	}
	err, newPVZ := s.PVZRepo.CreatePVZ(PVZ)
	if err != nil {
		logger.Logger.Error("Error when creating PVZ")
		c.JSON(400, models.Error{Message: err.Error()})
		return
	}
	c.JSON(201, newPVZ)
}
