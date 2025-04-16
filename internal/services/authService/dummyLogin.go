package authService

import (
	"github.com/RicliZz/avito-internship-pvz-service/pkg/JWT"
	"github.com/gin-gonic/gin"
	"log"
)

type DummyLoginService struct {
}

func NewDummyLogin() *DummyLoginService {
	return &DummyLoginService{}
}

type DummyLoginParams struct {
	Role string ` json:"role" binding:"required,oneof=employee moderator"` //role must be employee / moderator, else error
}

func (s *DummyLoginService) DummyLogin(c *gin.Context) {
	log.Println("Dummy login service start")
	params := &DummyLoginParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		log.Println("Invalid request")
		c.JSON(400, gin.H{"description": "Неверный запрос"})
		return
	}
	token, err := JWT.CreateDummyJWT(params.Role)
	if err != nil {
		c.JSON(400, gin.H{"description": err.Error()})
		return
	}
	c.JSON(200, gin.H{"token": token})
}
