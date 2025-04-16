package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

type APIServer struct {
	httpServer *http.Server
}

func NewAPIServer(router *gin.Engine) *APIServer {
	return &APIServer{
		httpServer: &http.Server{
			Addr:         os.Getenv("API_ADDR") + os.Getenv("API_PORT"),
			Handler:      router.Handler(),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *APIServer) Start() {
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

func (s *APIServer) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	switch ctx.Err() {
	case context.DeadlineExceeded:
		log.Println("Timeout shutting down server")
	case nil:
		log.Println("Shutdown completed before timeout.")
	default:
		log.Println("Shutdown ended with:", ctx.Err())
	}

	log.Println("Shutdown complete")
	return nil
}
