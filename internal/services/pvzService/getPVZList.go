package pvzService

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"time"
)

func (s *PVZService) GetPVZList(c *gin.Context) {
	queryParams := models.QueryParamForGetPVZList{
		Page:      1,
		Limit:     10,
		StartDate: time.Now().AddDate(0, 0, -30),
		EndDate:   time.Now(),
	}
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		logger.Logger.Debugw("Validation failed",
			"startDate", queryParams.StartDate,
			"endDate", queryParams.EndDate)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var allPVZ []models.ListPVZResponse
	allPVZ, err := s.PVZRepo.GetListPVZ(queryParams)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, allPVZ)
}
