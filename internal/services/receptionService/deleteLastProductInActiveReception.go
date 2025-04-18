package receptionService

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *ReceptionService) DeleteLastProductInActiveReception(c *gin.Context) {
	logger.Logger.Info("DeleteLastProductInActiveReception service was started")
	stringPVZID := c.Param("pvzId")
	uuidPVZID, err := uuid.Parse(stringPVZID)
	if err != nil {
		logger.Logger.Debug("Validation failed")
		c.JSON(400, models.Error{Message: "Invalid request"})
	}
	if err = s.ReceptionRepo.DeleteLastProduct(uuidPVZID); err != nil {
		logger.Logger.Info("Error deleting last product")
		c.JSON(400, err.Error())
		return
	}
	c.JSON(200, "Товар удалён")
}
