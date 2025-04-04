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

// Server encapsula toda la configuración del servidor HTTP
type Server struct {
	httpServer *http.Server
	router     *gin.Engine
}

// ServerConfig contiene todas las dependencias necesarias para el servidor
type ServerConfig struct {
	CarController *controllers.CarController
	Port          string
}

// NewServer crea una nueva instancia del servidor
func NewServer(config *ServerConfig) *Server {
	router := gin.Default()

	// Configurar rutas
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

// Start inicia el servidor HTTP
func (s *Server) Start() error {
	// Canal para manejar señales de apagado
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor en una goroutine
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error al iniciar el servidor: %v\n", err)
		}
	}()

	// Esperar señal de apagado
	<-quit
	log.Println("Apagando servidor...")

	// Contexto con timeout para el apagado
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Intentar apagar el servidor gracefully
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Servidor apagado correctamente")
	return nil
}
