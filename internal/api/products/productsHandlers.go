package products

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/services"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/middleware"
	"github.com/gin-gonic/gin"
	"os"
)

type ProductsHandlers struct {
	ProductService services.ProductService
}

func NewProductsHandlers(ProductService services.ProductService) *ProductsHandlers {
	return &ProductsHandlers{
		ProductService: ProductService,
	}
}

func (h *ProductsHandlers) InitProductsHandlers(router *gin.RouterGroup) {
	productsRouter := router.Group("/products")
	productsRouter.Use(middleware.RequestCounterMiddleware())
	productsRouter.Use(middleware.CheckRoleMiddleware(os.Getenv("JWT_SECRET"), "employee"))
	{
		productsRouter.POST("", h.ProductService.AddProductInReception)
	}
}
