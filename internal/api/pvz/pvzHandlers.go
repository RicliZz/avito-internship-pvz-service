package pvz

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/services"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/middleware"
	"github.com/gin-gonic/gin"
	"os"
)

type PVZHandler struct {
	pvzService services.PVZService
}

func NewPVZHandler(pvzService services.PVZService) *PVZHandler {
	return &PVZHandler{
		pvzService: pvzService,
	}
}

func (h *PVZHandler) InitPVZHandlers(router *gin.RouterGroup) {
	router.Use(middleware.CheckRoleMiddleware(os.Getenv("JWT_SECRET"), "moderator"))
	{
		router.POST("/pvz", h.pvzService.CreateNewPVZ)
	}
}
