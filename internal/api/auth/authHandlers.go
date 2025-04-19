package authentication

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/services"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService services.AuthenticationService
}

func NewAuthHandler(AuthService services.AuthenticationService) AuthHandler {
	return AuthHandler{
		AuthService: AuthService,
	}
}

func (h *AuthHandler) InitAuthHandlers(router *gin.RouterGroup) {
	router.Use(middleware.RequestCounterMiddleware())
	{
		router.POST("/dummyLogin", h.AuthService.DummyLogin)
		router.POST("/register", h.AuthService.Register)
		router.POST("/login", h.AuthService.Login)
	}
}
