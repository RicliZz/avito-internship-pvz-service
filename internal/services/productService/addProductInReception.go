package productService

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/internal/repositories"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/gin-gonic/gin"
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
	logger.Logger.Info("AddProductInReception service was started")
	var product models.AddProductRequest
	if err := c.ShouldBindJSON(&product); err != nil {
		logger.Logger.Debugw("Validation failed",
			"productType", product.Type,
			"pvzID", product.PvzID)
		c.JSON(400, gin.H{"description": "Неверный запрос или нет активной приёмки"})
		return
	}

	err, receptionID := s.ReceptionRepository.FindLastActiveReception(product.PvzID)
	if err != nil {
		c.JSON(400, gin.H{"description": "Неверный запрос или нет активной приёмки"})
		return
	}
	err, newProduct := s.ProductRepo.AddProductInActiveReception(receptionID, product.Type)
	if err != nil {
		c.JSON(400, gin.H{"description": err.Error()})
		return
	}
	c.JSON(201, newProduct)
}
