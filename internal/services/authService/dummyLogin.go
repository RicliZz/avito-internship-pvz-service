package authService

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/JWT"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/gin-gonic/gin"
)

type DummyLoginService struct {
}

func (s *DummyLoginService) DummyLogin(c *gin.Context) {
	logger.Logger.Info("DummyLogin service was started")
	params := &models.DummyLoginParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Logger.Debugw("Validation failed",
			"role", params.Role)
		c.JSON(400, models.Error{Message: err.Error()})
		return
	}
	token, err := JWT.CreateJWT(params.Role)
	if err != nil {
		c.JSON(400, models.Error{Message: err.Error()})
		return
	}
	c.JSON(200, token)
}
