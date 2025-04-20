package receptionService

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
)

func (s *ReceptionService) CloseLastReception(c *gin.Context) {
	logger.Logger.Info("CloseLastReception service was started")
	stringPVZID := c.Param("pvzId")
	uuidPVZID, err := uuid.Parse(stringPVZID)
	if err != nil {
		logger.Logger.Debug("Validation failed")
		c.JSON(400, models.Error{Message: "Invalid request"})
		return
	}
	err, closedReception := s.ReceptionRepo.CloseLastReception(uuidPVZID)
	if err != nil {
		log.Println("Failed to close reception")
		c.JSON(400, models.Error{Message: "Invalid request"})
		return
	}
	c.JSON(200, closedReception)
}
