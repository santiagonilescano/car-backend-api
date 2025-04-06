package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"car-service/cmd/api/controllers"
	"car-service/cmd/api/routes"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
	router     *gin.Engine
}

type ServerConfig struct {
	CarController *controllers.CarController
	Port          string
}

func NewServer(config *ServerConfig) *Server {
	router := gin.Default()

	routesConfig := &routes.Config{
		CarController: config.CarController,
	}
	routes.SetupRoutes(router, routesConfig)

	return &Server{
		router: router,
		httpServer: &http.Server{
			Addr:    ":" + config.Port,
			Handler: router,
		},
	}
}

func (s *Server) Start() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error al iniciar el servidor: %v\n", err)
		}
	}()

	<-quit
	log.Println("Apagando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Servidor apagado correctamente")
	return nil
}
