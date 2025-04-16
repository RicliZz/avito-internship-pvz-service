package authService

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/internal/repositories"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/JWT"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/pass"
	"github.com/gin-gonic/gin"
	"log"
)

type AuthLogin struct {
	DummyLoginService
	authDB repositories.AuthenticationRepository
}

func NewAuthLogin(authDB repositories.AuthenticationRepository) *AuthLogin {
	return &AuthLogin{
		authDB: authDB,
	}
}

func (s *AuthLogin) Login(c *gin.Context) {
	log.Println("Сервис Login")

	params := &models.LoginParams{}
	if err := c.ShouldBind(params); err != nil {
		log.Println("Не прошла валидация")
		c.JSON(400, gin.H{"description": "Неверный запрос"})
		return
	}

	password, role, err := s.authDB.GetUserByEmail(params.Email)
	if err != nil {
		log.Println("Ошибка при получении пароля для сравнения из БД")
		c.JSON(400, gin.H{"description": "Неверный запрос"})
		return
	}
	compare := pass.ComparePassWithHash(params.Password, password)
	if !compare {
		log.Println("Пароли не одинаковы")
		c.JSON(401, gin.H{"description": "Неверные учётные данные"})
		return
	}
	token, err := JWT.CreateJWT(role)
	if err != nil {
		log.Println("Ошибка при создании токена")
		return
	}
	c.JSON(200, gin.H{
		"description": "Успешная авторизация",
		"Token":       token},
	)
}

func (s *AuthLogin) Register(c *gin.Context) {
	log.Println("Сервис регистрации")
	params := &models.RegisterParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		log.Println("Не прошла валидация")
		c.JSON(400, gin.H{"description": "Неверный запрос"})
		return
	}

	//hash password and change password on hash in model
	params.Password = pass.CreateHash(params.Password)

	//create new user in db
	err := s.authDB.Register(*params)
	if err != nil {
		log.Println("Register repository fail")
		c.JSON(400, gin.H{"description": "Неверный запрос"})
		return
	}
	c.JSON(201, gin.H{"description": "Пользователь создан"})
}
