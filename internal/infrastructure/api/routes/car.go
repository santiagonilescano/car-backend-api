package routes

import (
	"car-service/internal/infrastructure/api/handlers"

	"github.com/gin-gonic/gin"
)

// CarRoutes configura las rutas relacionadas con cars
func SetupCarRoutes(router *gin.RouterGroup, carHandler *handlers.CarHandler) {
	cars := router.Group("/cars")
	{
		cars.POST("", carHandler.CreateCar)
		cars.GET("", carHandler.GetCars)
		// Aquí irán más rutas relacionadas con cars:
		// cars.GET("", carHandler.ListCars)
		// cars.GET("/:id", carHandler.GetCar)
		// cars.PUT("/:id", carHandler.UpdateCar)
		// cars.DELETE("/:id", carHandler.DeleteCar)
	}
}
