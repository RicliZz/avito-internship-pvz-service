package reception

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/services"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/middleware"
	"github.com/gin-gonic/gin"
	"os"
)

type ReceptionHandlers struct {
	ReceptionService services.ReceptionService
}

func NewReceptionHandlers(receptionService services.ReceptionService) *ReceptionHandlers {
	return &ReceptionHandlers{
		ReceptionService: receptionService,
	}
}

func (h *ReceptionHandlers) InitReceptionHandlers(router *gin.RouterGroup) {
	receptionRouter := router.Group("/receptions")
	receptionRouter.Use(middleware.RequestCounterMiddleware())
	receptionRouter.Use(middleware.CheckRoleMiddleware(os.Getenv("JWT_SECRET"), "employee"))
	{
		receptionRouter.POST("", h.ReceptionService.CreateReception)
	}
}
