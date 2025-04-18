package pvzService

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (s *PVZService) GetPVZList(c *gin.Context) {
	queryParams := models.QueryParamForGetPVZList{
		Page:  1,
		Limit: 10,
	}
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		log.Println("Ошибка при парсе параметров")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var allPVZ []models.ListPVZResponse
	err, allPVZ := s.PVZRepo.GetListPVZ(queryParams)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": allPVZ})
}
