package products

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/services"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/middleware"
	"github.com/gin-gonic/gin"
	"os"
)

type Product struct {
	ProductService services.ProductService
}

func NewProductsHandlers(ProductService services.ProductService) *Product {
	return &Product{
		ProductService: ProductService,
	}
}

func (h *Product) InitProductsHandlers(router *gin.RouterGroup) {
	productsRouter := router.Group("/products")
	productsRouter.Use(middleware.CheckRoleMiddleware(os.Getenv("JWT_SECRET"), "employee"))
	{
		productsRouter.POST("", h.ProductService.AddProductInReception)
	}
}
