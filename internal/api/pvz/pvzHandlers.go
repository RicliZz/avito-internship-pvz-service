package pvz

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/services"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/middleware"
	"github.com/gin-gonic/gin"
	"os"
)

type PVZHandler struct {
	pvzService       services.PVZService
	receptionService services.ReceptionService
}

func NewPVZHandler(pvzService services.PVZService, receptionService services.ReceptionService) *PVZHandler {
	return &PVZHandler{
		pvzService:       pvzService,
		receptionService: receptionService,
	}
}

func (h *PVZHandler) InitPVZHandlers(router *gin.RouterGroup) {
	pvzModeratorRouter := router.Group("/pvz")
	pvzModeratorRouter.Use(middleware.CheckRoleMiddleware(os.Getenv("JWT_SECRET"), "moderator"))

	{
		pvzModeratorRouter.POST("", h.pvzService.CreatePVZ)
		pvzModeratorRouter.GET("", h.pvzService.GetPVZList)
		pvzModeratorRouter.GET("/grpc", h.pvzService.GetPVZListFromRPCServer)
	}

	pvzEmployeeRouter := router.Group("/pvz")
	pvzEmployeeRouter.Use(middleware.CheckRoleMiddleware(os.Getenv("JWT_SECRET"), "employee"))

	{
		pvzEmployeeRouter.POST("/:pvzId/delete_last_product", h.receptionService.DeleteLastProductInActiveReception)
		pvzEmployeeRouter.POST("/:pvzId/close_last_reception", h.receptionService.CloseLastReception)
	}

}
