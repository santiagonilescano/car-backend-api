package routes

import (
	"car-service/cmd/api/controllers"

	"github.com/gin-gonic/gin"
)

type Config struct {
	CarController *controllers.CarController
}

func SetupRoutes(router *gin.Engine, config *Config) {
	v1 := router.Group("/api/v1")
	SetupCarRoutes(v1, *config.CarController)
}
