package app

import (
	"context"
	authentication "github.com/RicliZz/avito-internship-pvz-service/internal/api/auth"
	"github.com/RicliZz/avito-internship-pvz-service/internal/server"
	"github.com/RicliZz/avito-internship-pvz-service/internal/services/authService"
	"github.com/gin-gonic/gin"
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
	//Engine GIN
	r := gin.Default()

	//Initialize services
	dummyLoginService := authService.NewDummyLogin()

	//New Handlers
	api := r.Group("")
	authHandlers := authentication.NewAuthHandler(dummyLoginService)

	//Init Handlers
	authHandlers.InitAuthHandlers(api)

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
