package receptionService

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/internal/repositories"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type ReceptionService struct {
	ReceptionRepo repositories.ReceptionRepo
}

func NewReceptionService(receptionRepo repositories.ReceptionRepo) *ReceptionService {
	return &ReceptionService{
		ReceptionRepo: receptionRepo,
	}
}

func (s *ReceptionService) CreateReception(c *gin.Context) {
	log.Println("Создание приёмки")
	var reception models.CreateReceptionRequest
	if err := c.ShouldBindJSON(&reception); err != nil {
		log.Println("Ошибка при парсе JSON")
		c.JSON(http.StatusBadRequest, gin.H{"description": "Неверный запрос или есть незакрытая приёмка"})
		return
	}
	err, newReception := s.ReceptionRepo.CreateReception(reception)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"description": "Неверный запрос или есть незакрытая приемка"})
		return
	}
	c.JSON(201, gin.H{"Приёмка создана": newReception})
}
