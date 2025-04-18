package authService

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/internal/repositories"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/JWT"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/pass"
	"github.com/gin-gonic/gin"
	"log"
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
		c.JSON(400, gin.H{"description": "Неверный запрос"})
		return
	}

	err, password, role := s.authDB.GetUserByEmail(user.Email)
	if err != nil {
		c.JSON(400, gin.H{"description": "Неверный запрос"})
		return
	}
	compare := pass.ComparePassWithHash(user.Password, password)
	if !compare {
		logger.Logger.Debug("Password not compare")
		c.JSON(401, gin.H{"description": "Неверные учётные данные"})
		return
	}
	token, err := JWT.CreateJWT(role)
	if err != nil {
		logger.Logger.Error("Failed create JWT token")
		return
	}
	c.JSON(200, gin.H{
		"description": "Успешная авторизация",
		"Token":       token},
	)
}

func (s *AuthLogin) Register(c *gin.Context) {
	logger.Logger.Info("Сервис Register")
	var user models.RegisterParams
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Не прошла валидация")
		c.JSON(400, gin.H{"description": "Неверный запрос"})
		return
	}

	//Хэшируем пароль и заменяем изначальный на его хэш в модели
	user.Password = pass.CreateHash(user.Password)

	//Создание пользователя
	err := s.authDB.Register(user)
	if err != nil {
		log.Println("Ошибка при создании пользователя в БД")
		c.JSON(400, gin.H{"description": "Неверный запрос"})
		return
	}
	c.JSON(201, gin.H{"description": "Пользователь создан"})
}
