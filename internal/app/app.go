package app

import (
	"context"
	"fmt"
	authentication "github.com/RicliZz/avito-internship-pvz-service/internal/api/auth"
	"github.com/RicliZz/avito-internship-pvz-service/internal/api/products"
	"github.com/RicliZz/avito-internship-pvz-service/internal/api/pvz"
	"github.com/RicliZz/avito-internship-pvz-service/internal/api/reception"
	"github.com/RicliZz/avito-internship-pvz-service/internal/repositories/authRepo"
	"github.com/RicliZz/avito-internship-pvz-service/internal/repositories/productRepo"
	"github.com/RicliZz/avito-internship-pvz-service/internal/repositories/pvzRepo"
	"github.com/RicliZz/avito-internship-pvz-service/internal/repositories/receptionRepo"
	"github.com/RicliZz/avito-internship-pvz-service/internal/server"
	"github.com/RicliZz/avito-internship-pvz-service/internal/services/authService"
	"github.com/RicliZz/avito-internship-pvz-service/internal/services/productService"
	"github.com/RicliZz/avito-internship-pvz-service/internal/services/pvzService"
	"github.com/RicliZz/avito-internship-pvz-service/internal/services/receptionService"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	//init config .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//connect to PostgreSQL
	conn, err := pgx.Connect(context.Background(), os.Getenv("POSTGRESQL_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	//Engine GIN
	r := gin.Default()

	//Initialize repositories
	authRepository := authRepo.NewAuthRepository(conn)
	PVZRepository := pvzRepo.NewPVZRepository(conn)
	receptionRepository := receptionRepo.NewReceptionRepository(conn)
	productRepository := productRepo.NewProductRepository(conn)

	//Initialize services
	loginService := authService.NewAuthLogin(authRepository)
	PVZService := pvzService.NewPVZService(PVZRepository)
	ReceptionService := receptionService.NewReceptionService(receptionRepository)
	ProductService := productService.NewProductService(receptionRepository, productRepository)

	//New Handlers
	api := r.Group("")
	authHandlers := authentication.NewAuthHandler(loginService)
	PVZHandlers := pvz.NewPVZHandler(PVZService)
	receptionHandlers := reception.NewReceptionHandlers(ReceptionService)
	productHandlers := products.NewProductsHandlers(ProductService)

	//Init Handlers
	authHandlers.InitAuthHandlers(api)
	PVZHandlers.InitPVZHandlers(api)
	receptionHandlers.InitReceptionHandlers(api)
	productHandlers.InitProductsHandlers(api)

	// Initialize and configure the HTTP server
	srv := server.NewAPIServer(r)
	log.Printf("Starting server on %s%s", os.Getenv("SERVER_ADDRESS"), os.Getenv("PORT"))

	//start server
	go srv.Start()

	//close server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}

}
