package receptionService

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
)

func (s *ReceptionService) DeleteLastProductInActiveReception(c *gin.Context) {
	log.Println("Запуск сервиса по удалению последнего товара в приёмке")
	stringPVZID := c.Param("pvzId")
	uuidPVZID, err := uuid.Parse(stringPVZID)
	if err != nil {
		log.Println("Ошибка при парсе параметра в формат UUID")
		c.JSON(400, err)
	}
	if err = s.ReceptionRepo.DeleteLastProduct(uuidPVZID); err != nil {
		log.Println("Ошибка в удалении")
		c.JSON(400, err.Error())
		return
	}
	c.JSON(200, "Товар удалён")
}
