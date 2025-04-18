package receptionService

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
)

func (s *ReceptionService) CloseLastReception(c *gin.Context) {
	log.Println("Запуск сервиса для закрытия последней открытой приёмки")
	stringPVZID := c.Param("pvzId")
	uuidPVZID, err := uuid.Parse(stringPVZID)
	if err != nil {
		log.Println("Ошибка при парсе параметра в формат UUID")
		c.JSON(400, err)
	}
	err, reception := s.ReceptionRepo.CloseLastReception(uuidPVZID)
	if err != nil {
		log.Println("Не удалось закрыть приёмку")
		c.JSON(400, err)
		return
	}
	c.JSON(200, reception)
}
