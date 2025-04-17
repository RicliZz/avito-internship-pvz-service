package services

import "github.com/gin-gonic/gin"

type AuthenticationService interface {
	DummyLogin(ctx *gin.Context)
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type PVZService interface {
	CreatePVZ(ctx *gin.Context)
}

type ReceptionService interface {
	CreateReception(ctx *gin.Context)
}

type ProductService interface {
	AddProductInReception(ctx *gin.Context)
	FindLastActiveReception(ctx *gin.Context)
}
