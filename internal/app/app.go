package app

import (
	"context"
	authentication "github.com/RicliZz/avito-internship-pvz-service/internal/api/auth"
	"github.com/RicliZz/avito-internship-pvz-service/internal/api/products"
	"github.com/RicliZz/avito-internship-pvz-service/internal/api/pvz"
	"github.com/RicliZz/avito-internship-pvz-service/internal/api/reception"
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/internal/repositories/authRepo"
	"github.com/RicliZz/avito-internship-pvz-service/internal/repositories/productRepo"
	"github.com/RicliZz/avito-internship-pvz-service/internal/repositories/pvzRepo"
	"github.com/RicliZz/avito-internship-pvz-service/internal/repositories/receptionRepo"
	"github.com/RicliZz/avito-internship-pvz-service/internal/server"
	"github.com/RicliZz/avito-internship-pvz-service/internal/services/authService"
	"github.com/RicliZz/avito-internship-pvz-service/internal/services/productService"
	"github.com/RicliZz/avito-internship-pvz-service/internal/services/pvzService"
	"github.com/RicliZz/avito-internship-pvz-service/internal/services/receptionService"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	defer logger.Logger.Sync()
	//Загрузка в переменные окружения .env
	err := godotenv.Load()
	if err != nil {
		logger.Logger.Fatal("Error loading .env file")
	}

	//Подключение к Постгресу
	conn, err := pgx.Connect(context.Background(), os.Getenv("POSTGRESQL_URL"))
	if err != nil {
		logger.Logger.Infow("Unable to connect to database:",
			"error", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	r := gin.Default()
	//Кастомные валидаторы
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("datesForGetPVZList", models.DatesForGetPVZList)
	}
	//Инициализация репозиториев(БД)
	authRepository := authRepo.NewAuthRepository(conn)
	PVZRepository := pvzRepo.NewPVZRepository(conn)
	receptionRepository := receptionRepo.NewReceptionRepository(conn)
	productRepository := productRepo.NewProductRepository(conn)

	//Инициализация сервисов
	loginService := authService.NewAuthLogin(authRepository)
	PVZService := pvzService.NewPVZService(PVZRepository)
	ReceptionService := receptionService.NewReceptionService(receptionRepository)
	ProductService := productService.NewProductService(receptionRepository, productRepository)

	//Инициализация ручек
	api := r.Group("")
	authHandlers := authentication.NewAuthHandler(loginService)
	PVZHandlers := pvz.NewPVZHandler(PVZService, ReceptionService)
	receptionHandlers := reception.NewReceptionHandlers(ReceptionService)
	productHandlers := products.NewProductsHandlers(ProductService)

	//Привязка ручек
	authHandlers.InitAuthHandlers(api)
	PVZHandlers.InitPVZHandlers(api)
	receptionHandlers.InitReceptionHandlers(api)
	productHandlers.InitProductsHandlers(api)

	//Инициализация и конфигурация HTTP сервера
	srv := server.NewAPIServer(r)

	//Старт сервера
	go srv.Start()

	//Выключение
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		logger.Logger.Fatalw("Shutdown error",
			"error", err)
	}

}
