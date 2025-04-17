package productService

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/internal/repositories"
	"github.com/gin-gonic/gin"
	"log"
)

type ProductService struct {
	ReceptionRepository repositories.ReceptionRepo
	ProductRepo         repositories.ProductRepo
}

func NewProductService(ReceptionRepository repositories.ReceptionRepo, ProductRepo repositories.ProductRepo) *ProductService {
	return &ProductService{
		ReceptionRepository: ReceptionRepository,
		ProductRepo:         ProductRepo,
	}
}

func (s *ProductService) AddProductInReception(c *gin.Context) {
	log.Println("Запуск сервиса для добавления продукта в приёмку")
	var product models.AddProductRequest
	if err := c.ShouldBindJSON(&product); err != nil {
		log.Println("Ошибка при парсинге")
		c.JSON(400, gin.H{"description": "Неверный запрос или нет активной приёмки"})
		return
	}

	err, receptionID := s.ReceptionRepository.FindLastActiveReception(product.PvzID)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"description": "Неверный запрос или нет активной приёмки"})
		return
	}
	err, newProduct := s.ProductRepo.AddProductInActiveReception(receptionID, product.Type)
	if err != nil {
		log.Println("Ошибка при внесении информации о продукте в приёмку")
		c.JSON(400, gin.H{"description": err.Error()})
		return
	}
	c.JSON(201, newProduct)
}
