package pvz

import "github.com/gin-gonic/gin"

type PVZHandler struct {
}

func NewPVZHandler() *PVZHandler {
	return &PVZHandler{}
}

func (h *PVZHandler) InitPVZHandlers(router *gin.RouterGroup) {
	router.Use()
}
