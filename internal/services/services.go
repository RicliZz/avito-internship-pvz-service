package services

import "github.com/gin-gonic/gin"

type AuthenticationService interface {
	DummyLogin(ctx *gin.Context)
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}
