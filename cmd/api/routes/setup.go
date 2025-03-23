package routes

import (
	"car-service/cmd/api/controllers"

	"github.com/gin-gonic/gin"
)

// Config contiene todas las dependencias necesarias para las rutas
type Config struct {
	CarController *controllers.CarController
	// Aquí se agregarán más handlers según sea necesario
}

// SetupRoutes configura todas las rutas de la API
func SetupRoutes(router *gin.Engine, config *Config) {
	// Grupo base para la API v1
	v1 := router.Group("/api/v1")

	// Configurar rutas de cars
	SetupCarRoutes(v1, *config.CarController)

	// Aquí se agregarán más configuraciones de rutas
	// SetupUserRoutes(v1, config.UserHandler)
	// SetupAuthRoutes(v1, config.AuthHandler)
	// etc.
}
