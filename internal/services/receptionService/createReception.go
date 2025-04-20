package receptionService

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var CountCreatedReception = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "created_reception_count",
	Help: "Total number of created Receptions",
})

func (s *ReceptionService) CreateReception(c *gin.Context) {
	logger.Logger.Info("CreateReception service was started")
	var reception models.CreateReceptionRequest
	if err := c.ShouldBindJSON(&reception); err != nil {
		logger.Logger.Debug("Validation failed")
		c.JSON(400, models.Error{Message: "Invalid request"})
		return
	}
	newReception, err := s.ReceptionRepo.CreateReception(reception)
	if err != nil {
		logger.Logger.Info("Error creating reception")
		c.JSON(400, models.Error{Message: "Invalid request"})
		return
	}
	CountCreatedReception.Inc()
	c.JSON(201, newReception)
}
