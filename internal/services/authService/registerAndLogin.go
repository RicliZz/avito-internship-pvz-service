package authService

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/internal/repositories"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/JWT"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/pass"
	"github.com/gin-gonic/gin"
)

type AuthLogin struct {
	DummyLoginService
	authDB repositories.AuthenticationRepo
}

func NewAuthLogin(authDB repositories.AuthenticationRepo) *AuthLogin {
	return &AuthLogin{
		authDB: authDB,
	}
}

func (s *AuthLogin) Login(c *gin.Context) {
	logger.Logger.Info("Login service was started")

	user := &models.LoginParams{}
	if err := c.ShouldBind(user); err != nil {
		logger.Logger.Debug("Validation failed",
			"email", user.Email,
			"password", user.Password)
		c.JSON(401, models.Error{Message: "Invalid credentials"})
		return
	}

	err, password, role := s.authDB.GetUserByEmail(user.Email)
	if err != nil {
		c.JSON(401, models.Error{Message: "Invalid credentials"})
		return
	}
	compare := pass.ComparePassWithHash(user.Password, password)
	if !compare {
		logger.Logger.Debug("Password not compare")
		c.JSON(401, models.Error{Message: "Invalid credentials"})
		return
	}
	token, err := JWT.CreateJWT(role)
	if err != nil {
		logger.Logger.Error("Failed create JWT token")
		return
	}
	c.JSON(200, token)
}

func (s *AuthLogin) Register(c *gin.Context) {
	logger.Logger.Info("Register service was started")
	var user models.RegisterParams
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.Logger.Debug("Validation failed")
		c.JSON(400, models.Error{Message: "Invalid request"})
		return
	}

	//Хэшируем пароль и заменяем изначальный на его хэш в модели
	user.Password = pass.CreateHash(user.Password)

	//Создание пользователя
	err, newUser := s.authDB.Register(user)
	if err != nil {
		logger.Logger.Error("Error when creating user")
		c.JSON(400, models.Error{Message: err.Error()})
		return
	}
	c.JSON(201, newUser)
}
