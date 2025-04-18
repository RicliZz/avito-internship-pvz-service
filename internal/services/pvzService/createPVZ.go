package pvzService

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/gin-gonic/gin"
	"log"
)

func (s *PVZService) CreatePVZ(c *gin.Context) {
	log.Println("Началось создание нового ПВЗ")
	var PVZ models.CreatePVZRequest
	if err := c.ShouldBindJSON(&PVZ); err != nil {
		log.Println("Ошибка при валидации")
		c.JSON(400, gin.H{"description": "Неверный запрос"})
		return
	}
	err, newPVZ := s.PVZRepo.CreatePVZ(PVZ)
	if err != nil {
		log.Println("Ошибка при создании нового ПВЗ")
		c.JSON(400, gin.H{"description": err})
		return
	}
	c.JSON(201, gin.H{"ПВЗ создан": newPVZ})
}
